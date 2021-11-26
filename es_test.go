package test

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
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

func TestGetMomentCount(t *testing.T) {
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
