package elasticsearch

import (
	"flag"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/robfig/config"
)

var (
	es         *elasticsearch.Client
	esAddr     string = "http://localhost:9200" // es 地址及端口
	esIndex    string = "my_index"              // index 前缀
	configFile        = flag.String("configfile", "../conf/es.conf", "General configuration file")
)

// EsConf es配置
type EsConf struct {
	IP       string
	Port     string
	Index    string
	Username string
	Passwd   string
	Es       *elasticsearch.Client
}

func init() {
	var err error
	conf := make(map[string]string)
	cfg, err := config.ReadDefault(*configFile) //读取配置文件，并返回其Config

	if err != nil {
		log.Fatalf("Fail to find %v,%v", *configFile, err)
	}
	if cfg.HasSection("es_goods") { //判断配置文件中是否有section（一级标签）
		options, err := cfg.SectionOptions("es_goods") //获取一级标签的所有子标签options（只有标签没有值）
		if err == nil {
			for _, v := range options {
				optionValue, err := cfg.String("es_goods", v) //根据一级标签section和option获取对应的值
				if err == nil {
					conf[v] = optionValue
				}
			}
		}
	}
	config := elasticsearch.Config{}
	esConf := EsConf{}
	esConf.IP = conf["es_hostname"]
	esConf.Port = conf["es_port"]
	esConf.Index = conf["es_index"]
	esConf.Username = conf["es_username"]
	esConf.Passwd = conf["es_passwd"]
	address := strings.Join([]string{
		"http",
		"//" + esConf.IP,
		esConf.Port,
	}, ":")
	esIndex = esConf.Index
	config.Addresses = []string{address}
	config.Username = esConf.Username
	config.Password = esConf.Passwd
	es, err = elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// Search search docs
func Search(body string) (*esapi.Response, error) {
	res, err := es.Search(
		es.Search.WithIndex(esIndex),
		es.Search.WithBody(strings.NewReader(body)),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}

// CreateDoc add doc with es client
func CreateDoc(id string, body string) (*esapi.Response, error) {
	res, err := es.Create(esIndex, id, strings.NewReader(body))

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}

// UpdateDoc update doc with es client
func UpdateDoc(id string, body string) (*esapi.Response, error) {
	res, err := es.Update(esIndex, id, strings.NewReader(body))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}

// DeleteDoc delete doc with es client
func DeleteDoc(id string) (*esapi.Response, error) {
	res, err := es.Delete(esIndex, id)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}

// GetDoc get sigle doc from es
func GetDoc(id string) (*esapi.Response, error) {
	res, err := es.Get(esIndex, id)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}

// SQLDoc exec sql query
func SQLDoc(body string) (*esapi.Response, error) {
	res, err := es.SQL.Query(strings.NewReader(body))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}

	return res, nil
}
