package service

import (
	"fmt"
	"git.devops.com/go/golib/hfmysql"
	"os"
)

func init() {
	uri := os.Getenv("mysql_uri")
	uri = fmt.Sprintf("%s%s", uri, "tm_dating")
	err := hfmysql.Init(10, 5, uri)
	if err != nil {
		fmt.Printf("connect mysql failed, err: %v\n", err)
		os.Exit(1)
	}
}
