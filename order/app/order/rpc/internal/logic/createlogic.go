package logic

import (
	"context"
	"database/sql"
	"fmt"

	"zero-dt/order/app/order/rpc/internal/model"
	"zero-dt/order/app/order/rpc/internal/svc"
	"zero-dt/order/app/order/rpc/pb"

	"github.com/dtm-labs/client/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *pb.CreateReq) (*pb.CreateResp, error) {
	l.Logger.Infof("创建订单 in: %+v", in)

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		order := new(model.Order)
		order.GoodsId = in.GoodsId
		order.Num = in.Num
		order.UserId = in.UserId

		_, err := l.svcCtx.OrderModel.InsertWithTx(tx, order)
		if err != nil {
			return fmt.Errorf("创建订单失败: err: %v, order: %+v", err, order)
		}

		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateResp{}, nil
}
