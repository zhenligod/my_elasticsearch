package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v6"
	"github.com/robfig/config"
	log "github.com/sirupsen/logrus"
)

var es *elasticsearch.Client
var esAddr string = "http://ip:port" // es 地址及端口
var esIndex string = "my_index"      // index 前缀，表示获取job*的所有index
var configFile = flag.String("configfile", "conf/es.conf", "General configuration file")

// EsConf es配置
type EsConf struct {
	IP    string
	Port  string
	Index string
	Es    *elasticsearch.Client
}

func init() {
	var err error
	conf := make(map[string]string)
	cfg, err := config.ReadDefault(*configFile) //读取配置文件，并返回其Config

	if err != nil {
		log.Fatalf("Fail to find %v,%v", *configFile, err)
	}
	if cfg.HasSection("es") { //判断配置文件中是否有section（一级标签）
		options, err := cfg.SectionOptions("es") //获取一级标签的所有子标签options（只有标签没有值）
		if err == nil {
			for _, v := range options {
				optionValue, err := cfg.String("es", v) //根据一级标签section和option获取对应的值
				if err == nil {
					conf[v] = optionValue
				}
			}
		}
	}
	log.Println(conf)
	config := elasticsearch.Config{}
	esConf := EsConf{}
	esConf.IP = conf["es_hostname"]
	esConf.Port = conf["es_port"]
	esConf.Index = conf["es_index"]
	address := strings.Join([]string{
		"http",
		"//" + esConf.IP,
		esConf.Port,
	}, ":")
	esIndex = esConf.Index
	config.Addresses = []string{address}
	es, err = elasticsearch.NewClient(config)
	if err != nil {
		log.Error(err.Error())
	}
}

// SearchByJob 通过job名称获取es日志
func SearchByJob(job string) (*string, error) {
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
