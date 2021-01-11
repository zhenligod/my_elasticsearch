package main

import (
	"log"

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
}
