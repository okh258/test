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
