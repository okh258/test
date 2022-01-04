package services

import (
	"context"
	"fmt"
	"git.devops.com/go/golib/odm"
	"git.devops.com/go/golib/odm/dynamo"
	"git.devops.com/wsim/hflib/dao"
	"git.devops.com/wsim/hflib/logs"
	"os"
	"test/model"
	"test/models"
)

func init() {
	uri := os.Getenv("ddb_uri")
	initDb("uat", "dynamodb", uri)
}

func initDb(env string, dbtype string, dbpath string) {
	db, err := odm.Open(dbtype, dbpath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if db, ok := db.(*odm.ODMDB).DialectDB.(*dynamo.DB); ok {
		//db.SetTablePrefix(strings.ToLower(env) + "_")
		db.SetTablePrefix(fmt.Sprintf("%s_%s_", env, "tm"))
	}
	dao.SetDB(db)
}

func GetUserDriverInfo(ctx context.Context, uid int64) (*model.UserDevice, error) {
	var result []*model.UserDevice
	builder := dao.DB().Query(ctx).Table("user_device")
	builder.Index("last_user_id_index").HashKey("last_user_id").Hash(uid)
	_, err := builder.Limit(1).List(&result)
	if err != nil || len(result) < 1 {
		return nil, err
	}
	return result[0], err
}

func GetUserDriverMap(ctx context.Context, uid int64) (*map[string]interface{}, error) {
	var result []*map[string]interface{}
	builder := dao.DB().Query(ctx).Table("user_device")
	builder.Index("last_user_id_index").HashKey("last_user_id").Hash(uid)
	_, err := builder.List(&result)
	if err != nil || len(result) < 1 {
		return nil, err
	}
	return result[0], err
}

func GetUserDriver(ctx context.Context, deviceToken string) (*model.UserDevice, error) {
	var r model.UserDevice
	return &r, dao.DB().Get(ctx, &r, deviceToken)
}

func GetUserDriverItem(ctx context.Context, uid int64) (*model.UserDevice, error) {
	var result model.UserDevice
	err := dao.DB().Table("user_device").GetItem(ctx, "last_user_id", nil, nil, &result)
	if err != nil {
		logs.Errorf(ctx, "DDB GetUserInfo err:%v", err.Error())
		return nil, err
	}
	return &result, nil
}

func GetDetailUsers(ctx context.Context, uids []interface{}) (map[int64]*models.UserBaseInfo, error) {
	var userInfos []*models.UserBaseInfo

	opt := &odm.BatchGet{TableName: "user_user"}
	for _, uid := range uids {
		opt.Keys = append(opt.Keys, odm.Key{
			HashKey:  uid,
			RangeKey: nil,
		})
	}
	err := dao.DB().Table(models.UserBaseInfo{}).GetDB().BatchGetItem(ctx, []*odm.BatchGet{opt}, nil, &userInfos)
	if err != nil {
		logs.Errorf(ctx, "DDB GetDetailUsers err: %v", err)
		return nil, err
	}
	userInfoMap := make(map[int64]*models.UserBaseInfo)
	for _, info := range userInfos {
		userInfoMap[info.UserId] = info
	}
	return userInfoMap, nil
}
