package svc

import (
	"zero-dt/order/app/order/api/internal/config"
	"zero-dt/order/app/order/rpc/order"
	"zero-dt/order/app/stock/rpc/stock"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc order.Order
	StockRpc stock.Stock
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc: order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		StockRpc: stock.NewStock(zrpc.MustNewClient(c.StockRpcConf)),
	}
}
