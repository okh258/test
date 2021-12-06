package base

import (
	"reflect"
	"testing"
)

// 比较结构体的值是否相等
func TestEqual(t *testing.T) {
	d := &DynamicCount{LockOrderNum: 1, SkillNum: 10, CensorNickNameNum: 100}
	b := &DynamicCount{LockOrderNum: 1, SkillNum: 10, CensorNickNameNum: 100}
	t.Logf("\nold dynamic: %+v, \nnew dynamic: %+v", d, b)
	if reflect.DeepEqual(d, b) {
		t.Log("结构体值相等")
		goto a1
	}
	t.Log("结构体值不相等")

a1:
	x := *d
	x.LockOrderNum = 2
	if reflect.DeepEqual(x, d) {
		t.Log("结构体值相等")
		return
	}
	t.Log("结构体值不相等")
}

type DynamicCount struct {
	LockOrderNum         int64 `json:"lockOrderNum"`
	SkillNum             int64 `json:"skillNum"`
	RealNameAuthNum      int64 `json:"realNameAuthNum"`
	CensorMomentNum      int64 `json:"censorMomentNum"`
	CensorCommentNum     int64 `json:"censorCommentNum"`
	CensorMsgNum         int64 `json:"censorMsgNum"`
	CensorCensorVideoNum int64 `json:"censorCensorVideoNum"`
	CensorNickNameNum    int64 `json:"censorNickNameNum"`
	CensorPortraitNum    int64 `json:"censorPortraitNum"`
	CensorSignatureNum   int64 `json:"censorSignatureNum"`
	CensorFeedbackNum    int64 `json:"censorFeedbackNum"`
	CensorTipOffNum      int64 `json:"censorTipOffNum"`
}
