package svc

import (
	"zero-dt/order/app/order/rpc/internal/config"
	"zero-dt/order/app/order/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
