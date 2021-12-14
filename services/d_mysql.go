package services

import (
	"context"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	"test/db"
	"test/model"
	"test/models"
)

var _dService *DService

func NewDService() *DService {
	if _dService == nil {
		_dService = &DService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameDating),
		}
		orm.RegisterModel(&model.CapitalInflowCensus{}, &model.CapitalOutflowCensus{}, &models.BusinessOrderCensus{})
	}
	return _dService
}

type DService struct {
	o orm.Ormer
}

// GetInviteeUserCount 获取有邀请码的用户数量
func (s *DService) GetInviteeUserCount(start, end int64) int64 {
	var resultCount int64
	o := s.o
	sql := "SELECT COUNT(DISTINCT invitee_uid) FROM t_invitee_info WHERE create_time >= ? AND create_time <= ?"
	err := o.Raw(sql, start, end).QueryRow(&resultCount)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(context.TODO(), "GetInviteeUserCount: %v", err.Error())
		return 0
	}
	return resultCount
}

// GetOnlineOrderAmount 获取线上订单总金额
// 平台历史线上订单总价 = 线上品类+视频通话时送礼物总花费的时币
// 用户线上单总支付（时币）= 线上品类+视频通话时送礼物总花费的时币-优惠券减免金额，单位为时币
// trade_type 交易类型 1：收入 2：支出
// is_online 是否线上类型 1线上 0线下
// category_id 分类编号 8881：视频礼物
func (s *DService) GetOnlineOrderAmount(ctx context.Context, tradeType, start, end int64) int64 {
	var total int64
	sql := "SELECT SUM(h.amount) amount FROM t_inviter_trade_history h " +
		"LEFT JOIN t_category c ON c.category_id = h.category_id " +
		"WHERE trade_type = 1 AND (c.is_online = 1 OR h.category_id = 8881)"
	r := s.o.Raw(sql, tradeType)
	if start > 0 && end > 0 {
		sql += "AND create_time >= ? AND create_time <= ?"
		r = s.o.Raw(sql, tradeType, start, end)
	}
	err := r.QueryRow(&total)
	if err != nil {
		logs.Errorf(ctx, "GetOnlineOrderInAmount: %v", err.Error())
		return 0
	}
	return total
}

// GetOfflineOrderAmount 获取线下订单总金额
// 平台历史线下订单格总价=用户线下单订单金额
// 用户线下单总支付（元）= 用户线下单订单金额-优惠券减免金额
// trade_type 交易类型 1：收入 2：支出
// is_online 是否线上类型 1线上 0线下
func (s *DService) GetOfflineOrderAmount(ctx context.Context, tradeType, start, end int64) int64 {
	var total int64
	sql := "SELECT SUM(h.amount) amount FROM t_inviter_trade_history h " +
		"LEFT JOIN t_category c ON c.category_id = h.category_id  " +
		"WHERE h.trade_type = ? AND c.is_online = 0"
	r := s.o.Raw(sql, tradeType)
	if start > 0 && end > 0 {
		sql += "AND create_time >= ? AND create_time <= ?"
		r = s.o.Raw(sql, tradeType, start, end)
	}
	err := r.QueryRow(&total)
	if err != nil {
		logs.Errorf(ctx, "GetOfflineOrderInAmount: %v", err.Error())
		return 0
	}
	return total
}

// GetAmountByType 根据类型获取总金额
// trade_type 交易类型 1：收入 2：支出 TODO 该字段建议加索引
// profit_type 收益来源类型 1：订单 2：礼物
func (s *DService) GetAmountByType(ctx context.Context, tradeType, profitType, start, end int64) int64 {
	var total int64
	sql := "SELECT SUM(amount) amount FROM t_inviter_trade_history WHERE trade_type = ? AND profit_type = ?"
	r := s.o.Raw(sql, tradeType, profitType)
	if start > 0 && end > 0 {
		sql += "AND create_time >= ? AND create_time <= ?"
		r = s.o.Raw(sql, tradeType, profitType, start, end)
	}
	err := r.QueryRow(&total)
	if err != nil {
		logs.Errorf(ctx, "GetAmountByType: %v", err.Error())
		return 0
	}
	return total
}
