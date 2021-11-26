package main

import (
	"log"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	uri := "tcp://47.100.222.170:10004/yixiang?username=&compress=true&debug=true"

	var err error
	db, err = gorm.Open(clickhouse.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database, got error %v", err)
	}
}
