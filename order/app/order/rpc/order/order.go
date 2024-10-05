// Code generated by goctl. DO NOT EDIT.
// Source: order.proto

package order

import (
	"context"

	"zero-dt/order/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateReq  = pb.CreateReq
	CreateResp = pb.CreateResp

	Order interface {
		Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error)
		CreateRollback(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}

func (m *defaultOrder) Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.Create(ctx, in, opts...)
}

func (m *defaultOrder) CreateRollback(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.CreateRollback(ctx, in, opts...)
}
