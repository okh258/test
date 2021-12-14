package services

import (
	"context"
	"git.devops.com/wsim/hflib/logs"
	"github.com/astaxie/beego/orm"
	"test/models"
)

// AddCommentCensusLog 新增评论统计记录
func (s *AService) AddCommentCensusLog(ctx context.Context, log *models.CommentCensusLog) error {
	orm.RegisterModel(&models.CommentCensusLog{})
	_, err := s.o.Insert(log)
	if err != nil {
		logs.Errorf(ctx, "AddCommentCensusLog err:", err)
		return err
	}
	return nil
}

func (s *AService) DelCommentCensusLog(ctx context.Context, id, status int64) error {
	sql := "UPDATE `t_comment_census_log` SET `status` = ? WHERE `comment_id` = ?"
	_, err := s.o.Raw(sql, status, id).Exec()
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "DelCommentCensusLog failed, err: %v", err)
		return err
	}
	return nil
}

// AddMomentCensusLog 新增动态统计记录
func (s *AService) AddMomentCensusLog(ctx context.Context, log *models.MomentCensusLog) error {
	orm.RegisterModel(&models.MomentCensusLog{})
	_, err := s.o.Insert(log)
	if err != nil {
		logs.Errorf(ctx, "AddMomentCensusLog err:", err)
		return err
	}
	return nil
}

func (s *AService) DelMomentCensusLog(ctx context.Context, id, status int64) error {
	sql := "UPDATE `t_moment_census_log` SET `status` = ? WHERE `moment_id` = ?"
	_, err := s.o.Raw(sql, status, id).Exec()
	if err != nil && err != orm.ErrNoRows {
		logs.Errorf(ctx, "DelMomentCensusLog failed, err: %v", err)
		return err
	}
	return nil
}
