package domain

import (
	"context"
)

type TxRepository interface {
	ExecTx(ctx context.Context, fn func(
		productRepo ProductRepository,
		stockRepo StockRepository,
	) error) error
	ExecStockMovementTx(ctx context.Context, fn func(
		stockMovement StockMovementRepository,
		stockRepo StockRepository,
	) error) error
	// ExecStockMovementBulkTx(ctx context.Context, fn func(
	// 	stockMovement StockMovementRepository,
	// 	stockRepo StockRepository,
	// ) error) error
	ExecOrderTx(ctx context.Context, fn func(
		orderRepo OrderRepository,
		stockMovement StockMovementRepository,
		stockRepo StockRepository,
	) error) error
}
