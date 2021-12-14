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

var _dynamicCensusService *DynamicCensusService

func NewDynamicCensusService() *DynamicCensusService {
	if _dynamicCensusService == nil {
		_dynamicCensusService = &DynamicCensusService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameAdmin),
		}
		orm.RegisterModel(&models.OrderCountCensus{})
	}
	return _dynamicCensusService
}

type DynamicCensusService struct {
	o orm.Ormer
}

// GetDynamicCensusCount 根据时间查找用户统计数量
// start end 任意不传, 查询当天统计信息
func (s *DynamicCensusService) GetDynamicCensusCount(ctx context.Context, start, end *time.Time) (*models.OrderCountCensus, error) {
	o := s.o
	var result *models.OrderCountCensus
	sql := "SELECT * FROM t_order_count_census WHERE TO_DAYS(create_date) = TO_DAYS(CURDATE()) LIMIT 1"
	r := o.Raw(sql)
	if start != nil && end != nil {
		sql = "SELECT %s%s%s FROM t_order_count_census WHERE TO_DAYS(create_date) >= TO_DAYS(?) AND TO_DAYS(create_date) <= TO_DAYS(?)"
		c1 := "SUM(order_count) order_count, SUM(order_cancel_count) order_cancel_count, SUM(appoint_count) appoint_count, SUM(appoint_cancel_count) appoint_cancel_count, SUM(one_by_one_count) one_by_one_count, SUM(one_by_one_cancel_count) one_by_one_cancel_count, "
		c2 := "MAX(order_total) order_total, MAX(order_cancel_total) order_cancel_total, MAX(appoint_total) appoint_total, MAX(appoint_cancel_total) appoint_cancel_total, MAX(one_by_one_total) one_by_one_total, MAX(one_by_one_cancel_total) one_by_one_cancel_total, "
		c3 := "MAX(create_date) create_date, MAX(create_time) create_time, MAX(update_time) update_time "
		sql = fmt.Sprintf(sql, c1, c2, c3)
		r = o.Raw(sql).SetArgs(start, end)
	}
	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "GetUserCensus failed, err: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *DynamicCensusService) UpsertDynamicCensusCount(ctx context.Context, data *models.OrderCountCensus) error {
	o := s.o
	var count int64
	sql := "SELECT count(1) FROM t_order_count_census WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	err := o.Raw(sql, data.CreateDate).QueryRow(&count)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "UpsertDynamicCensusCount.get failed, err: %v", err.Error())
		return err
	}
	sql = "UPDATE `t_order_count_census` SET `order_count` = ?, `order_cancel_count` = ?, `appoint_count` = ?, `appoint_cancel_count` = ?, `one_by_one_count` = ?, `one_by_one_cancel_count` = ?, `order_total` = ?, `order_cancel_total` = ?, `appoint_total` = ?, `appoint_cancel_total` = ?, `one_by_one_total` = ?, `one_by_one_cancel_total` = ?, `create_time` = ?, `update_time` = ? WHERE TO_DAYS(create_date) = TO_DAYS(?)"
	if count < 1 {
		sql = "INSERT INTO `t_order_count_census` (`order_count`, `order_cancel_count`, `appoint_count`, `appoint_cancel_count`, `one_by_one_count`, `one_by_one_cancel_count`, `order_total`, `order_cancel_total`, `appoint_total`, `appoint_cancel_total`, `one_by_one_total`, `one_by_one_cancel_total`, `create_time`, `update_time`, `create_date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	}
	result, err := o.Raw(sql,
		data.OrderCount, data.OrderCancelCount, data.AppointCount, data.AppointCancelCount, data.OneByOneCount, data.OneByOneCancelCount,
		data.OrderTotal, data.OrderCancelTotal, data.AppointTotal, data.AppointCancelTotal, data.OneByOneTotal, data.OneByOneCancelTotal,
		data.CreateTime, data.UpdateTime, data.CreateDate).Exec()
	if err != nil {
		logs.Error(ctx, "upsert UpsertDynamicCensusCount failed, err:", err)
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// GetMomentCount 根据时间获取非删除的动态
// status 0：正常显示  1:隐藏  -1：删除
func (s *DynamicCensusService) GetMomentCount(status, start, end int64) (*models.ContentView, error) {
	o := s.o
	var result models.ContentView
	s1 := "SELECT COUNT(`comment_id`) FROM t_comment_census_log WHERE `status` = ? and create_time >= ? AND create_time <= ?"
	s2 := "SELECT COUNT(*) FROM t_moment_census_log"
	sql := fmt.Sprintf("select (%v) moment_count, (%v) moment_total", s1, s2)
	err := o.Raw(sql, status, start, end).QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(context.TODO(), "GetVipCensusCount: %v", err.Error())
		return &result, err
	}
	return &result, nil
}

// GetCommentCount 根据时间获取评论的数量
// status 0：正常 -1：删除
func (s *DynamicCensusService) GetCommentCount(status, start, end int64) (*models.ContentView, error) {
	o := s.o
	var result models.ContentView
	s1 := "SELECT COUNT(`comment_id`) FROM t_comment_census_log WHERE `status` = ? and create_time >= ? AND create_time <= ?"
	s2 := "SELECT COUNT(*) FROM t_comment_census_log"
	sql := fmt.Sprintf("select (%v) comment_count, (%v) comment_total", s1, s2)
	err := o.Raw(sql, status, start, end).QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(context.TODO(), "GetVipCensusCount: %v", err.Error())
		return &result, err
	}
	return &result, nil
}

// GetOrderCount 根据时间获取订单的数量
// notStatus 1 待完成 4 已完成 5 已取消
func (s *DynamicCensusService) GetOrderCount(status, start, end int64) int64 {
	b := elastic.NewBoolQuery()
	if status != 0 {
		b.Must(elastic.NewTermQuery("status", status))
	}
	if start < 1 && end > 0 {
		b.Must(elastic.NewRangeQuery("create_time").Lte(end))
	} else if start > 0 && end > 0 {
		b.Must(elastic.NewRangeQuery("create_time").Gte(start).Lte(end))
	}

	aggregation := elastic.NewValueCountAggregation().Field("order_id") // 设置统计字段
	searchResult, err := elasticsearch.ES().Search().Index(elasticsearch.GetIndexName(models.IndexESDatingOrder)).
		Query(b).
		Aggregation("total", aggregation).
		Sort("create_time", false).
		Size(0).
		Do(context.TODO())
	if err != nil {
		logs.Errorf(context.TODO(), "GetOrderCensus failed, err: %v", err)
		return 0
	}
	agg, found := searchResult.Aggregations.ValueCount("total")
	if !found {
		return 0
	}
	return int64(*agg.Value)
}

// GetAppointCount 根据时间获取发单的数量
// notStatus 0 选人中  2 成单 3 用户取消 4超时取消 5后台取消
func (s *DynamicCensusService) GetAppointCount(status []int, start, end int64) int64 {
	b := elastic.NewBoolQuery()
	if len(status) > 0 {
		s := make([]interface{}, len(status))
		for index, value := range status {
			s[index] = value
		}
		b.Must(elastic.NewTermsQuery("status", s...))
	}
	if start < 1 {
		b.Must(elastic.NewRangeQuery("create_time").Lte(end))
	} else {
		b.Must(elastic.NewRangeQuery("create_time").Gte(start).Lte(end))
	}

	aggregation := elastic.NewValueCountAggregation().Field("dating_id") // 设置统计字段
	resp, err := elasticsearch.ES().Search().Index(elasticsearch.GetIndexName(models.IndexESDating)).
		Query(b).
		Aggregation("total", aggregation).
		Sort("create_time", false).
		Size(0).
		Do(context.TODO())
	if err != nil {
		logs.Errorf(context.TODO(), "GetAppointCount failed, err: %v", err)
		return 0
	}
	agg, found := resp.Aggregations.ValueCount("total")
	if !found {
		return 0
	}
	return int64(*agg.Value)
}

// GetOneByOneCount 根据时间获取一对一发单的数量
// notStatus 0 待接单 1 已接单 2拒绝接单 3超时未接单
func (s *DynamicCensusService) GetOneByOneCount(status []int, start, end int64) int64 {
	b := elastic.NewBoolQuery()
	if len(status) > 0 {
		s := make([]interface{}, len(status))
		for index, value := range status {
			s[index] = value
		}
		b.Must(elastic.NewTermsQuery("status", s...))
	}
	if start < 1 && end > 0 {
		b.Must(elastic.NewRangeQuery("create_time").Lte(end))
	} else if start > 0 && end > 0 {
		b.Must(elastic.NewRangeQuery("create_time").Gte(start).Lte(end))
	}

	resp, err := elasticsearch.ES().Search().Index(elasticsearch.GetIndexName(models.IndexESDatingAgain)).
		Query(b).
		Sort("create_time", false).
		Size(0).
		Do(context.TODO())
	if err != nil {
		logs.Errorf(context.TODO(), "GetOneByOneCount failed, err: %v", err)
		return 0
	}
	return resp.TotalHits()
}
