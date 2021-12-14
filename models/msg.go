package models

type CommentCensusLog struct {
	CommentId   int64 `json:"comment_id" orm:"pk"` // 评论ID comment_id
	CommentUid  int64 `json:"comment_uid"`         // 评论用户 uid
	CommentType int64 `json:"comment_type"`        // 评论类型: 0：文本 1：语音 2：图片 3：视频 4：涂鸦  5:表情贴纸
	Status      int64 `json:"status"`              // 状态: 0 正常 -1 删除
	CreateTime  int64 `json:"create_time"`         // 统计记录创建时间
}

func (m *CommentCensusLog) TableName() string {
	return "t_comment_census_log"
}

type MomentCensusLog struct {
	MomentId   int64 `json:"moment_id" orm:"pk"` // 动态ID moment_id
	MomentUid  int64 `json:"moment_uid"`         // 动态用户 uid
	MomentType int64 `json:"moment_type"`        // 动态类型 0：普通动态， 1：只绑话题的动态， 2：只绑贴吧的动态，3：有绑话题并且绑贴吧
	Status     int64 `json:"status"`             // 状态: 0 正常显示, 1 隐藏， -1 删除
	HavePhoto  int64 `json:"have_photo"`         // 动态是否含照片(0：无，1：有）
	DataType   int64 `json:"data_type"`          // 动态数据类型 0：随便发（原写两句），1：说照片，2:涂一涂（原涂鸦），3：上传照片, 4:文字个性签名，5：语音个性签名, 6：分享链接，7：我的状态（原心情动态），8：随手拍，9：发照片，10：上传文件，11：秘密记事本、12：写小记，13：分享写笔记、14：说故事、15：分享音乐、16：分享录音乐、17:签证视频
	CreateTime int64 `json:"create_time"`        // 统计记录创建时间
}

func (m *MomentCensusLog) TableName() string {
	return "t_moment_census_log"
}
