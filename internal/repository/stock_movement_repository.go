package repository

import (
	"context"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockMovementReposity struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) domain.StockMovementRepository {
	return &stockMovementReposity{db: db}
}

func NewStockMovementRepositoryWithTx(tx *gorm.DB) domain.StockMovementRepository {
	return &stockMovementReposity{db: tx}
}

func (s *stockMovementReposity) Create(ctx context.Context, stockMovement *domain.StockMovement) error {
	return s.db.WithContext(ctx).Create(stockMovement).Error
}

func (s *stockMovementReposity) FindByMovementType(ctx context.Context, movementType string) ([]domain.StockMovement, error) {
	var stockMovements []domain.StockMovement
	err := s.db.WithContext(ctx).Preload("Stock").Preload("Product").Where("movement_type = ?", movementType).Find(&stockMovements).Error
	if err != nil {
		return nil, err
	}
	return stockMovements, nil
}

func (s *stockMovementReposity) FindByStockMovementID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	var stockMovement domain.StockMovement
	err := s.db.WithContext(ctx).Preload("Stock").Preload("Product").Where("id = ?", id).First(&stockMovement).Error
	if err != nil {
		return nil, err
	}
	return &stockMovement, nil
}

func (s *stockMovementReposity) GetStockMovement(ctx context.Context) ([]domain.StockMovement, error) {
	var stockMovements []domain.StockMovement
	err := s.db.WithContext(ctx).Preload("Stock").Preload("Product").Find(&stockMovements).Error
	if err != nil {
		return nil, err
	}
	return stockMovements, nil
}
