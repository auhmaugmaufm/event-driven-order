package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"gorm.io/gorm"
)

type txRepository struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) domain.TxRepository {
	return &txRepository{db: db}
}

func (t *txRepository) ExecTx(ctx context.Context, fn func(
	productRepo domain.ProductRepository,
	stockRepo domain.StockRepository,
) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(
			&productRepository{db: tx},
			&stockRepository{db: tx},
		)
	})
}

func (t *txRepository) ExecStockMovementTx(ctx context.Context, fn func(
	stockMovementRepo domain.StockMovementRepository,
	stockRepo domain.StockRepository,
) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(
			&stockMovementReposity{db: tx},
			&stockRepository{db: tx},
		)
	})
}
