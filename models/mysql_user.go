package models

import (
	"time"
)

type UserCountCensus struct {
	CreateDate          time.Time `json:"create_date" orm:"pk"`
	VerifyCount         int64     `json:"verify_count"`   // 本日输入邀请码的注册人数
	RegisterCount       int64     `json:"register_count"` // 注册量
	RegisterMaleCount   int64     `json:"register_male_count"`
	RegisterFemaleCount int64     `json:"register_female_count"`
	CoachCount          int64     `json:"coach_count"`
	CoachMaleCount      int64     `json:"coach_male_count"`
	CoachFemaleCount    int64     `json:"coach_female_count"`
	ConsumerCount       int64     `json:"consumer_count"`
	ConsumerMaleCount   int64     `json:"consumer_male_count"`
	ConsumerFemaleCount int64     `json:"consumer_female_count"`
	IosCount            int64     `json:"ios_count"`
	AndroidCount        int64     `json:"android_count"`
	VerifyTotal         int64     `json:"verify_total"` // 输入邀请码的注册总人数
	RegisterTotal       int64     `json:"register_total"`
	CoachTotal          int64     `json:"coach_total"`
	ConsumerTotal       int64     `json:"consumer_total"`
	IosTotal            int64     `json:"ios_total"`
	AndroidTotal        int64     `json:"android_total"`
	CreateTime          int64     `json:"create_time"`
	UpdateTime          int64     `json:"update_time"`
}

func (m *UserCountCensus) TableName() string {
	return "t_user_count_census"
}
