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

type CreateRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRollbackLogic {
	return &CreateRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRollbackLogic) CreateRollback(in *pb.CreateReq) (*pb.CreateResp, error) {
	l.Logger.Infof("创建订单回滚 in: %+v", in)

	order, err := l.svcCtx.OrderModel.FindLastOneByUserIdGoodsId(in.UserId, in.GoodsId)

	l.Logger.Infof("创建订单回滚 order: %+v, err: %v", order, err)

	if err != nil && err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, err.Error())
	}

	
	if order != nil {
		l.Logger.Infof("需要回滚的订单: %+v", order)

		barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
			order.RowState = -1
			if err := l.svcCtx.OrderModel.UpdateWithTx(tx, order); err != nil {
				return fmt.Errorf("创建订单回滚失败: err: %v, order: %+v", err, order)
			}

			return nil
		}); err != nil {
			l.Logger.Errorf("创建订单回滚失败: err: %v, order: %+v", err, order)

			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.CreateResp{}, nil
}
