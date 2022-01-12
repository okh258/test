package services

import (
	"test/models"
	"testing"
)

func TestGetAllCalls(t *testing.T) {

	var params = models.CallVideoOrderReq{
		CoachUid:      0,
		CoachNickName: "",
		UserId:        0,
		UserNickName:  "",
		CoachUidList:  nil,
		UserIdList:    nil,
		Status:        -1,
		OrderId:       0,
		StartTime:     0,
		EndTime:       0,
		PageNum:       45,
		PageSize:      20,
	}
	list, err := NewDService().GetAllCalls(ctx, params)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("len: %v, list: %+v", len(list), list)
}

func TestGetSignatureCount(t *testing.T) {
	list, err := NewDService().GetSignatureCount(ctx, 11, 1, 356503167878144)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("len: %v", list)
	list -= 1
	if list > 2 {
		t.Logf("true")
		return
	}
	t.Logf("false")
}

func TestGetCategoryListByIds(t *testing.T) {
	list, err := NewDService().GetCategoryListByIds(ctx, []int64{34, 11})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("len: %v, list: %+v", len(list), list)
}

func TestPageRechargeGoldLog(t *testing.T) {
	p := &models.RechargeInviteReq{
		InviteeUid: 0,
		InviterUid: 0,
		Status:     0,
		PageNum:    0,
		PageSize:   0,
	}

	list, total, err := NewDService().PageRechargeGoldLog(ctx, p)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("len: %v, list: %+v", total, list)
}
