package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/zhenligod/my_elasticsearch/elasticsearch"
)

func main() {
	body := `{
		"sort": { "date": "desc"} 
	}`
	res, err := elasticsearch.Search(body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)

	doc := map[string]string{
		"date": time.Now().Format("2006-01-02"),
	}

	jsonDoc, _ := json.Marshal(doc)
	res, err = elasticsearch.CreateDoc("5", string(jsonDoc))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)
}
