package models

import "time"

const (
	VipTypeOneMonth   int64 = 1  // 一个月
	VipTypeThreeMonth int64 = 3  // 三个月
	VipTypeOneYear    int64 = 12 // 一年
	VipTypeTwoYear    int64 = 24 // 两年
)

type VipCensusLog struct {
	OrderId    int64 `json:"order_id" orm:"pk"`
	Uid        int64 `json:"uid"`
	Duration   int64 `json:"duration"` // vip开通时长, 单位: 月
	CreateTime int64 `json:"create_time"`
}

func (m *VipCensusLog) TableName() string {
	return "t_vip_census_log"
}

type VipCensusView struct {
	VipTotal        int64   `json:"vip_total"`         // 会员总量
	VipAvg          float64 `json:"vip_avg"`           // 会员开通平均时长
	VipCount        int64   `json:"vip_count"`         // 会员开通次数
	OneMonthCount   int64   `json:"one_month_count"`   // 1月会员开通次数
	ThreeMonthCount int64   `json:"three_month_count"` // 3月会员开通次数
	OneYearCount    int64   `json:"one_year_count"`    // 1年会员开通次数
	TwoYearCount    int64   `json:"two_year_count"`    // 2年会员开通次数
}

type DynamicCountCensus struct {
	CreateDate          *time.Time `json:"create_date" orm:"pk"`    // 统计记录时间, 保留到日
	MomentCount         int64      `json:"moment_count"`            // 本日动态数量
	CommentCount        int64      `json:"comment_count"`           // 本日评论数量
	OrderCount          int64      `json:"order_count"`             // 本日订单数量
	OrderCancelCount    int64      `json:"order_cancel_count"`      // 本日已取消订单数量
	AppointCount        int64      `json:"appoint_count"`           // 本日发单数量
	AppointCancelCount  int64      `json:"appoint_cancel_count"`    // 本日已取消发单数量
	OneByOneCount       int64      `json:"one_by_one_count"`        // 本日一对一下单数量
	OneByOneCancelCount int64      `json:"one_by_one_cancel_count"` // 本日取消一对一下单数量
	MomentTotal         int64      `json:"moment_total"`            // 动态总数量
	CommentTotal        int64      `json:"comment_total"`           // 评论总数量
	OrderTotal          int64      `json:"order_total"`             // 订单总数量
	OrderCancelTotal    int64      `json:"order_cancel_total"`      // 已取消订单总数量
	AppointTotal        int64      `json:"appoint_total"`           // 发单总数量
	AppointCancelTotal  int64      `json:"appoint_cancel_total"`    // 已取消发单总数量
	OneByOneTotal       int64      `json:"one_by_one_total"`        // 一对一下单总数量
	OneByOneCancelTotal int64      `json:"one_by_one_cancel_total"` // 已取消一对一下单总数量
	CreateTime          int64      `json:"create_time"`
	UpdateTime          int64      `json:"update_time"`
}
