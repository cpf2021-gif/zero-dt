package model

import (
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		withSession(session sqlx.Session) OrderModel

		FindLastOneByUserIdGoodsId(userId, goodsId int64) (*Order, error)
		InsertWithTx(tx *sql.Tx, data *Order) (sql.Result, error)
		UpdateWithTx(tx *sql.Tx, data *Order) error
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn),
	}
}

func (m *customOrderModel) withSession(session sqlx.Session) OrderModel {
	return NewOrderModel(sqlx.NewSqlConnFromSession(session))
}

func (c *customOrderModel) FindLastOneByUserIdGoodsId(userId int64, goodsId int64) (*Order, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `goods_id` = ? order by `id` desc limit 1", orderRows, c.table)
	var resp Order
	err := c.conn.QueryRow(&resp, query, userId, goodsId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customOrderModel) InsertWithTx(tx *sql.Tx, data *Order) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", c.table, orderRowsExpectAutoSet)
	return tx.Exec(query, data.UserId, data.GoodsId, data.Num, data.RowState)
}

func (c *customOrderModel) UpdateWithTx(tx *sql.Tx, data *Order) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", c.table, orderRowsWithPlaceHolder)
	_, err := tx.Exec(query, data.UserId, data.GoodsId, data.Num, data.RowState, data.Id)
	return err
}
