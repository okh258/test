package service

import (
	"context"
	"fmt"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	"test/db"
	"test/models"
	"time"
)

var _aService *AService

func NewAService() *AService {
	if _aService == nil {
		_aService = &AService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameAdmin),
		}
		orm.RegisterModel(&models.VipCensusLog{})
	}
	return _aService
}

type AService struct {
	o orm.Ormer
}

// GetUserCensusCount 根据时间查找用户统计数量
// start end 任意不传, 查询当天统计信息
func (s *AService) GetUserCensusCount(ctx context.Context, start, end *time.Time) (*models.UserCountCensus, error) {
	orm.RegisterModel(&models.UserCountCensus{})
	o := s.o
	var result *models.UserCountCensus
	sql := "SELECT * FROM t_user_count_census WHERE TO_DAYS(create_date) = TO_DAYS(CURDATE()) LIMIT 1"
	r := o.Raw(sql)
	if start != nil && end != nil {
		sql = "SELECT %s%s%s%s FROM t_user_count_census WHERE TO_DAYS(create_date) >= TO_DAYS(?) AND TO_DAYS(create_date) <= TO_DAYS(?)"
		c1 := "SUM(vcode_count) verify_count, SUM(register_count) register_count, SUM(register_male_count) register_male_count, SUM(register_female_count) register_female_count, "
		c2 := "SUM(coach_count) coach_count, SUM(coach_male_count) coach_male_count, SUM(coach_female_count) coach_female_count, "
		c3 := "SUM(consumer_count) consumer_count, SUM(consumer_male_count) consumer_male_count, SUM(consumer_female_count) consumer_female_count, SUM(ios_count) ios_count, SUM(android_count) android_count, "
		c4 := "MAX(vcode_total) as verify_total, MAX(register_total) register_total, MAX(coach_total) coach_total, MAX(consumer_total) consumer_total, MAX(ios_total) ios_total, MAX(android_total) android_total, MAX(create_date) create_date, MAX(create_time) create_time, MAX(update_time) update_time "
		sql = fmt.Sprintf(sql, c1, c2, c3, c4)
		r = o.Raw(sql).SetArgs(start, end)
	}
	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "GetUserCensus failed, err: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *AService) UpsertUserCountCensus1(ctx context.Context, data *models.UserCountCensus) error {
	orm.RegisterModel(&models.UserCountCensus{})
	o := s.o
	_, err := o.InsertOrUpdate(data)
	if err != nil {
		logs.Errorf(ctx, "UpsertUserCountCensus failed, err: %v", err)
	}
	return nil
}

func (s *AService) UpsertUserCensus(ctx context.Context, data *models.UserCountCensus) error {
	o := s.o
	var count int64
	sql := "SELECT count(1) FROM t_user_count_census WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	err := o.Raw(sql, data.CreateDate).QueryRow(&count)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "UserCountCensus.get failed, err: %v", err.Error())
		return err
	}
	sql = "UPDATE t_user_count_census SET `vcode_count` = ?, `register_count` = ?, `register_male_count` = ?, `register_female_count` = ?, `coach_count` = ?, `coach_male_count` = ?, `coach_female_count` = ?, `consumer_count` = ?, `consumer_male_count` = ?, `consumer_female_count` = ?, `ios_count` = ?, `android_count` = ?, `vcode_total` = ?, `register_total` = ?, `coach_total` = ?, `consumer_total` = ?, `ios_total` = ?, `android_total` = ?, `create_time` = ?, `update_time` = ? WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	if count < 1 {
		sql = "INSERT INTO `t_user_count_census` (`vcode_count`, `register_count`, `register_male_count`, `register_female_count`, `coach_count`, `coach_male_count`, `coach_female_count`, `consumer_count`, `consumer_male_count`, `consumer_female_count`, `ios_count`, `android_count`, `vcode_total`, `register_total`, `coach_total`, `consumer_total`, `ios_total`, `android_total`, `create_time`, `update_time`, `create_date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	}
	result, err := o.Raw(sql, data.VerifyCount, data.RegisterCount, data.RegisterMaleCount, data.RegisterFemaleCount,
		data.CoachCount, data.CoachMaleCount, data.CoachFemaleCount,
		data.ConsumerCount, data.ConsumerMaleCount, data.ConsumerFemaleCount, data.IosCount, data.AndroidCount,
		data.VerifyTotal, data.RegisterTotal, data.CoachTotal, data.ConsumerTotal, data.IosTotal, data.AndroidTotal, data.CreateTime, data.UpdateTime, data.CreateDate).Exec()
	if err != nil {
		logs.Error(ctx, "upsert UserCountCensus failed, err:", err)
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (s *AService) AddVipCensusLog(ctx context.Context, log *models.VipCensusLog) error {
	_, err := s.o.Insert(log)
	if err != nil {
		logs.Errorf(ctx, "AddVipCensusLog err:", err)
		return err
	}
	return nil
}

// GetInviteeUserCount 获取有邀请码的用户数量
func (s *AService) GetInviteeUserCount(start, end int64) (*models.VipCensusView, error) {
	var result models.VipCensusView
	o := s.o
	column := fmt.Sprintf("COUNT(DISTINCT uid) vip_total, AVG(duration) vip_avg, COUNT(1) vip_count, SUM(IF(duration=%v,1,0)) one_month_count, SUM(IF(duration=%v,1,0)) three_month_count, SUM(IF(duration=%v,1,0)) one_year_count, SUM(IF(duration=%v,1,0)) two_year_count", models.VipTypeOneMonth, models.VipTypeThreeMonth, models.VipTypeOneYear, models.VipTypeTwoYear)
	sql := "SELECT %s FROM t_vip_census_log WHERE create_time >= ? AND create_time <= ?"
	err := o.Raw(fmt.Sprintf(sql, column), start, end).QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(context.TODO(), "GetInviteeUserCount: %v", err.Error())
		return &result, err
	}
	return &result, nil
}

// GetDynamicCensusCount 根据时间查找用户统计数量
// start end 任意不传, 查询当天统计信息
func (s *AService) GetDynamicCensusCount(ctx context.Context, start, end *time.Time) (*models.DynamicCountCensus, error) {
	o := s.o
	var result *models.DynamicCountCensus
	sql := "SELECT * FROM t_dynamic_count_census WHERE TO_DAYS(create_date) = TO_DAYS(CURDATE()) LIMIT 1"
	r := o.Raw(sql)
	if start != nil && end != nil {
		sql = "SELECT %s%s%s%s%s FROM t_dynamic_count_census WHERE TO_DAYS(create_date) >= TO_DAYS(?) AND TO_DAYS(create_date) <= TO_DAYS(?)"
		c1 := "SUM(moment_count) moment_count, SUM(comment_count) comment_count, "
		c2 := "SUM(order_count) order_count, SUM(order_cancel_count) order_cancel_count, SUM(appoint_count) appoint_count, SUM(appoint_cancel_count) appoint_cancel_count, SUM(one_by_one_count) one_by_one_count, SUM(one_by_one_cancel_count) one_by_one_cancel_count, "
		c3 := "MAX(moment_total) moment_total, MAX(comment_total) comment_total, "
		c4 := "MAX(order_total) order_total, MAX(order_cancel_total) order_cancel_total, MAX(appoint_total) appoint_total, MAX(appoint_cancel_total) appoint_cancel_total, MAX(one_by_one_total) one_by_one_total, MAX(one_by_one_cancel_total) one_by_one_cancel_total, "
		c5 := "MAX(create_date) create_date, MAX(create_time) create_time, MAX(update_time) update_time "
		sql = fmt.Sprintf(sql, c1, c2, c3, c4, c5)
		r = o.Raw(sql).SetArgs(start, end)
	}
	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "GetUserCensus failed, err: %v", err)
		return nil, err
	}
	return result, nil
}
