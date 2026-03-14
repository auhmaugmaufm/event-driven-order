package dto

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementRequest struct {
	ProductID    uuid.UUID `json:"product_id"`
	MovementType string    `json:"movement_type"`
	Quantity     int       `json:"quantity"`
}

type StockMovementResponse struct {
	ID           uuid.UUID `json:"stock_id"`
	ProductID    uuid.UUID `json:"product_id"`
	MovementType string    `json:"movement_type"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
