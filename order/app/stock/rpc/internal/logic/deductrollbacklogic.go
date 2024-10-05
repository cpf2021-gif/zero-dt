package logic

import (
	"context"
	"database/sql"
	"fmt"

	"zero-dt/order/app/stock/rpc/internal/svc"
	"zero-dt/order/app/stock/rpc/pb"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductRollbackLogic {
	return &DeductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductRollbackLogic) DeductRollback(in *pb.DeductReq) (*pb.DeductResp, error) {
	l.Logger.Infof("扣减库存回滚 in: %+v", in)

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		if err := l.svcCtx.StockModel.AddStock(tx, in.GoodsId, in.Num); err != nil {
			return fmt.Errorf("扣减库存回滚失败: err: %v, in: %+v", err, in)
		}
		return nil
	}); err != nil {
		l.Logger.Errorf("扣减库存回滚失败: err: %v, in: %+v", err, in)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeductResp{}, nil
}
