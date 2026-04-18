package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type StockAdjustment struct {
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:0"`
}

type StockRepository interface {
	Create(ctx context.Context, stock *Stock) error
	IncreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error
	DecreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error
	DecreaseStockBulkWithTx(ctx context.Context, stockAdjustments []StockAdjustment) error
	GetProductStock(ctx context.Context, productId uuid.UUID) (*Stock, error)
	GetStocks(ctx context.Context, pagination *Pagination) ([]Stock, int64, error)
}
