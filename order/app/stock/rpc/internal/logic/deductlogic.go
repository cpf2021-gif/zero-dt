package logic

import (
	"context"
	"database/sql"

	"zero-dt/order/app/stock/rpc/internal/model"
	"zero-dt/order/app/stock/rpc/internal/svc"
	"zero-dt/order/app/stock/rpc/pb"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductLogic {
	return &DeductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductLogic) Deduct(in *pb.DeductReq) (*pb.DeductResp, error) {
	l.Logger.Infof("扣减库存 in: %+v", in)

	stock, err := l.svcCtx.StockModel.FindOneByGoodsId(l.ctx, in.GoodsId)
	if err != nil && err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if stock == nil || stock.Num < in.Num {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		sqlResult, err := l.svcCtx.StockModel.DecuctStock(tx, in.GoodsId, in.Num)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		affected, err := sqlResult.RowsAffected()
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if affected <= 0 {
			return status.Error(codes.Aborted, dtmcli.ResultFailure)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.DeductResp{}, nil
}
