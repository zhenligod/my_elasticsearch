package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v6"
	log "github.com/sirupsen/logrus"
)

var es *elasticsearch.Client
var esAddr string = "http://ip:port" // es 地址及端口
var esIndex string = "job*"          // index 前缀，表示获取job*的所有index

// EsConf es配置
type EsConf struct {
	IP    string
	Port  string
	Index string
	Es    *elasticsearch.Client
}

func init() {
	var err error
	config := elasticsearch.Config{}
	config.Addresses = []string{esAddr}
	es, err = elasticsearch.NewClient(config)
	if err != nil {
		log.Error(err.Error())
	}
}

// SearchByJob 通过job名称获取es日志
func SearchByJob(job string, conf EsConf) (*string, error) {
	var (
		buf bytes.Buffer
		r   map[string]interface{}
		bt  bytes.Buffer
	)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"kubernetes.labels.tf-job-name": job, //具体的查询条件，可以根据日志格式进行修改
			},
		},
		"_source": map[string]interface{}{
			"includes": []interface{}{
				"log",
			},
		},
		"sort": map[string]interface{}{ // 排序
			"@timestamp": map[string]interface{}{
				"order": "asc",
			},
			"_id": map[string]interface{}{
				"order": "asc",
			},
		},
	}
	//fmt.Println(query)
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Errorf("Error encoding query: %s", err)
		return nil, err
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(esIndex),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
		es.Search.WithSize(100000),
	)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	fmt.Println(res)
	if res.IsError() {
		log.Error(res)
		return nil, errors.New(fmt.Sprint(res))
	}
	if res.StatusCode != 200 {
		log.Errorf("request error: %d", res.StatusCode)
		return nil, err
	}
	//fmt.Println(res)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}
	log.Infof("log length: %d", len(r["hits"].(map[string]interface{})["hits"].([]interface{})))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		bt.WriteString(fmt.Sprintf("%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["log"]))
	}
	logs := bt.String()
	return &logs, nil
}
