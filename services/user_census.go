package services

import (
	"context"
	"fmt"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	"github.com/olivere/elastic/v7"
	"test/db"
	"test/elasticsearch"
	"test/models"
	"time"
)

var _userCensusService *UserCensusService

func NewUserCensusService() *UserCensusService {
	if _userCensusService == nil {
		_userCensusService = &UserCensusService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameAdmin),
		}
		orm.RegisterModel(&models.UserCountCensus{}, &models.VipCensusLog{})
	}
	return _userCensusService
}

type UserCensusService struct {
	o orm.Ormer
}

func (s *UserCensusService) GetUserCountBy(params models.SearchUserRequestParams) int64 {
	esClient := elasticsearch.ES()
	indexName := elasticsearch.GetIndexName(models.IndexNameUser)
	query := elastic.NewBoolQuery()
	if params.Gender > 0 {
		query.Must(elastic.NewTermQuery("gender", params.Gender))
	}
	if params.IsCoach != 0 {
		if params.IsCoach == -1 {
			query.Must(elastic.NewTermQuery("user_identity", 0))
		} else {
			query.Must(elastic.NewTermQuery("user_identity", 1))
		}
	}
	if params.OsType != 0 {
		query.Must(elastic.NewTermQuery("os_type", params.OsType))
	}
	if params.EndTime != 0 && params.StartTime == 0 {
		query.Must(elastic.NewRangeQuery("register_time").Lte(params.EndTime))
	} else if params.EndTime != 0 && params.StartTime != 0 {
		query.Must(elastic.NewRangeQuery("register_time").Gte(params.StartTime).Lte(params.EndTime))
	}

	resp, err := esClient.Search().Index(indexName).
		Query(query).
		Sort("register_time", false).
		Size(0).
		Do(context.TODO())
	if err != nil {
		logs.Errorf(context.TODO(), "GetUserCountBy failed, err: %v", err)
		return 0
	}

	return resp.TotalHits()
}

// GetUserCensusCount 根据时间查找用户统计数量
// start end 任意不传, 查询当天统计信息
func (s *UserCensusService) GetUserCensusCount(ctx context.Context, start, end *time.Time) (*models.UserCountCensus, error) {
	if start == nil {
		return nil, fmt.Errorf("start time can't be nil")
	}
	o := s.o
	var result *models.UserCountCensus
	sql := "SELECT * FROM t_user_count_census WHERE TO_DAYS(create_date) = TO_DAYS(?) LIMIT 1"
	r := o.Raw(sql, start)
	if start != nil && end != nil {
		sql = "SELECT %s%s%s%s FROM t_user_count_census WHERE TO_DAYS(create_date) >= TO_DAYS(?) AND TO_DAYS(create_date) <= TO_DAYS(?)"
		c1 := "SUM(verify_count) verify_count, SUM(register_count) register_count, SUM(register_male_count) register_male_count, SUM(register_female_count) register_female_count, "
		c2 := "SUM(coach_count) coach_count, SUM(coach_male_count) coach_male_count, SUM(coach_female_count) coach_female_count, "
		c3 := "SUM(consumer_count) consumer_count, SUM(consumer_male_count) consumer_male_count, SUM(consumer_female_count) consumer_female_count, SUM(ios_count) ios_count, SUM(android_count) android_count, "
		c4 := "MAX(verify_total) verify_total, MAX(register_total) register_total, MAX(coach_total) coach_total, MAX(consumer_total) consumer_total, MAX(ios_total) ios_total, MAX(android_total) android_total, MAX(create_date) create_date, MAX(create_time) create_time, MAX(update_time) update_time "
		sql = fmt.Sprintf(sql, c1, c2, c3, c4)
		r = o.Raw(sql, start, end)
	}
	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "GetUserCensus failed, err: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *UserCensusService) UpsertUserCensus(ctx context.Context, data *models.UserCountCensus) error {
	o := s.o
	var count int64
	sql := "SELECT count(1) FROM t_user_count_census WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	err := o.Raw(sql, data.CreateDate).QueryRow(&count)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "UserCountCensus.get failed, err: %v", err.Error())
		return err
	}
	sql = "UPDATE t_user_count_census SET `verify_count` = ?, `register_count` = ?, `register_male_count` = ?, `register_female_count` = ?, `coach_count` = ?, `coach_male_count` = ?, `coach_female_count` = ?, `consumer_count` = ?, `consumer_male_count` = ?, `consumer_female_count` = ?, `ios_count` = ?, `android_count` = ?, `verify_total` = ?, `register_total` = ?, `coach_total` = ?, `consumer_total` = ?, `ios_total` = ?, `android_total` = ?, `create_time` = ?, `update_time` = ? WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	if count < 1 {
		sql = "INSERT INTO `t_user_count_census` (`verify_count`, `register_count`, `register_male_count`, `register_female_count`, `coach_count`, `coach_male_count`, `coach_female_count`, `consumer_count`, `consumer_male_count`, `consumer_female_count`, `ios_count`, `android_count`, `verify_total`, `register_total`, `coach_total`, `consumer_total`, `ios_total`, `android_total`, `create_time`, `update_time`, `create_date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
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

func (s *UserCensusService) GetVipCensusLog(ctx context.Context, log *models.VipCensusLog) error {
	_, err := s.o.Insert(log)
	if err != nil {
		logs.Errorf(ctx, "AddVipCensusLog err:", err)
		return err
	}
	return nil
}

// GetVipCensusCount 获取VIP统计信息
func (s *UserCensusService) GetVipCensusCount(start, end int64) (*models.VipCensusView, error) {
	var result models.VipCensusView
	o := s.o
	column := fmt.Sprintf("COUNT(DISTINCT uid) vip_total, AVG(duration) vip_avg, COUNT(1) vip_count, SUM(IF(duration=%v,1,0)) one_month_count, SUM(IF(duration=%v,1,0)) three_month_count, SUM(IF(duration=%v,1,0)) one_year_count, SUM(IF(duration=%v,1,0)) two_year_count", models.VipTypeOneMonth, models.VipTypeThreeMonth, models.VipTypeOneYear, models.VipTypeTwoYear)
	sql := "SELECT %s FROM t_vip_census_log WHERE create_time >= ? AND create_time <= ?"
	err := o.Raw(fmt.Sprintf(sql, column), start, end).QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(context.TODO(), "GetVipCensusCount: %v", err.Error())
		return &result, err
	}
	return &result, nil
}
