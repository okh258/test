package timer

import (
	"context"
	"test/models"
	"test/services"
	"test/util"
	"time"

	"git.devops.com/wsim/hflib/logs"
)

type CensusService struct {
}

var _censusService = &CensusService{}

func NewCensusService() *CensusService {
	return _censusService
}

func getCensus(ctx context.Context, now time.Time) func() error {
	logs.Infof(ctx, "start census task, now: %v", now.Format("2006-04-02 15:04:05"))
	//// 获取用户统计信息, 注册数量, 服务者数量, 非服务者数量
	NewCensusService().GetUserCensus(ctx, now)
	// 订单数量, 发单数量, 一对一下单数量
	//NewCensusService().GetOrderCensus(ctx, now)
	// 业务数据统计
	//NewCensusService().GetBusinessCensus(ctx, now)
	return nil
}

// GetUserCensus 获取用户统计信息
func (s *CensusService) GetUserCensus(ctx context.Context, now time.Time) {
	// 获取当天的用户统计信息
	oldUserCensus, err := services.NewUserCensusService().GetUserCensusCount(ctx, &now, nil)
	if err != nil {
		logs.Errorf(ctx, "GetUserCensus failed, date: %v, err: %v", now.Format(time.RFC3339), err)
		return
	}

	// 无存量数据处理
	if oldUserCensus == nil || oldUserCensus.CreateDate.IsZero() {
		userCensus := &models.UserCountCensus{
			CreateDate: now,
			CreateTime: now.UnixNano() / 1000000,
			UpdateTime: now.UnixNano() / 1000000,
		}
		p := models.SearchUserRequestParams{
			StartTime: GetZero(now).UnixNano() / 1000000,
			EndTime:   userCensus.UpdateTime,
		}
		// 本日输入邀请码的注册人数
		userCensus.VerifyCount = services.NewDService().GetInviteeUserCount(p.StartTime*1000, p.EndTime*1000)
		// 注册数量
		userCensus.RegisterCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 1
		userCensus.RegisterMaleCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 2
		userCensus.RegisterFemaleCount = services.NewUserCensusService().GetUserCountBy(p)
		// 服务者数量
		p.IsCoach = 1
		userCensus.CoachCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 1
		userCensus.RegisterMaleCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 2
		userCensus.RegisterFemaleCount = services.NewUserCensusService().GetUserCountBy(p)
		// 非服务者数量
		p.IsCoach = -1
		userCensus.ConsumerCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 1
		userCensus.ConsumerMaleCount = services.NewUserCensusService().GetUserCountBy(p)
		p.Gender = 2
		userCensus.ConsumerFemaleCount = services.NewUserCensusService().GetUserCountBy(p)
		// 注册总数
		pTotal := models.SearchUserRequestParams{
			StartTime: 0,
			EndTime:   now.UnixNano() / 1000000,
		}
		// 总输入邀请码的注册人数
		userCensus.VerifyTotal = services.NewDService().GetInviteeUserCount(0, p.EndTime*1000)
		// 注册总数量
		userCensus.RegisterTotal = services.NewUserCensusService().GetUserCountBy(pTotal)
		// 服务者总数量
		pTotal.IsCoach = 1
		userCensus.CoachTotal = services.NewUserCensusService().GetUserCountBy(pTotal)
		// 非服务者总数量
		pTotal.IsCoach = -1
		userCensus.ConsumerTotal = services.NewUserCensusService().GetUserCountBy(pTotal)
		logs.Infof(ctx, "UpsertUserCensus: %+v", userCensus)
		err = services.NewUserCensusService().UpsertUserCensus(ctx, userCensus)
		if err != nil {
			logs.Errorf(ctx, "upsertUserCountCensus failed, date: %+v, err: %v", userCensus, err)
			return
		}
		return
	}
	// 有存量数据处理
	userCensus := *oldUserCensus
	p := models.SearchUserRequestParams{
		StartTime: userCensus.UpdateTime,
		EndTime:   now.UnixNano() / 1000000,
	}
	// 本日输入邀请码的注册人数
	userCensus.VerifyCount += services.NewDService().GetInviteeUserCount(p.StartTime*1000, p.EndTime*1000)
	// 注册数量
	userCensus.RegisterCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 1
	userCensus.RegisterMaleCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 2
	userCensus.RegisterFemaleCount += services.NewUserCensusService().GetUserCountBy(p)
	// 服务者数量
	p.IsCoach = 1
	userCensus.CoachCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 1
	userCensus.RegisterMaleCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 2
	userCensus.RegisterFemaleCount += services.NewUserCensusService().GetUserCountBy(p)
	// 非服务者数量
	p.IsCoach = -1
	userCensus.ConsumerCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 1
	userCensus.ConsumerMaleCount += services.NewUserCensusService().GetUserCountBy(p)
	p.Gender = 2
	userCensus.ConsumerFemaleCount += services.NewUserCensusService().GetUserCountBy(p)
	// 总数量
	pTotal := models.SearchUserRequestParams{
		StartTime: 0,
		EndTime:   now.UnixNano() / 1000000,
	}
	// 总输入邀请码的注册人数
	userCensus.VerifyTotal = services.NewDService().GetInviteeUserCount(0, pTotal.EndTime*1000)
	// 注册总数量
	userCensus.RegisterTotal = services.NewUserCensusService().GetUserCountBy(pTotal)
	// 服务者总数量
	pTotal.IsCoach = 1
	userCensus.CoachTotal = services.NewUserCensusService().GetUserCountBy(pTotal)
	// 非服务者总数量
	pTotal.IsCoach = -1
	userCensus.ConsumerTotal = services.NewUserCensusService().GetUserCountBy(pTotal)

	userCensus.UpdateTime = now.UnixNano() / 1000000
	logs.Infof(ctx, "UpsertUserCensus: %+v", userCensus)
	err = services.NewUserCensusService().UpsertUserCensus(ctx, &userCensus)
	if err != nil {
		logs.Errorf(ctx, "upsertUserCountCensus failed, date: %+v, err: %v", userCensus, err)
		return
	}
}

// GetOrderCensus 获取订单统计数
func (s *CensusService) GetOrderCensus(ctx context.Context, now time.Time) {
	// 获取当天的统计信息
	oldCensus, err := services.NewDynamicCensusService().GetDynamicCensusCount(ctx, &now, nil)
	if err != nil {
		logs.Errorf(ctx, "GetOrderCensus failed, date: %v, err: %v", now.Format(time.RFC3339), err)
		return
	}

	// 无存量数据处理
	if oldCensus == nil || oldCensus.CreateDate == nil {
		census := &models.OrderCountCensus{
			CreateDate: &now,
			CreateTime: now.UnixNano() / 1000,
			UpdateTime: now.UnixNano() / 1000,
		}
		logs.Infof(ctx, "dynamicCensus: %+v", census)
		// 本日数量
		census.OrderCount = services.NewDynamicCensusService().GetOrderCount(0, GetZero(now).UnixNano()/1000, census.UpdateTime)
		census.OrderCancelCount = services.NewDynamicCensusService().GetOrderCount(5, GetZero(now).UnixNano()/1000, census.UpdateTime)
		census.AppointCount = services.NewDynamicCensusService().GetAppointCount(nil, GetZero(now).UnixNano()/1000, census.UpdateTime)
		census.AppointCancelCount = services.NewDynamicCensusService().GetAppointCount([]int{3, 4, 5}, GetZero(now).UnixNano()/1000, census.UpdateTime)
		census.OneByOneCount = services.NewDynamicCensusService().GetOneByOneCount(nil, GetZero(now).UnixNano()/1000, census.UpdateTime)
		census.OneByOneCancelCount = services.NewDynamicCensusService().GetOneByOneCount([]int{2, 3}, GetZero(now).UnixNano()/1000, census.UpdateTime)
		// 总数量
		census.OrderTotal = services.NewDynamicCensusService().GetOrderCount(0, 0, census.UpdateTime)
		census.OrderCancelTotal = services.NewDynamicCensusService().GetOrderCount(5, 0, census.UpdateTime)
		census.AppointTotal = services.NewDynamicCensusService().GetAppointCount(nil, 0, census.UpdateTime)
		census.AppointCancelTotal = services.NewDynamicCensusService().GetAppointCount([]int{3, 4, 5}, 0, census.UpdateTime)
		census.OneByOneTotal = services.NewDynamicCensusService().GetOneByOneCount(nil, 0, census.UpdateTime)
		census.OneByOneCancelTotal = services.NewDynamicCensusService().GetOneByOneCount([]int{2, 3}, 0, census.UpdateTime)
		logs.Infof(ctx, "UpsertDynamicCensus: %+v", census)
		err = services.NewDynamicCensusService().UpsertDynamicCensusCount(ctx, census)
		if err != nil {
			logs.Errorf(ctx, "UpsertDynamicCensusCount failed, date: %+v, err: %v", census, err)
			return
		}
		return
	}
	// 有存量数据处理
	census := *oldCensus
	logs.Infof(ctx, "dynamicCensus: %+v", census)
	// 本日数量
	census.OrderCount += services.NewDynamicCensusService().GetOrderCount(0, census.UpdateTime, now.UnixNano()/1000)
	census.OrderCancelCount += services.NewDynamicCensusService().GetOrderCount(5, census.UpdateTime, now.UnixNano()/1000)
	census.AppointCount += services.NewDynamicCensusService().GetAppointCount(nil, census.UpdateTime, now.UnixNano()/1000)
	census.AppointCancelCount += services.NewDynamicCensusService().GetAppointCount([]int{3, 4, 5}, census.UpdateTime, now.UnixNano()/1000)
	census.OneByOneCount += services.NewDynamicCensusService().GetOneByOneCount(nil, census.UpdateTime, now.UnixNano()/1000)
	census.OneByOneCancelCount += services.NewDynamicCensusService().GetOneByOneCount([]int{2, 3}, census.UpdateTime, now.UnixNano()/1000)
	// 总数量
	census.OrderTotal = services.NewDynamicCensusService().GetOrderCount(0, 0, now.UnixNano()/1000)
	census.OrderCancelTotal = services.NewDynamicCensusService().GetOrderCount(5, 0, now.UnixNano()/1000)
	census.AppointTotal = services.NewDynamicCensusService().GetAppointCount(nil, 0, now.UnixNano()/1000)
	census.AppointCancelTotal = services.NewDynamicCensusService().GetAppointCount([]int{3, 4, 5}, 0, now.UnixNano()/1000)
	census.OneByOneTotal = services.NewDynamicCensusService().GetOneByOneCount(nil, 0, now.UnixNano()/1000)
	census.OneByOneCancelTotal = services.NewDynamicCensusService().GetOneByOneCount([]int{2, 3}, 0, now.UnixNano()/1000)

	census.UpdateTime = now.UnixNano() / 1000
	logs.Infof(ctx, "UpsertDynamicCensus: %+v", census)
	err = services.NewDynamicCensusService().UpsertDynamicCensusCount(ctx, &census)
	if err != nil {
		logs.Errorf(ctx, "UpsertDynamicCensusCount failed, date: %+v, err: %v", census, err)
		return
	}
}

// GetBusinessCensus 业务数据统计
func (s *CensusService) GetBusinessCensus(ctx context.Context, now time.Time) {
	start := GetZero(now).UnixNano() / 1000000
	end := now.UnixNano() / 1000000
	// 获取当天的统计信息
	oldCensus, err := services.NewAService().GetBusinessOrderCensus(ctx, start, end)
	if err != nil {
		logs.Errorf(ctx, "GetBusinessOrderCensus failed, date: %v, err: %v", now.Format(time.RFC3339), err)
		return
	}
	census := &models.BusinessOrderCensus{
		CreateTime: util.Timestamp(),
		UpdateTime: util.Timestamp(),
	}
	if oldCensus != nil {
		census.Id = oldCensus.Id
		census.CreateTime = oldCensus.CreateTime
	}
	// 当天平台历史线上单订单总价
	census.OnlineOrderInAmount = services.NewDService().GetOnlineOrderAmount(ctx, 1, start, end)
	// 当天平台历史线下订单格总价
	census.OfflineOrderInAmount = services.NewDService().GetOfflineOrderAmount(ctx, 1, start, end)
	// 当天用户线上单总支付
	census.OnlineOrderOutAmount = services.NewDService().GetOnlineOrderAmount(ctx, 2, start, end)
	// 当天用户线下单总支付
	census.OfflineOrderOutAmount = services.NewDService().GetOfflineOrderAmount(ctx, 2, start, end)
	// 当天礼物总收入
	census.GiftInAmount = services.NewDService().GetAmountByType(ctx, 1, 2, start, end)
	// 当天送出礼物总消费（时币）
	census.GiftOutAmount = services.NewDService().GetAmountByType(ctx, 2, 2, start, end)

	// 平台历史线上单订单总价
	census.OnlineOrderInAmountTotal = services.NewDService().GetOnlineOrderAmount(ctx, 1, 0, 0)
	// 平台历史线下订单格总价
	census.OfflineOrderInAmountTotal = services.NewDService().GetOfflineOrderAmount(ctx, 1, 0, 0)
	// 用户线上单总支付
	census.OnlineOrderOutAmountTotal = services.NewDService().GetOnlineOrderAmount(ctx, 2, 0, 0)
	// 用户线下单总支付
	census.OfflineOrderOutAmountTotal = services.NewDService().GetOfflineOrderAmount(ctx, 2, 0, 0)
	// 礼物总收入
	census.GiftInAmountTotal = services.NewDService().GetAmountByType(ctx, 1, 2, 0, 0)
	// 送出礼物总消费（时币）
	census.GiftOutAmountTotal = services.NewDService().GetAmountByType(ctx, 2, 2, 0, 0)
	// 数据落库
	err = services.NewAService().UpsertBusinessOrderCensus(ctx, census)
	if err != nil {
		logs.Errorf(ctx, "UpsertBusinessOrderCensus failed, date: %+v, err: %v", census, err)
		return
	}
}

func GetZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func InSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
