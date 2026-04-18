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

func (s *stockMovementReposity) CreateBulk(ctx context.Context, stockMovements []domain.StockMovement) error {
	return s.db.WithContext(ctx).CreateInBatches(stockMovements, 100).Error
}

func (s *stockMovementReposity) FindByMovementType(ctx context.Context, movementType string, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit
	var stockMovements []domain.StockMovement
	err := s.db.WithContext(ctx).Where("movement_type = ?", movementType).Limit(limit).Offset(offset).Preload("Stock").Preload("Product").Find(&stockMovements).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := s.db.WithContext(ctx).Model(&domain.StockMovement{}).Where("movement_type = ?", movementType).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	return stockMovements, totalOrder, nil
}

func (s *stockMovementReposity) FindByStockMovementID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	var stockMovement domain.StockMovement
	err := s.db.WithContext(ctx).Preload("Stock").Preload("Stock.Product").Where("id = ?", id).First(&stockMovement).Error
	if err != nil {
		return nil, err
	}
	return &stockMovement, nil
}

func (s *stockMovementReposity) GetStockMovement(ctx context.Context, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	var stockMovements []domain.StockMovement
	err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Preload("Stock").Preload("Stock.Product").Find(&stockMovements).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := s.db.WithContext(ctx).Model(&domain.StockMovement{}).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	return stockMovements, totalOrder, nil
}
