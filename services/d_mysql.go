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
		orm.RegisterModel(&models.CallVideoOrder{}, &model.CapitalInflowCensus{}, &model.CapitalOutflowCensus{},
			&models.BusinessOrderCensus{}, &models.CategoryCoachSignature{}, &models.SkillCategory{},
			&models.RechargeGoldLog{})
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

func (s *DService) GetAllCalls(ctx context.Context, params models.CallVideoOrderReq) ([]*models.CallVideoOrder, error) {
	var resultList []*models.CallVideoOrder
	query := s.o.QueryTable("t_call_video_order")
	if len(params.UserIdList) > 0 {
		query = query.Filter("uid__in", params.UserIdList)
	}
	if len(params.CoachUidList) > 0 {
		query = query.Filter("coach_uid__in", params.CoachUidList)
	}
	if params.Status == 5 {
		query = query.Filter("status", params.Status)
	} else if params.Status == -1 {
		query = query.Filter("status__gte", 5)
	}
	if params.OrderId > 0 {
		query = query.Filter("order_id", params.OrderId)
	}
	if params.StartTime > 0 {
		query = query.Filter("create_time__gte", params.StartTime)
	}
	if params.EndTime > 0 {
		query = query.Filter("create_time__lte", params.EndTime)
	}
	query = query.OrderBy("-create_time")

	query = query.Offset((params.PageNum - 1) * params.PageSize).Limit(params.PageSize)

	_, err := query.All(&resultList)
	if err != nil {
		logs.Error(ctx, "CallVideoOrderService:GetAllCalls err:", err)
		return nil, err
	}
	return resultList, nil
}

// GetSignatureCount 获取服务者设置的标签数量
func (s *DService) GetSignatureCount(ctx context.Context, categoryID, isSelected, coachUid int64) (int64, error) {
	q := orm.NewOrm().QueryTable("t_category_coach_signature").Filter("coach_uid", coachUid)
	if isSelected != -1 {
		// 是否选中: 1：选中 0：未选中
		q = q.Filter("is_selected", isSelected)
	}
	if categoryID != 0 {
		q = q.Filter("category_id", categoryID)
	}

	count, err := q.Count()
	if err != nil {
		logs.Error(ctx, "CategoryCoachSignatureService：GetSignatureCount err:", err)
		return count, err
	}
	return count, nil
}

func (s *DService) GetCategoryListByIds(ctx context.Context, categoryIds []int64) (map[int64]*models.SkillCategory, error) {
	if len(categoryIds) == 0 {
		return nil, nil
	}

	var result []*models.SkillCategory
	query := s.o.QueryTable("t_category")
	if len(categoryIds) > 0 {
		query = query.Filter("category_id__in", categoryIds)
	}
	_, err := query.OrderBy("-sort_num").All(&result)
	if err != nil {
		logs.Error(ctx, "GetCategoryListByIds failed, err:", err)
		return nil, err
	}

	var data = make(map[int64]*models.SkillCategory)
	for _, val := range result {
		data[val.CategoryId] = val
	}
	return data, nil
}

func (s *DService) PageRechargeGoldLog(ctx context.Context, params *models.RechargeInviteReq) ([]*models.RechargeGoldLog, int64, error) {
	// 初始化页码, 页大小
	if params.PageNum == 0 {
		params.PageNum = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 1
	}
	q := orm.NewOrm().QueryTable("t_recharge_gold_log")
	if params.Status > 0 {
		q = q.Filter("status", params.Status)
	}
	if params.InviteeUid > 0 {
		q = q.Filter("uid", params.InviteeUid)
	}
	if params.InviterUid > 0 {
		q = q.Filter("inviter_uid", params.InviterUid)
	}

	var result []*models.RechargeGoldLog
	_, err := q.Limit(params.PageSize, (params.PageNum-1)*params.PageSize).
		OrderBy("-create_time").
		All(&result)
	if err != nil {
		logs.Errorf(ctx, "PageRechargeGoldLog failed, err: %v", err)
		return nil, 0, err
	}

	totalNum, err := q.Count()
	if err != nil {
		logs.Errorf(ctx, "PageRechargeGoldLog total failed, err: %v", err)
		return nil, 0, err
	}

	return result, totalNum, nil
}
