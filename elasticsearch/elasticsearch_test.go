// Package elasticsearch provides ...
package elasticsearch

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestAddDoc(t *testing.T) {
	doc := map[string]interface{}{
		"created_at": time.Now().Format("2006-01-02"),
		"id":         10000,
		"price":      1 * 3,
		"tag":        "book",
		"title":      "book_" + strconv.Itoa(1),
	}

	jsonDoc, _ := json.Marshal(doc)
	res, err := CreateDoc(strconv.Itoa(10000), string(jsonDoc))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)
	t.Log("add doc test ok!!!")
}

func TestSearchDoc(t *testing.T) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "book_10000",
			},
		},
	}
	jsonDoc, _ := json.Marshal(query)
	res, err := Search(string(jsonDoc))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)
	t.Log("search doc test ok!!!")
}

func TestGetDoc(t *testing.T) {
	id := "10000"
	res, err := GetDoc(id)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)
	t.Log("get doc test ok!!!")
}

func TestDelDoc(t *testing.T) {
	res, err := DeleteDoc(strconv.Itoa(10000))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(res)
	t.Log("del doc test ok!!!")
}
