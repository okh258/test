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

func (m *CategoryCoachSignature) TableName() string {
	return "t_category_coach_signature"
}

//用户订单个性标签
type CategoryCoachSignature struct {
	//编号
	Id int64 `orm:"auto" json:"id"`
	//个性标签名称
	SignatureName string `json:"signature_name"`
	//分类编号
	CategoryId int64 `json:"category_id"`
	//用户uid
	CoachUid int64 `json:"coach_uid"`
	//创建时间
	CreateTime int64 `json:"create_time"`
	//状态 1：审核中 2：通过 3：未通过
	Status int64 `json:"status"`
	//是否被选中 1：选中 0：未选中
	IsSelected int64 `json:"is_selected"`
	//来源编号 1：系统推荐标签 2：自选
	SourceId int64 `json:"source_id"`
	//是否显示 0：不显示 1：显示
	IsShow int64 `json:"is_show"`
}

// 发约品类
type SkillCategory struct {
	// 品类id
	CategoryId int64 `orm:"auto" json:"category_id"`
	// 品类名
	CategoryName string `json:"category_name"`
	// 品类描述
	CategoryDesc string `json:"category_desc"`
	// 品类类型 1线上服务 2休闲娱乐 3运动陪练 4教育培训 5时尚达人 6独家专属
	CategoryType int64 `json:"category_type"`
	// 搜索关键词
	SearchKeyword string `json:"search_keyword"`
	// 是否是线上 1 线上 0线下
	IsOnline int64 `json:"is_online"`
	// 是否饮酒
	IsDrink int64 `json:"is_drink"`
	// 最小时长(秒)
	MinDuration int64 `json:"min_duration"`
	// 最大时长
	MaxDuration int64 `json:"max_duration"`
	// 价格范围
	PriceRange string `json:"price_range"`
	// 排序
	SortNum int64 `json:"sort_num"`
	// 背景颜色
	BackgroundColor string `json:"background_color"`
	// 状态 0正常 -1禁用
	Status int64 `json:"status"`
}

func (m *SkillCategory) TableName() string {
	return "t_category"
}

// RechargeInviteReq 邀请收益
type RechargeInviteReq struct {
	Uid        int64 `json:"uid"`         // 受邀人uid
	InviterUid int64 `json:"inviter_uid"` // 邀请人uid
	Status     int64 `json:"status"`      // 状态: 0未打款 1已打款
	PageNum    int64 `json:"page_num"`    // 页码
	PageSize   int64 `json:"page_size"`   // 页大小
}

type RechargeGoldLog struct { // 充值时币记录表
	Id         int64 `json:"id" orm:"auto"`
	OrderId    int64 `json:"order_id"`
	Status     int64 `json:"status"` // 状态: 0未打款  1已打款
	Amount     int64 `json:"amount"`
	Uid        int64 `json:"uid"`
	InviterUid int64 `json:"inviter_uid"` // 邀请人uid
	CreateTime int64 `json:"create_time"`
}

func (m *RechargeGoldLog) TableName() string {
	return "t_recharge_gold_log"
}
