package services

import (
	"context"
	"test/models"
	"test/util"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
)

func TestGetSkillCertCount(t *testing.T) {
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

func TestGetUserCensus(t *testing.T) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	count1, err := NewAService().GetUserCensusCount(context.Background(), &start, &now)
	if err != nil {
		t.Fatalf("get failed, err: %v", err)
	}
	t.Logf("count1: %+v", count1)
}

func TestGetUserCensusCount(t *testing.T) {
	s, _ := time.Parse("2006-01-02 15:04:05", "2021-12-05 00:00:00")
	e, _ := time.Parse("2006-01-02 15:04:05", "2021-12-06 00:00:00")
	r, err := NewAService().GetUserCensusCount(context.TODO(), &s, &e)
	if err != nil {
		t.Fatalf("get error, err: %v", err)
		return
	}
	t.Logf("census: %+v", r)
}

func TestUpsert(t *testing.T) {
	now := time.Now()
	d := &models.UserCountCensus{
		CreateDate:          now,
		RegisterCount:       2,
		RegisterMaleCount:   1,
		RegisterFemaleCount: 1,
		CoachCount:          1,
		CoachMaleCount:      1,
		CoachFemaleCount:    1,
		ConsumerCount:       1,
		ConsumerMaleCount:   1,
		ConsumerFemaleCount: 1,
		RegisterTotal:       1,
		CoachTotal:          1,
		ConsumerTotal:       1,
		CreateTime:          util.Timestamp(),
		UpdateTime:          util.Timestamp(),
	}
	err := NewAService().UpsertUserCensus(context.TODO(), d)
	if err != nil {
		t.Fatalf("upsert failed, err: %v", err)
	}
	t.Logf("insert ok...\n")
}

func TestAddVipCensusLog(t *testing.T) {
	log := &models.VipCensusLog{
		OrderId:    11111111,
		Uid:        22222220,
		Duration:   12,
		CreateTime: util.Timestamp(),
	}
	err := NewAService().AddVipCensusLog(context.TODO(), log)
	if err != nil {
		t.Logf("add failed, err: %v", err)
		return
	}
	t.Logf("add log ok. \n")
	count, err := GetSkillCertCount(t, context.Background(), 1, 0, util.MicroTime())
	if err != nil {
		t.Logf("GetSkillCertCount failed, err: %v", err)
		return
	}
	t.Logf("count: %v", count)
}

func TestGetInviteeUserCount(t *testing.T) {
	data, err := NewAService().GetInviteeUserCount(0, util.Timestamp())
	if err != nil {
		t.Fatalf("get error, err: %v\n", err)
	}
	t.Logf("data: %+v\n", data)
}

func TestGetDynamicCensusCount(t *testing.T) {
	s, _ := time.Parse("2006-01-02 15:04:05", "2021-12-01 00:00:00")
	e, _ := time.Parse("2006-01-02 15:04:05", "2021-12-05 00:00:00")

	data, err := NewAService().GetDynamicCensusCount(context.TODO(), &s, &e)
	if err != nil {
		t.Fatalf("get error, err: %v\n", err)
	}
	t.Logf("data: %+v\n", data)
}