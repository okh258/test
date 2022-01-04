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
	OsType     int64 `json:"os_type"`
	// 条数
	Num int64 `json:"num"`
	// 分页page
	Page int64 `json:"page"`
	// 开始时间
	StartTime int64 `json:"startTime"`
	// 结束时间
	EndTime int64 `json:"endTime"`
}

type UserBaseInfo struct {
	Avatar              string  `json:"avatar"`
	UserId              int64   `json:"user_id" odm:"PK"`
	Nickname            string  `json:"nickname"`
	Gender              int32   `json:"gender"`
	Mobile              string  `json:"mobile"`
	Usernumber          int64   `json:"usernumber"`
	Username            string  `json:"username"`
	Age                 int32   `json:"age"`
	Place               string  `json:"place"`
	Birthday            int64   `json:"birthday"`
	VipLevel            int32   `json:"vip_level"`
	VipExpireTime       int64   `json:"vip_expire_time"`
	Vest                []int64 `json:"vest"` // 用户马甲: 1测试号
	UserStatus          int32   `json:"user_status"`
	TextSignature       string  `json:"text_signature"`
	CertifyStatus       int32   `json:"certify_status"`
	UserIdentity        int32   `json:"user_identity"`
	RegisterTime        int64   `json:"register_time"`
	LastVipBuyTime      int64   `json:"last_vip_buy_time"`
	LatLng              string  `json:"lat_lng"`
	LoginIp             string  `json:"login_ip"`
	LoginTime           int64   `json:"login_time"`
	LoginOS             string  `json:"login_os"`
	BelongAgentUid      int64   `json:"belong_agent_uid"`
	BelongAgentNickname string  `json:"belong_agent_nickname"`
	InviteCount         int64   `json:"invite_count"`
	RegisterIp          string  `json:"register_ip"`
	RegisterOsType      int32   `json:"register_os_type"`
	RegisterDeviceToken string  `json:"register_device_token"`
	RegisterDeviceName  string  `json:"register_device_name"`
	RegisterDeviceType  string  `json:"register_device_type"`
}

func (db *UserBaseInfo) GetName() string {
	return "user_user"
}
