package test

import (
	"context"
	"encoding/json"
	"fmt"
	"git.devops.com/wsim/hflib/logs"
	"github.com/olivere/elastic/v7"
	"test/elasticsearch"
	"test/models"
	"test/util"
	"testing"
)

func TestGetLockOrderNum(t *testing.T) {
	ctx := context.Background()
	num := GetLockOrderNum(t, ctx, 0, util.MicroTime())
	t.Logf("num: %v", num)
}

// GetLockOrderNum 获取已锁订单数量
func GetLockOrderNum(t *testing.T, ctx context.Context, uid int64, now int64) int {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	b := elastic.NewBoolQuery()
	b.Must(elastic.NewTermQuery("is_locked", "1"))
	b.Must(elastic.NewRangeQuery("create_time").Gte(1636788149934318).Lte(now))
	src, err := b.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}

	t.Logf("%s", data)

	aggregation := elastic.NewValueCountAggregation().Field("order_id") // 设置统计字段
	searchResult, err := elasticClient.Search().Index("uat_tm_dating_order").
		Query(b).
		Aggregation("total", aggregation).
		Size(0).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Errorf("GetLockOrderNum failed, err: %v", err)
		return 0
	}
	// 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.ValueCount("total")
	if found {
		// 打印结果，注意：这里使用的是取值运算符
		t.Logf("lock order num: %v", *agg.Value)
	}
	return int(*agg.Value)
}

func TestGetCensorMsgCount(t *testing.T) {
	ctx := context.Background()
	//num1 := GetCensorMomentCount(t, ctx, 0, util.MicroTime())
	//num2 := GetCensorCommentCount(t, ctx, 0, util.MicroTime())
	num3 := GetCensorMsgCount(t, ctx, 0, util.MicroTime())
	//t.Logf("num1: %v", num1)
	//t.Logf("num2: %v", num2)
	t.Logf("num3: %v", num3)
}

const (
	CensorResultPass   = "pass"   // 通过
	CensorResultReview = "review" // 人工审核
	CensorResultBlock  = "block"  // 违规
)

// GetCensorMomentCount 获取需要审核的动态
func GetCensorMomentCount(t *testing.T, ctx context.Context, start, end int64) int {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	b := elastic.NewBoolQuery()
	b.Must(elastic.NewTermQuery("source", "censor"))
	b.MustNot(elastic.NewTermQuery("status", -1))
	b.Must(elastic.NewTermsQuery("censor_result", CensorResultReview, CensorResultBlock))
	b.Must(elastic.NewRangeQuery("update").Gte(start).Lte(end))

	aggregation := elastic.NewValueCountAggregation().Field("moment_id") // 设置统计字段
	searchResult, err := elasticClient.Search().Index("tm_uat_moment").
		Query(b).
		Aggregation("total", aggregation).
		Size(0).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Errorf("GetMomentCount failed, err: %v", err)
		return 0
	}
	// 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.ValueCount("total")
	if !found {
		return 0
	}
	return int(*agg.Value)
}

// GetCensorCommentCount 获取需要审核的动态
func GetCensorCommentCount(t *testing.T, ctx context.Context, start, end int64) int {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	b := elastic.NewBoolQuery()
	b.Must(elastic.NewTermQuery("source", "censor"))
	b.Must(elastic.NewTermsQuery("censor_result", CensorResultReview, CensorResultBlock))
	b.Must(elastic.NewRangeQuery("update").Gte(start).Lte(end))

	aggregation := elastic.NewValueCountAggregation().Field("comment_id") // 设置统计字段
	searchResult, err := elasticClient.Search().Index("tm_uat_comment").
		Query(b).
		Aggregation("total", aggregation).
		Size(0).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Errorf("GetMomentCount failed, err: %v", err)
		return 0
	}
	// 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.ValueCount("total")
	if !found {
		return 0
	}
	return int(*agg.Value)
}

// GetCensorMsgCount 获取需要审核的消息
func GetCensorMsgCount(t *testing.T, ctx context.Context, start, end int64) int {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	b := elastic.NewBoolQuery()
	b.Must(elastic.NewTermQuery("source", "censor"))
	b.MustNot(elastic.NewTermQuery("status", -1))
	b.Must(elastic.NewTermsQuery("censor_result", CensorResultReview, CensorResultBlock))
	b.Must(elastic.NewRangeQuery("update").Gte(start).Lte(end))

	aggregation := elastic.NewValueCountAggregation().Field("msg_id.keyword") // 设置统计字段
	searchResult, err := elasticClient.Search().Index("tm_uat_message").
		Query(b).
		Aggregation("total", aggregation).
		Size(0).
		Pretty(true).
		Do(ctx)
	if err != nil {
		t.Errorf("GetCensorMsgCount failed, err: %v", err)
		return 0
	}
	// 使用ValueCount函数和前面定义的聚合条件名称，查询结果
	agg, found := searchResult.Aggregations.ValueCount("total")
	if !found {
		return 0
	}
	return int(*agg.Value)
}

func TestGetUserCount(t *testing.T) {
	//d, _ := time.Parse("2006-01-02", "2021-10-23")
	//d, _ := time.Parse("2006-01-02", "2021-12-06")
	total := GetUserCount(models.SearchUserRequestParams{
		Gender:  0,
		IsCoach: 0,
		//StartTime: d.UnixNano() / 1000000,
		//StartTime: util.GetZero(time.Now()).UnixNano() / 1000000,
		StartTime: 0,
		EndTime:   util.Timestamp(),
	})
	t.Logf("total: %v", total)
}

func GetUserCount(params models.SearchUserRequestParams) int64 {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://172.31.15.17:9200/"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	indexName := "uat_tm_user"
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
	if params.EndTime != 0 && params.StartTime == 0 {
		query.Must(elastic.NewRangeQuery("register_time").Lte(params.EndTime))
	} else if params.EndTime != 0 && params.StartTime != 0 {
		query.Must(elastic.NewRangeQuery("register_time").Gte(params.StartTime).Lte(params.EndTime))
	}

	resp, err := elasticClient.Search().Index(indexName).
		Query(query).
		Sort("register_time", false).
		Size(0).
		Do(context.TODO())
	if err != nil {
		return 0
	}

	return resp.TotalHits()
}

func TestSetSkillAppearanceLevel(t *testing.T) {
	SetSkillAppearanceLevel(context.TODO(), 528041712999424, 2)
}

// SetSkillAppearanceLevel 设置服务者列表颜值信息
func SetSkillAppearanceLevel(ctx context.Context, uid, appearanceLevel int64) error {
	indexName := elasticsearch.GetIndexName("user_skill")
	_, err := elasticsearch.ES().UpdateByQuery().Query(elastic.NewTermQuery("uid", uid)).
		Index(indexName).
		Script(elastic.NewScript(fmt.Sprintf("ctx._source['appearance_level']=%v", appearanceLevel))).
		Do(context.TODO())
	if err != nil && !elastic.IsNotFound(err) {
		logs.Errorf(ctx, "update es coach info failed, err: %v", err)
		return err
	}
	return nil
}
