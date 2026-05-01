package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type stockMovementReposity struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewStockMovementRepository(db *gorm.DB, redis *redis.Client) domain.StockMovementRepository {
	return &stockMovementReposity{db: db, redis: redis}
}

func NewStockMovementRepositoryWithTx(tx *gorm.DB) domain.StockMovementRepository {
	return &stockMovementReposity{db: tx}
}

func (r *stockMovementReposity) Create(ctx context.Context, stockMovement *domain.StockMovement) error {
	err := r.db.WithContext(ctx).Create(stockMovement).Error
	if err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "orders:*")
}

func (r *stockMovementReposity) CreateBulk(ctx context.Context, stockMovements []domain.StockMovement) error {
	return r.db.WithContext(ctx).CreateInBatches(stockMovements, 100).Error
}

func (r *stockMovementReposity) FindByMovementType(ctx context.Context, movementType string, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit
	var stockMovements []domain.StockMovement
	err := r.db.WithContext(ctx).Where("movement_type = ?", movementType).Limit(limit).Offset(offset).Preload("Stock").Preload("Product").Find(&stockMovements).Error
	if err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.StockMovement{}).Where("movement_type = ?", movementType).Count(&totalOrder).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	return stockMovements, totalOrder, nil
}

func (r *stockMovementReposity) FindByStockMovementID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	var stockMovement domain.StockMovement
	err := r.db.WithContext(ctx).Preload("Stock").Preload("Stock.Product").Where("id = ?", id).First(&stockMovement).Error
	if err != nil {
		return nil, err
	}
	return &stockMovement, nil
}

func (r *stockMovementReposity) GetStockMovement(ctx context.Context, pagination *domain.Pagination) ([]domain.StockMovement, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	cachaKey := fmt.Sprintf("stockmoves:page:%d:limit:%d", page, limit)
	totalKey := "stockmoves:total"

	cached, err := r.redis.Get(ctx, totalKey).Result()
	if err == nil {
		var stockMovements []domain.StockMovement
		if err := json.Unmarshal([]byte(cached), &stockMovements); err == nil {
			totalStr, err := r.redis.Get(ctx, totalKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(totalStr, 10, 64)
				return stockMovements, total, nil
			}
		}
	}

	var stockMovements []domain.StockMovement
	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Preload("Stock").Preload("Stock.Product").Find(&stockMovements).Error
	if err != nil {
		return nil, 0, err
	}

	var totalStockMove int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.StockMovement{}).Count(&totalStockMove).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	if data, err := json.Marshal(stockMovements); err == nil {
		r.redis.Set(ctx, cachaKey, data, 5*time.Minute)
	}

	return stockMovements, totalStockMove, nil
}
