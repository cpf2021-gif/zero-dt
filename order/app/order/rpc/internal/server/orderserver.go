// Code generated by goctl. DO NOT EDIT.
// Source: order.proto

package server

import (
	"context"

	"zero-dt/order/app/order/rpc/internal/logic"
	"zero-dt/order/app/order/rpc/internal/svc"
	"zero-dt/order/app/order/rpc/pb"
)

type OrderServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedOrderServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
	}
}

func (s *OrderServer) Create(ctx context.Context, in *pb.CreateReq) (*pb.CreateResp, error) {
	l := logic.NewCreateLogic(ctx, s.svcCtx)
	return l.Create(in)
}

func (s *OrderServer) CreateRollback(ctx context.Context, in *pb.CreateReq) (*pb.CreateResp, error) {
	l := logic.NewCreateRollbackLogic(ctx, s.svcCtx)
	return l.CreateRollback(in)
}