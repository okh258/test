package service

import (
	"context"
	"git.devops.com/go/golib/util"
	"git.devops.com/wsim/hflib/logs"
	"test/model"
	"testing"
	"time"
)

func TestAddCapitalOutflowCensus(t *testing.T) {
	ctx := context.TODO()

	zero := GetZero(time.Now())                                // 今天零点
	zeroStam := zero.UnixNano() / 100000                       // 今天零点
	nextZeroStam := zero.Add(24*time.Hour).UnixNano() / 100000 // 明天零点

	old, err := NewCensusService().GetCapitalOutflowCensus(ctx, 1, zeroStam, nextZeroStam)
	if err != nil {
		logs.Errorf(ctx, "GetCapitalOutflowCensus failed, err: %v", err)
		return
	}
	if old != nil {
		t.Logf("old != nil")
	}
	t.Logf("old == nil")

	d := &model.CapitalOutflowCensus{
		FromType:    1,
		Amount:      10,
		AmountTotal: 100,
		CreateTime:  util.Timestamp(),
		UpdateTime:  util.Timestamp(),
	}
	err = NewCensusService().InsertCapitalOutflowCensus(d)
	if err != nil {
		logs.Errorf(ctx, "UpsertCapitalOutflowCensus failed, err: %v", err)
	}
	d.Amount = 11
	err = NewCensusService().UpsertCapitalFlowCensus(d)
	if err != nil {
		logs.Errorf(ctx, "UpsertCapitalOutflowCensus failed, err: %v", err)
	}
}

func TestAddCapitalInflowCensus(t *testing.T) {
	ctx := context.TODO()
	info := model.Order{
		GoodsId:    1,
		PayType:    3,
		PayAmount:  3000,
		RealAmount: 3000,
	}

	logs.Infof(ctx, "AddCapitalInflowCensus orderInfo: %+v", info)
	zero := GetZero(time.Now())                                 // 今天零点
	zeroStam := zero.UnixNano() / 1000000                       // 今天零点
	nextZeroStam := zero.Add(24*time.Hour).UnixNano() / 1000000 // 明天零点
	old, err := GetOldCapitalInflowCensus(ctx, info.PayType, zeroStam, nextZeroStam, zero)
	if err != nil {
		logs.Errorf(ctx, "getOldCapitalOutflowCensus failed, err: %v", err)
		return
	}
	new := &model.CapitalInflowCensus{
		PayType:    info.PayType,
		CreateTime: util.Timestamp(),
		UpdateTime: util.Timestamp(),
	}
	// 1: 金币 2: 钱包
	switch info.GoodsId {
	case GoodsId_Gold:
		new.GoldAmount = info.RealAmount
		new.GoldAmountTotal = info.RealAmount
		new.GoldPayAmount = info.PayAmount
		new.GoldPayAmountTotal = info.PayAmount
	case GoodsId_Wallet:
		new.WalletAmount = info.RealAmount
		new.WalletAmountTotal = info.RealAmount
		new.WalletPayAmount = info.PayAmount
		new.WalletPayAmountTotal = info.PayAmount
	}
	// 当天第一笔, 重置当日金额
	if old != nil && !(zeroStam <= old.CreateTime && old.CreateTime < nextZeroStam) {
		old.GoldAmount = 0
		old.GoldPayAmount = 0
		old.WalletAmount = 0
		old.WalletPayAmount = 0
	}
	if old != nil {
		new.GoldAmount += old.GoldAmount
		new.GoldPayAmount += old.GoldPayAmount
		new.WalletAmount += old.WalletAmount
		new.WalletPayAmount += old.WalletPayAmount
		new.GoldAmountTotal += old.GoldAmountTotal
		new.GoldPayAmountTotal += old.GoldPayAmountTotal
		new.WalletAmountTotal += old.WalletAmountTotal
		new.WalletPayAmountTotal += old.WalletPayAmountTotal
		new.CreateTime = old.CreateTime
	}
	err = NewCensusService().UpsertCapitalFlowCensus(new)
	if err != nil {
		logs.Errorf(ctx, "UpsertCapitalInFlowCensus failed, err: %v", err)
	}
}
