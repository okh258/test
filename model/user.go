package model

type UserDevice struct {
	DeviceToken   string `json:"device_token"`
	DeviceName    string `json:"device_name"`
	LastUserId    int64  `json:"last_user_id"`
	LastLoginTime int64  `json:"last_login_time"`
}
