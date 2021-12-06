package db

import (
	"context"
	"database/sql"
	"fmt"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

const (
	AliasNameAdmin  = "tm_admin"
	AliasNameDating = "tm_dating"
)

const (
	maxIdleConn = 10
	maxOpenConn = 5
	maxLifeTime = time.Hour * 7
)

func RegisterDataBase(ctx context.Context, aliasName string) (*sql.DB, error) {
	uri := os.Getenv("mysql_uri")
	uri = fmt.Sprintf("%s%s", uri, aliasName)
	err := orm.RegisterDataBase(aliasName, "mysql", uri, maxIdleConn, maxOpenConn)
	if err != nil {
		return nil, err
	}

	DB, err := orm.GetDB(aliasName)
	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(maxLifeTime)

	return DB, nil
}

func NewOrmWithDB(ctx context.Context, aliasName string) orm.Ormer {
	DB, err := orm.GetDB(aliasName)
	if err != nil {
		DB, err = RegisterDataBase(ctx, aliasName)
		if err != nil {
			logs.Error(ctx, fmt.Sprintf("orm db alias %s have not register, %s", aliasName, err))
			panic(fmt.Sprintf("orm db alias %s have not register, %s", aliasName, err))
		}
	}

	ormer, _ := orm.NewOrmWithDB("mysql", aliasName, DB)
	return ormer
}
