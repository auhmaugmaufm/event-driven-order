package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) domain.StockRepository {
	return &stockRepository{db: db}
}

func NewStockRepositoryWithTx(tx *gorm.DB) domain.StockRepository {
	return &stockRepository{db: tx}
}

func (r *stockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

func (r *stockRepository) IncreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

func (r *stockRepository) DecreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (r *stockRepository) DecreaseStockBulkWithTx(ctx context.Context, stockAdjustments []domain.StockAdjustment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, adj := range stockAdjustments {
			result := tx.Model(&domain.Stock{}).
				Where("product_id = ? AND quantity >= ?", adj.ProductID, adj.Quantity).
				Update("quantity", gorm.Expr("quantity - ?", adj.Quantity))

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (r *stockRepository) GetProductStock(ctx context.Context, productId uuid.UUID) (*domain.Stock, error) {
	var stock domain.Stock
	err := r.db.WithContext(ctx).Preload("Product").Where("product_id = ?", productId).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) GetStocks(ctx context.Context, pagination *domain.Pagination) ([]domain.Stock, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	var stocks []domain.Stock
	err := r.db.WithContext(ctx).Limit(pagination.Limit).Offset(offset).Preload("Product").Find(&stocks).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.Stock{}).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	return stocks, totalOrder, nil
}
