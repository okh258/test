package test

import (
	"context"
	"test/util"
	"testing"

	"git.devops.com/go/golib/odm"
	"git.devops.com/wsim/hflib/dao"
)

func InitDDB(t *testing.T) {
	db, err := odm.Open("dynamodb", "AccessKey=123;SecretKey=456;Token=789;Region=localhost;Endpoint=http://127.0.0.1:8000")
	//db, err := odm.Open("dynamodb", "AccessKey=123;SecretKey=456;Token=789;Region=cn-northwest-1;Endpoint=http://dynamodb.cn-northwest-1.amazonaws.com")
	if err != nil {
		t.Fatalf("connect dynamodb failed, err: %v", err)
	}
	dao.SetDB(db)
}

func TestGetRealNameAuthCount(t *testing.T) {
	InitDDB(t)
	count, err := GetRealNameAuthCount(t, context.Background(), 1, 0, util.MicroTime())
	if err != nil {
		t.Logf("GetSkillCertCount failed, err: %v", err)
		return
	}
	t.Logf("count: %v", count)
}

func GetRealNameAuthCount(t *testing.T, ctx context.Context, status, start, end int64) (int64, error) {
	count, hasMore, err := dao.DB().Query(ctx).Table("realname_auth").
		Index("status-create_time-index").
		HashKey("status").Hash(status).
		RangeKey("create_time").Between(start, end).
		Count()
	if err != nil {
		t.Errorf("GetRealNameAuthCount: %v", err.Error())
		return 0, err
	}
	t.Logf("count: %v, hasMore: %v", count, hasMore)
	return count, nil
}

// GetCensorVideoCount 获取需要审核的视频墙
func GetCensorVideoCount(t *testing.T, ctx context.Context, status, start, end int64) (int64, error) {
	count, hasMore, err := dao.DB().Query(ctx).Table("user_video").
		Index("video_status-create_time-index").
		HashKey("video_status").Hash(status).
		RangeKey("create_time").Between(start, end).
		Count()
	if err != nil {
		t.Errorf("GetCensorVideoCount: %v", err.Error())
		return 0, err
	}
	t.Logf("count: %v, hasMore: %v", count, hasMore)
	return count, nil
}
