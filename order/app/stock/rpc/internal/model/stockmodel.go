package model

import (
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StockModel = (*customStockModel)(nil)

type (
	// StockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockModel.
	StockModel interface {
		stockModel
		withSession(session sqlx.Session) StockModel

		DecuctStock(tx *sql.Tx, goodsId, num int64) (sql.Result, error)
		AddStock(tx *sql.Tx, goodsId, num int64) error
	}

	customStockModel struct {
		*defaultStockModel
	}
)

// NewStockModel returns a model for the database table.
func NewStockModel(conn sqlx.SqlConn) StockModel {
	return &customStockModel{
		defaultStockModel: newStockModel(conn),
	}
}

func (m *customStockModel) withSession(session sqlx.Session) StockModel {
	return NewStockModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customStockModel) DecuctStock(tx *sql.Tx, goodsId, num int64) (sql.Result, error) {
	query := fmt.Sprintf("update %s set `num` = `num` - ? where `goods_id` = ? and num >= ?", m.table)
	return tx.Exec(query, num, goodsId, num)
}

func (m *customStockModel) AddStock(tx *sql.Tx, goodsId, num int64) error {
	query := fmt.Sprintf("update %s set `num` = `num` + ? where `goods_id` = ?", m.table)
	_, err := tx.Exec(query, num, goodsId)
	return err
}