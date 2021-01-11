package main

import (
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
)

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://9.135.127.17:9200"},
		Transport: nil,
		Username:  "elastic",
		Password:  "paNBq*2693br",
	})

	if err != nil {
		log.Fatalf("Error getting es: %s", err)
	}

	body := `{
		"sort": { "date": "desc"} 
	}`

	res, err := es.Search(
		es.Search.WithIndex("my_index"),
		es.Search.WithBody(strings.NewReader(body)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)

	defer res.Body.Close()

}
