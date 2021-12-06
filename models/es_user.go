package models

type SearchUserRequestParams struct {
	// 用户uid
	Uid int64 `json:"uid"`
	// userNumber
	UserNumber int64 `json:"usernumber"`
	// 性别 0全部 1男 2女
	Gender int64 `json:"gender"`
	// 手机号
	Mobile int64 `json:"mobile"`
	// 昵称
	Nickname string `json:"nickname"`
	// 是否vip -1否 0全部 1是
	IsVip int32 `json:"is_vip"`
	// 是否服务者 -1否 0全部 1是
	IsCoach int64 `json:"is_coach"`
	// 用户状态 -1封禁 1正常
	UserStatus int64 `json:"user_status"`
	// 条数
	Num int64 `json:"num"`
	// 分页page
	Page int64 `json:"page"`
	// 开始时间
	StartTime int64 `json:"startTime"`
	// 结束时间
	EndTime int64 `json:"endTime"`
}
