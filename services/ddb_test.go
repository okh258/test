package services

import (
	"context"
	"testing"
)

func TestGetUserDriverInfo(t *testing.T) {
	info, err := GetUserDriverMap(context.TODO(), 530671533814784)
	if err != nil {
		t.Fatalf("get info failed, err: %v", info)
		return
	}
	t.Logf("info: %+v\n", info)
	info1, err := GetUserDriver(context.TODO(), "687efc393f500436")
	if err != nil {
		t.Fatalf("get info failed, err: %v", info)
		return
	}
	t.Logf("info: %+v\n", info1)
	info2, err := GetUserDriverItem(context.TODO(), 530671533814784)
	if err != nil {
		t.Fatalf("get info failed, err: %v", info)
		return
	}
	t.Logf("info: %+v\n", info2)
}

func TestGetDetailUsers(t *testing.T) {
	uids := []int64{500944337634304, 530671533814784}
	ids := make([]interface{}, len(uids))
	for index, value := range uids {
		ids[index] = value
	}

	count, err := GetDetailUsers(context.TODO(), ids)
	if err != nil {
		t.Logf("TestGetDetailUsers failed, err: %v", err)
		return
	}
	t.Logf("count: %v", count)
}
