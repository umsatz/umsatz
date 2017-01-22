package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

func newLink(h *http.Header, rel string, href string, method string) link {
	baseUri := h.Get("X-Requested-Uri")
	if strings.HasSuffix(baseUri, "/") {
		baseUri = baseUri[:len(baseUri)-1]
	}
	link := link{
		rel,
		baseUri + href,
		method,
	}
	return link
}

type backupEntry struct {
	Backup
	Links []link `json:"_links"`
}

type Backup struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Path      string `json:"-"`
}

func (b *Backup) BaseName() string {
	index := strings.LastIndex(b.Name, ".")
	return b.Name[:index]
}

type Provider interface {
	List() ([]*Backup, error)
}

type BackupRestorer interface {
	Restore(*Backup) error
}

type AnsibleBackupRestorer struct {
	AnsiblePlaybookPath     string // path to ansible-playbook
	ProvisioningDirectory   string // provisioning directory
	BackupConfigurationPath string // path to backup.json
}

func (r *AnsibleBackupRestorer) Restore(b *Backup) error {
	cmd := exec.Command(r.AnsiblePlaybookPath, "-i", "hosts", "restore.yml", "-e", "@"+r.BackupConfigurationPath, "-e", "archive="+b.Path)
	cmd.Dir = r.ProvisioningDirectory

	fmt.Printf("restoring…\n")

	output, err := cmd.Output()
	fmt.Printf("error: %v\nansible: %v", err, string(output))
	return nil
}

type FileSystemProvider struct {
	BackupDirectory string
}

func (fs *FileSystemProvider) List() ([]*Backup, error) {
	backups := make([]*Backup, 0)

	files, err := ioutil.ReadDir(fs.BackupDirectory)
	if err != nil {
		fmt.Printf("unable to list: %v", err)
		return nil, nil
	}

	for i, file := range files {
		if filepath.Ext(file.Name()) != ".tar" {
			continue
		}

		backups = append(backups, &Backup{
			Id:        i,
			Name:      file.Name(),
			CreatedAt: strings.Split(file.Name()[7:], ".")[0],
			Path:      fs.BackupDirectory + "/" + file.Name(),
		})
	}

	return backups, nil
}

func newBackupProvider(baseDir string) (Provider, error) {
	return &FileSystemProvider{baseDir}, nil
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func listBackups(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		enc := json.NewEncoder(w)
		if req.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"method": "mismatch"})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		var backups []*Backup
		var err error
		backups, err = p.List()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"backup": "list error"})
			return
		}

		var entries []backupEntry = make([]backupEntry, len(backups))
		for i, backup := range backups {
			entries[i] = backupEntry{
				*backups[i],
				[]link{
					newLink(&req.Header, "restore", fmt.Sprintf("/%v/restore", backup.BaseName()), "POST"),
					newLink(&req.Header, "delete", fmt.Sprintf("/%v", backup.BaseName()), "DELETE"),
				},
			}
		}

		var bytes []byte
		bytes, err = json.Marshal(entries)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"backup": "serialization error"})
			return
		}

		io.WriteString(w, string(bytes))
	})
}

func deleteBackup(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		enc := json.NewEncoder(w)

		if req.Method != "DELETE" {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"method": "mismatch"})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		parts := strings.Split(req.URL.Path, "/")

		if len(parts) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"parts": "too many"})
			return
		}

		var backups []*Backup
		var err error
		backups, err = p.List()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"backups": "list error"})
			return
		}

		backupName := parts[1]

		for _, backup := range backups {
			if backup.BaseName() == backupName {
				os.Remove(backup.Path)

				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"backup": "unknown"})
	})
}

func restoreBackup(p Provider, restorer BackupRestorer) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		enc := json.NewEncoder(w)

		if req.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"method": "mismatch"})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		parts := strings.Split(req.URL.Path, "/")

		if len(parts) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"parts": "too many"})
			return
		}

		var backups []*Backup
		var err error
		backups, err = p.List()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(map[string]string{"backups": "list error"})
			return
		}

		backupName := parts[1]

		for _, backup := range backups {
			if backup.BaseName() == backupName {
				fmt.Printf("starting restore\n")
				restorer.Restore(backup)
				fmt.Printf("restore done\n")

				w.WriteHeader(http.StatusAccepted)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"backup": "unknown"})
	}
	return http.HandlerFunc(handler)
}

func newBackupServer(provider Provider, restorer BackupRestorer) http.Handler {
	r := http.DefaultServeMux

	r.HandleFunc("/", logHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/restore") {
			restoreBackup(provider, restorer).ServeHTTP(w, req)
		} else {
			if req.Method == "GET" {
				listBackups(provider).ServeHTTP(w, req)
			} else if req.Method == "DELETE" {
				deleteBackup(provider).ServeHTTP(w, req)
			}
		}
	})))

	return http.Handler(r)
}

// Set by make file on build
var (
	Version string
	Commit  string
)

func main() {
	var (
		backupRoot            = flag.String("backup.root", "/tmp", "FileSystem backup root directory")
		httpAddress           = flag.String("http.addr", ":8080", "HTTP listen address")
		backupConfig          = flag.String("backup.config", "", "Path to backup configuration file")
		ansiblePlaybook       = flag.String("ansible.playbook", "", "Path to ansible-playbook executable")
		provisioningDirectory = flag.String("provisioning.directory", "", "Path to provisioning directory")
		printVersion          = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	p, err := newBackupProvider(*backupRoot)
	if err != nil {
		log.Fatal(err)
	}

	ansibleRestorer := AnsibleBackupRestorer{
		AnsiblePlaybookPath:     *ansiblePlaybook,
		ProvisioningDirectory:   *provisioningDirectory,
		BackupConfigurationPath: *backupConfig,
	}

	log.Printf("listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, newBackupServer(p, &ansibleRestorer)))
}
