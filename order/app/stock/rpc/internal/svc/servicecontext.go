package svc

import (
	"zero-dt/order/app/stock/rpc/internal/config"
	"zero-dt/order/app/stock/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	StockModel model.StockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		StockModel: model.NewStockModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
