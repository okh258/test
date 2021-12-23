package models

// CallVideoOrderReq 视频通话查询req
type CallVideoOrderReq struct {
	//服务者uid
	CoachUid int64 `json:"coach_uid"`
	//服务者昵称
	CoachNickName string `json:"coach_nick_name"`
	//用户编号
	UserId int64 `json:"uid"`
	//用户昵称
	UserNickName string `json:"user_nick_name"`
	//服务者uid 集合 -- 昵称可能查出多个用户
	CoachUidList []int64 `json:"coach_uid_list"`
	//用户uid集合 -- 昵称可能查出多个用户
	UserIdList []int64 `json:"user_id_list"`
	//状态
	Status int64 `json:"status"`
	//订单编号
	OrderId int64 `json:"order_id"`
	//开始时间
	StartTime int64 `json:"start_time"`
	//结束时间
	EndTime int64 `json:"end_time"`
	//页码
	PageNum int64 `json:"page_num"`
	//每页行数
	PageSize int64 `json:"page_size"`
}

type CallVideoOrder struct {
	// 主键自增
	CallVideoOrderId int64 `orm:"auto" json:"call_video_order_id"`
	// 关联视频订单id
	OrderId int64 `json:"order_id"`
	// 用户uid
	Uid int64 `json:"uid"`
	// 服务者uid
	CoachUid int64 `json:"coach_uid"`
	// 价格
	Price int64 `json:"price"`
	// 状态 1 正在发起视频 2用户主动取消 3服务者拒接 4超时没响应取消 5 成功接听 6视频无心跳
	Status int64 `json:"status"`
	// 取消/拒接方 ,1是用户， 2是服务者
	RefuseFrom int64 `json:"refuse_from"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
	// channelId
	ChannelId string `json:"channel_id"`
	// 是否vip 0否1是
	IsVip int64 `json:"is_vip"`
	// 取消时间
	CancelTime int64 `json:"cancel_time"`
	// 相关信息记录
	RelateInfo string `json:"relate_info"`
	// 视频开始时间
	StartTime int64 `json:"start_time"`
	// 视频结束时间
	EndTime int64 `json:"end_time"`
	// 实时消费金额
	PromPrice int64 `json:"prom_price"`
	// 总消费金额
	TotalAmount int64 `json:"total_amount"`
	// 通话时长
	TotalCallTime int64 `json:"total_call_time"`
	// 时长不足提示充值时间
	TimerNotify int64 `json:"timer_notify"`
	// 后台审核状态 0 待审核 1 审核通过 2 拒绝
	CheckStatus int64 `json:"check_status"`
	// 审核时间
	CheckTime int64 `json:"check_time"`
	// 视频用户端截图
	UserScreenShot string `json:"user_screen_shot"`
	// 视频服务者端截图
	CoachScreenShot string `json:"coach_screen_shot"`
	//用户是否开启模糊
	IsVague        int64 `json:"is_vague"`
	IsIllegal      int   `json:"is_illegal"`       //用户是否违规
	CoachIsIllegal int   `json:"coach_is_illegal"` //服务者是否违规

	// 仅页面展示字段
	Nickname      string `json:"nickname" orm:"-"`       // 用户昵称
	CoachNickname string `json:"coach_nickname" orm:"-"` // 服务者昵称
}

func (m *CallVideoOrder) TableName() string {
	return "t_call_video_order"
}
