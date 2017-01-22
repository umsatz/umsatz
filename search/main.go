// TODO bleve based full text search.
// whenever the data of a position changes, we receive an NSQ event, and request
// the content of the attachment, if any.
// we then insert basic position informations alongside any conversion informations
// based on the position id
// also we provide a http search endpoint
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/blevesearch/bleve"
)

type positionEvent struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type position struct {
	ID         string
	Attachment string
	// TODO add other fields fetched via api
}

func (p *positionEvent) positionID() string {
	var ps = strings.Split(p.URL, "/")
	return ps[len(ps)-1]
}

func main() {
	var config = nsq.NewConfig()
	var consumer, err = nsq.NewConsumer("positions", "ch", config)

	if err != nil {
		log.Fatalf("%v", err)
	}

	mapping := bleve.NewIndexMapping()
	// TODO store needs to be configurable to be saved between deployments
	index, err := bleve.New("positions.bleve", mapping)

	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		var buf = bytes.NewBuffer(m.Body)
		var d = json.NewDecoder(buf)
		var e positionEvent
		d.Decode(&e)

		log.Printf("%v", e.URL)

		// var resp, err = http.Get(e.URL)
		// if err != nil {
		// 	log.Printf("%v", resp.Body)
		// } else {
		// 	log.Fatalf("%v", err)
		// }

		var p = position{
			ID: e.positionID(),
		}
		if e.Type == "create" {
			index.Index(e.positionID(), p)
		} else if e.Type == "update" {
			index.Index(e.positionID(), p)
		} else if e.Type == "destroy" {
			index.Delete(e.positionID())
		}
		m.Finish()
		return nil
	}))

	if err = consumer.ConnectToNSQDs([]string{"127.0.0.1:4150"}); err != nil {
		log.Fatalf("unable to connect to nsqd: %v", err)
	}

	for {
		select {
		case <-time.After(time.Second * 10):
		}
	}
}
