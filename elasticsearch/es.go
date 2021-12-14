package elasticsearch

import (
	"fmt"
	"github.com/olivere/elastic/v7"
)

func ES() *elastic.Client {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return elasticClient
}

func GetIndexName(indexName string) string {
	return fmt.Sprintf("uat_tm_%v", indexName)
}
