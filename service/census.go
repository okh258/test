package service

import (
	"context"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	"test/db"
	"test/model"
	"time"
)

const (
	GoodsId_Gold             = 1
	GoodsId_VIP              = 2
	GoodsId_SVIP             = 3
	GoodsId_MoneyPower       = 4
	GoodsId_DatingOrder      = 5
	GoodsId_DatingOrderAgain = 6
	GoodsId_Wallet           = 7
)

var _censusService *CensusService

func NewCensusService() *CensusService {
	if _censusService == nil {
		_censusService = &CensusService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameAdmin),
		}
		orm.RegisterModel(&model.CapitalInflowCensus{}, &model.CapitalOutflowCensus{})
	}
	return _censusService
}

type CensusService struct {
	o orm.Ormer
}

// GetCapitalOutflowCensus 根据时间获取1条最新的统计信息, 时间遵循左闭右开
// fromType 1：服务收益，2：邀请收益
func (s *CensusService) GetCapitalOutflowCensus(ctx context.Context, fromType, start, end int64) (*model.CapitalOutflowCensus, error) {
	var result *model.CapitalOutflowCensus
	sql := "SELECT * FROM t_capital_outflow_census WHERE create_time >= ? AND create_time < ? ORDER BY create_time DESC LIMIT 1"
	r := s.o.Raw(sql, start, end)
	if fromType > 0 {
		sql = "SELECT * FROM t_capital_outflow_census WHERE create_time >= ? AND create_time < ? AND from_type = ? ORDER BY create_time DESC LIMIT 1"
		r.SetArgs(fromType)
	}

	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "getCapitalOutflowCensus err:", err)
		return nil, err
	}
	return result, nil
}

func (s *CensusService) UpsertCapitalFlowCensus(data interface{}) error {
	if data == nil {
		return nil
	}
	o := s.o
	_, err := o.InsertOrUpdate(data)
	if err != nil {
		logs.Errorf(context.TODO(), "UpsertCapitalFlowCensus err:", err)
		return err
	}
	return nil
}

func (s *CensusService) InsertCapitalOutflowCensus(data *model.CapitalOutflowCensus) error {
	_, err := s.o.Insert(data)
	if err != nil {
		logs.Errorf(context.TODO(), "InsertCapitalOutflowCensus err:", err)
		return err
	}
	return nil
}

func GetZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetCapitalInflowCensus 根据时间获取1条最新统计的信息, 时间遵循左闭右开
// 支付类型，1：支付宝，2：微信，3：苹果内购
func (s *CensusService) GetCapitalInflowCensus(ctx context.Context, payType, start, end int64) (*model.CapitalInflowCensus, error) {
	o := s.o
	var result *model.CapitalInflowCensus
	sql := "SELECT * FROM t_capital_inflow_census WHERE create_time >= ? AND create_time < ? ORDER BY create_time DESC LIMIT 1"
	r := o.Raw(sql, start, end)
	if payType > 0 {
		sql = "SELECT * FROM t_capital_inflow_census WHERE create_time >= ? AND create_time < ? AND pay_type = ? ORDER BY create_time DESC LIMIT 1"
		r = o.Raw(sql, start, end, payType)
	}

	err := r.QueryRow(&result)
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "GetCapitalInflowCensus err:", err)
		return nil, err
	}
	return result, nil
}

func GetOldCapitalInflowCensus(ctx context.Context, payType, zeroStam, nextZeroStam int64, zero time.Time) (*model.CapitalInflowCensus, error) {
	// 获取今天存量数据
	old, err := NewCensusService().GetCapitalInflowCensus(ctx, payType, zeroStam, nextZeroStam)
	if err != nil {
		logs.Errorf(ctx, "GetCapitalInflowCensus failed, err: %v", err)
		return nil, err
	}
	// 获取昨天存量数据
	if old == nil {
		yesterday := zero.Add(-24*time.Hour).UnixNano() / 100000 // 昨天零点
		// 获取今天存量数据
		old, err = NewCensusService().GetCapitalInflowCensus(ctx, payType, yesterday, zeroStam)
		if err != nil {
			logs.Errorf(ctx, "GetCapitalInflowCensus failed, err: %v", err)
			return nil, err
		}
	}
	return old, nil
}
