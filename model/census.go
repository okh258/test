package model

type VipCensusLog struct {
	OrderId    int64 `json:"order_id" orm:"pk"`
	Uid        int64 `json:"uid"`
	Duration   int64 `json:"duration"` // vip开通时长, 单位: 月
	CreateTime int64 `json:"create_time"`
}

func (m *VipCensusLog) TableName() string {
	return "t_vip_census_log"
}

// CapitalInflowCensus 资金流入统计
type CapitalInflowCensus struct {
	CreateTime           int64 `json:"create_time" orm:"pk"`    // 创建时间
	UpdateTime           int64 `json:"update_time"`             // 更新时间
	PayType              int64 `json:"pay_type"`                // 支付类型，1：支付宝，2：微信，3：苹果内购
	GoldAmount           int64 `json:"gold_amount"`             // 当日购买时币原始金额
	GoldAmountTotal      int64 `json:"gold_amount_total"`       // 购买时币原始总金额
	WalletAmount         int64 `json:"wallet_amount"`           // 当日购买原始钱包余额
	WalletAmountTotal    int64 `json:"wallet_amount_total"`     // 购买原始钱包总金额
	GoldPayAmount        int64 `json:"gold_pay_amount"`         // 当日购买时币支付金额
	GoldPayAmountTotal   int64 `json:"gold_pay_amount_total"`   // 购买时币支付总金额
	WalletPayAmount      int64 `json:"wallet_pay_amount"`       // 当日购买原始支付余额
	WalletPayAmountTotal int64 `json:"wallet_pay_amount_total"` // 购买支付钱包总金额
}

func (m *CapitalInflowCensus) TableName() string {
	return "t_capital_inflow_census"
}

// CapitalOutflowCensus 资金流出统计
type CapitalOutflowCensus struct {
	FromType    int64 `json:"from_type"`            // 提现来源，1：服务收益，2：邀请收益
	Amount      int64 `json:"amount"`               // 当日提现总金额
	AmountTotal int64 `json:"amount_total"`         // 提现总金额
	CreateTime  int64 `json:"create_time" orm:"pk"` // 创建时间
	UpdateTime  int64 `json:"update_time"`          // 更新时间
}

func (m *CapitalOutflowCensus) TableName() string {
	return "t_capital_outflow_census"
}
