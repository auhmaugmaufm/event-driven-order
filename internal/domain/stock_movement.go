package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StockMovement struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StockID      uuid.UUID `json:"stock_id" gorm:"type:uuid;not null;index"`
	MovementType string    `json:"movement_type" gorm:"not null"`
	Quantity     int       `json:"quantity" gorm:"not null;default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreatetTme"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Stock *Stock `json:"stock,omitempty" gorm:"foreignKey:StockID"`
}

type StockMovementRepository interface {
	Create(ctx context.Context, stockMovement *StockMovement) error
	GetStockMovement(ctx context.Context) ([]StockMovement, error)
	FindByStockMovementID(ctx context.Context, id uuid.UUID) (*StockMovement, error)
	FindByMovementType(ctx context.Context, movementType string) ([]StockMovement, error)
}
