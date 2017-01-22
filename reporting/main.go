package main

//go:generate go-bindata -pkg main -o bindata.go sample.html

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"encoding/json"
	_ "expvar"
)

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.RemoteAddr,
			time.Now().Format("2006-01-02T15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}

// see play.golang.org/p/PpbPyXbtEs
type flushWritter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWritter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return n, err
}

func measureRuntime(label string, block func()) {
	start := time.Now()
	block()
	log.Printf(
		"%s\t%s",
		label,
		time.Since(start),
	)
}

type context struct {
	fiscalPeriod
	metaData
}

func previewPDF(rw http.ResponseWriter, req *http.Request) {
	var templateID = req.URL.Path
	var data, _ = fiscalPeriodTemplate(templateID)
	var c context

	measureRuntime("umsatz api", func() {
		var fp fiscalPeriod
		var err error
		if fp, err = UmsatzClient.FetchFiscalPeriod(templateID); err != nil {
			log.Fatalf("Error loading data: %v", err)
		}
		c.fiscalPeriod = fp
		c.metaData = UmsatzClient.getMetaData()
	})

	html, err := template.New("pdf").Funcs(templateHelperFuncs).Parse(string(data))

	err = html.Execute(rw, &c)
	if err != nil {
		log.Fatalf("Error evaluating the template: %v", err)
	}

	var fw = flushWritter{w: rw}
	if f, ok := rw.(http.Flusher); ok {
		fw.f = f
	}
}

type templateResponse struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func fiscalPeriodTemplate(id string) ([]byte, error) {
	var f, err = os.Open(fmt.Sprintf("%v/%v.html", TemplateDirectory, id))
	if err != nil {
		return Asset("sample.html")
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func templateHandler(rw http.ResponseWriter, req *http.Request) {
	var templateID, _ = strconv.Atoi(req.URL.Path)
	if req.Method == "GET" {
		var b, _ = fiscalPeriodTemplate(req.URL.Path)
		var t = templateResponse{
			ID:   templateID,
			Data: string(b),
		}

		var enc = json.NewEncoder(rw)
		enc.Encode(&t)
	} else if req.Method == "PUT" {
		var d = json.NewDecoder(req.Body)
		var t templateResponse
		d.Decode(&t)

		if err := ioutil.WriteFile(fmt.Sprintf("%v/%v.html", TemplateDirectory, templateID), ([]byte)(t.Data), 0644); err != nil {
			log.Printf("failed: %v", err)
		}

		var e = json.NewEncoder(rw)
		e.Encode(&t)
	}
}

func generatePDF(rw http.ResponseWriter, req *http.Request) {
	var templateID = req.URL.Path

	var tmp, err = ioutil.TempFile("", "sample")
	if err != nil {
		log.Fatalf("Error creating temp file")
	}

	// cleanup tempfile
	defer func(f *os.File) {
		f.Close()
		os.Remove(f.Name())
	}(tmp)

	var data, _ = fiscalPeriodTemplate(templateID)
	var c context

	measureRuntime("umsatz api", func() {
		var fp fiscalPeriod
		if fp, err = UmsatzClient.FetchFiscalPeriod(templateID); err != nil {
			log.Fatalf("Error loading data: %v", err)
		}
		c.fiscalPeriod = fp
		c.metaData = UmsatzClient.getMetaData()
	})

	var html *template.Template
	html, err = template.New("pdf").Funcs(templateHelperFuncs).Parse(string(data))

	err = html.Execute(tmp, &c)
	if err != nil {
		log.Fatalf("Error evaluating the template: %v", err)
	}

	var fw = flushWritter{w: rw}
	if f, ok := rw.(http.Flusher); ok {
		fw.f = f
	}

	cmd := exec.Command("weasyprint", "-f", "pdf", "-e", "utf-8", tmp.Name(), "-")
	cmd.Dir = "/tmp"
	cmd.Stdout = &fw

	measureRuntime("weasyprint", func() {
		cmd.Run()
	})
}

func authHandler(next http.Handler, authAddress string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// nginx adds the X-Requires-Authentication header.
		var authenticationRequired = r.Header.Get("X-Requires-Authentication") != ""
		if !authenticationRequired {
			next.ServeHTTP(w, r)
			return
		}

		var sessionKey = r.Header.Get("X-UMSATZ-SESSION")

		var resp, err = http.Get(authAddress + "/validate/" + sessionKey)
		if err != nil {
			log.Printf("err: %v", err)
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

func newServer() http.Handler {
	var r = http.DefaultServeMux

	r.Handle("/generate/", http.StripPrefix("/generate/", http.HandlerFunc(generatePDF)))
	r.Handle("/preview/", http.StripPrefix("/preview/", http.HandlerFunc(previewPDF)))
	r.Handle("/template/", http.StripPrefix("/template/", http.HandlerFunc(templateHandler)))

	return http.Handler(logHandler(r))
}

// Set by make file on build
var (
	Version           string
	Commit            string
	UmsatzClient      endpoint
	TemplateDirectory string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		httpAddress       = flag.String("http.addr", ":8080", "HTTP listen address")
		umsatzAddress     = flag.String("umsatz.addr", "http://127.0.0.1:8081", "HTTP listen address of umsatz api")
		authAddress       = flag.String("auth.addr", "http://127.0.0.1:8084", "HTTP listen address of umsatz auth api")
		templateDirectory = flag.String("templates", "", "template directory")
		printVersion      = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	UmsatzClient = newEndpoint(*umsatzAddress, *authAddress)

	if *printVersion {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	if *templateDirectory == "" {
		*templateDirectory = "./templates"
	}
	TemplateDirectory = *templateDirectory

	log.Printf("listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, authHandler(newServer(), *authAddress)))
}
