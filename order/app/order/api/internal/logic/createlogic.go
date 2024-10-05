package logic

import (
	"context"
	"fmt"

	"zero-dt/order/app/order/api/internal/svc"
	"zero-dt/order/app/order/api/internal/types"
	"zero-dt/order/app/order/rpc/order"
	orderpb "zero-dt/order/app/order/rpc/pb"
	stockpb "zero-dt/order/app/stock/rpc/pb"
	"zero-dt/order/app/stock/rpc/stock"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

var dtmServer = "etcd://localhost:2379/dtmservice"

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.QuickCreateReq) (resp *types.QuickCreateResp, err error) {
	orderRpcBusiServer, err := l.svcCtx.Config.OrderRpcConf.BuildTarget()
	if err != nil {
		return nil, fmt.Errorf("下单异常超时")
	}

	stockRpcBusiServer, err := l.svcCtx.Config.StockRpcConf.BuildTarget()
	if err != nil {
		return nil, fmt.Errorf("下单异常超时")
	}

	createOrderReq := &order.CreateReq{UserId: req.UserId, GoodsId: req.GoodsId, Num: req.Num}
	deductReq := &stock.DeductReq{GoodsId: req.GoodsId, Num: req.Num}

	gid := dtmgrpc.MustGenGid(dtmServer)
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(orderRpcBusiServer+orderpb.Order_Create_FullMethodName, orderRpcBusiServer+orderpb.Order_CreateRollback_FullMethodName, createOrderReq).
		Add(stockRpcBusiServer+stockpb.Stock_Deduct_FullMethodName, stockRpcBusiServer+stockpb.Stock_DeductRollback_FullMethodName, deductReq)

	saga.WaitResult = true
	err = saga.Submit()

	if err != nil {
		return nil, fmt.Errorf("submit data to dtm-server failed: %v", err)
	}

	return &types.QuickCreateResp{}, err
}
