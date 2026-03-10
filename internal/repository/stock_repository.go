package repository

import (
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

func (r *stockRepository) CreateWithTx(tx *gorm.DB, stock *domain.Stock) error {
	return tx.Create(stock).Error
}

func (r *stockRepository) IncreaseStock(productId uuid.UUID, quantity int) error {
	return r.db.Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

func (r *stockRepository) DecreaseStock(productId uuid.UUID, quantity int) error {
	return r.db.Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (r *stockRepository) GetProductStock(productId uuid.UUID) (*domain.Stock, error) {
	var stock domain.Stock
	err := r.db.Preload("Product").Where("product_id = ?", productId).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) GetStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Preload("Product").Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
