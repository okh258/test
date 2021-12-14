package services

import (
	"context"
	"github.com/astaxie/beego/orm"
	"test/db"
)

var _vipEventService *VipEventService

func NewVipEventService() *VipEventService {
	if _vipEventService == nil {
		_vipEventService = &VipEventService{
			o: db.NewOrmWithDB(context.Background(), db.AliasNameAdmin),
		}
	}
	return _vipEventService
}

type VipEventService struct {
	o orm.Ormer
}
