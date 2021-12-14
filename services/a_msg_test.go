package services

import (
	"context"
	"test/models"
	"test/util"
	"testing"
)

var ctx = context.TODO()

func TestAddCommentCensusLog(t *testing.T) {
	log := &models.CommentCensusLog{
		CommentId:   1,
		CommentUid:  1,
		CommentType: 1,
		Status:      1,
		CreateTime:  util.Timestamp(),
	}

	err := NewAService().AddCommentCensusLog(ctx, log)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("ok...")

	err = NewAService().DelCommentCensusLog(ctx, log.CommentId, 0)
	if err != nil {
		t.Fatalf("del err: %v", err)
	}
	t.Logf("del ok...")
}

func TestAddMomentCensusLog(t *testing.T) {
	log := &models.MomentCensusLog{
		MomentId:   1,
		MomentUid:  1,
		MomentType: 1,
		Status:     1,
		HavePhoto:  1,
		DataType:   1,
		CreateTime: util.Timestamp(),
	}

	err := NewAService().AddMomentCensusLog(ctx, log)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("ok...")

	err = NewAService().DelMomentCensusLog(ctx, log.MomentId, 0)
	if err != nil {
		t.Fatalf("del err: %v", err)
	}
	t.Logf("del ok...")
}
