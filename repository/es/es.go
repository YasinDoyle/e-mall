package es

import (
	"fmt"
	"log"

	elastic "github.com/elastic/go-elasticsearch"
	"github.com/sirupsen/logrus"

	conf "github.com/YasinDoyle/e-mall/config"
)

var EsClient *elastic.Client

func InitES() {
	eConfig := conf.Config.Es
	esConn := fmt.Sprintf("http://%s:%s", eConfig.EsHost, eConfig.EsPort)
	cfg := elastic.Config{
		Addresses: []string{esConn},
	}
	client, err := elastic.NewClient(cfg)
	if err != nil {
		log.Panic(err)
		return
	}
	EsClient = client
}

func EsHookLog() *ElasticHook {
	eConfig := conf.Config.Es
	hook, err := NewElasticHook(EsClient, eConfig.EsHost, logrus.DebugLevel, eConfig.EsIndex)
	if err != nil {
		log.Panic(err)
	}
	return hook
}
