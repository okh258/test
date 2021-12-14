package model

//order表结构
type Order struct {
	//用户ID
	Uid int64 `odm:"PK" json:"uid"`
	//订单ID
	OrderId int64 `odm:"SK" json:"order_id"`
	//商品ID
	GoodsId int64 `json:"goods_id"`
	//购买数量
	GoodsCount int64 `json:"goods_count"`
	//支付类型，1：支付宝，2：微信，3：内购
	PayType int64 `json:"pay_type"`
	//真实支付金额
	PayAmount int64 `json:"pay_amount"`
	//原始金额
	RealAmount int64 `json:"real_amount"`
	//货币代码
	Currency string `json:"currency"`
	//来源渠道
	Channel string `json:"channel"`
	//为了什么付费
	BuyFor string `json:"buy_for"`
	//状态，0：待处理，1：成功
	Status int64 `odm:"PK|status_index|All" json:"status"`
	//创建时间
	CreateTime int64 `odm:"SK|status_index|All" json:"create_time"`
	//支付时间
	PayTime int64 `json:"pay_time"`
	//客户端类型
	OsType string `json:"os_type"`
	//苹果内购产品id
	ProductId string `json:"product_id,omitempty"`
}
