package test

import (
	"context"
	"os"
	"test/util"
	"testing"

	"git.devops.com/go/golib/hfmysql"
	"github.com/astaxie/beego/orm"
)

func InitMysql(t *testing.T) {
	uri := os.Getenv("mysql_uri")
	err := hfmysql.Init(10, 5, uri+"/tm_dating")
	if err != nil {
		t.Fatalf("connect mysql failed, err: %v", err)
	}
}

func TestGetSkillCertCount(t *testing.T) {
	InitMysql(t)
	count, err := GetSkillCertCount(t, context.Background(), 1, 0, util.MicroTime())
	if err != nil {
		t.Logf("GetSkillCertCount failed, err: %v", err)
		return
	}
	t.Logf("count: %v", count)
}

func GetSkillCertCount(t *testing.T, ctx context.Context, status, start, end int64) (int64, error) {
	var resultCount int64
	o := orm.NewOrm()
	sql := "SELECT COUNT(1) FROM t_skill_certification WHERE `status` = ? AND create_time >= ? AND create_time <= ?"
	err := o.Raw(sql, status, start, end).QueryRow(&resultCount)
	if err != nil && err != orm.ErrNoRows {
		t.Errorf("GetSkillCertCount: %v", err.Error())
		return 0, err
	}
	return resultCount, nil
}

func TestGetCensorUserInfoCount(t *testing.T) {
	InitMysql(t)
	count1, _ := GetCensorUserInfoCount(t, context.Background(), 0, 1, 0, util.MicroTime())
	count2, _ := GetCensorUserInfoCount(t, context.Background(), 0, 2, 0, util.MicroTime())
	count3, _ := GetCensorUserInfoCount(t, context.Background(), 0, 3, 0, util.MicroTime())
	t.Logf("count1: %v", count1)
	t.Logf("count2: %v", count2)
	t.Logf("count3: %v", count3)
}

// GetCensorUserInfoCount 获取需要审核的用户信息
// status 0待审核 1通过 2拒绝
// contentType 1昵称 2头像 3个签
func GetCensorUserInfoCount(t *testing.T, ctx context.Context, status, contentType, start, end int64) (int64, error) {
	var resultCount int64
	o := orm.NewOrm()
	sql := "SELECT COUNT(1) FROM t_userinfo_check WHERE `status` = ? AND content_type = ? AND create_time >= ? AND create_time <= ?"
	err := o.Raw(sql, status, contentType, start, end).QueryRow(&resultCount)
	if err != nil && err != orm.ErrNoRows {
		t.Errorf("GetSkillCertCount: %v", err.Error())
		return 0, err
	}
	return resultCount, nil
}

func TestGetFeedbackCount(t *testing.T) {
	InitMysql(t)
	count1 := GetFeedbackCount(t, context.Background(), 0, util.MicroTime())
	count2 := GetTipOffCount(t, context.Background(), 0, util.MicroTime())
	t.Logf("count1: %v", count1)
	t.Logf("count2: %v", count2)
}

// GetFeedbackCount 获取意见反馈数量
func GetFeedbackCount(t *testing.T, ctx context.Context, start, end int64) int64 {
	var resultCount int64
	o := orm.NewOrm()
	sql := "SELECT COUNT(1) FROM feedback WHERE create_time >= ? AND create_time <= ?"
	err := o.Raw(sql, start, end).QueryRow(&resultCount)
	if err != nil && err != orm.ErrNoRows {
		t.Errorf("GetFeedbackCount: %v", err.Error())
		return 0
	}
	return resultCount
}

// GetTipOffCount 获取投诉举报数量
func GetTipOffCount(t *testing.T, ctx context.Context, start, end int64) int64 {
	var resultCount int64
	o := orm.NewOrm()
	sql := "SELECT COUNT(1) FROM tip_off WHERE create_time >= ? AND create_time <= ?"
	err := o.Raw(sql, start, end).QueryRow(&resultCount)
	if err != nil && err != orm.ErrNoRows {
		t.Errorf("GetTipOffCount: %v", err.Error())
		return 0
	}
	return resultCount
}
