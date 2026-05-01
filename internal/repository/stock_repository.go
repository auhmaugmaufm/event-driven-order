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

type stockRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewStockRepository(db *gorm.DB, redis *redis.Client) domain.StockRepository {
	return &stockRepository{db: db, redis: redis}
}

func NewStockRepositoryWithTx(tx *gorm.DB) domain.StockRepository {
	return &stockRepository{db: tx}
}

func (r *stockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	err := r.db.WithContext(ctx).Create(stock).Error
	if err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "stocks:*")
}

func (r *stockRepository) IncreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error {
	err := r.db.WithContext(ctx).Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
	if err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "stocks:*")
}

func (r *stockRepository) DecreaseStockWithTx(ctx context.Context, productId uuid.UUID, quantity int) error {
	err := r.db.WithContext(ctx).Model(&domain.Stock{}).Where("product_id = ?", productId).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
	if err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "stocks:*")
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
		return utils.InvalidateCache(ctx, r.redis, "stocks:*")
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

	cacheKey := fmt.Sprintf("products:page:%d:limit:%d", page, limit)
	totalKey := "products:total"

	cached, err := r.redis.Get(ctx, totalKey).Result()
	if err == nil {
		var stocks []domain.Stock
		if err := json.Unmarshal([]byte(cached), &stocks); err == nil {
			totalStr, err := r.redis.Get(ctx, totalKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(totalStr, 10, 64)
				return stocks, total, nil
			}
		}
	}

	var stocks []domain.Stock
	err = r.db.WithContext(ctx).Limit(pagination.Limit).Offset(offset).Preload("Product").Find(&stocks).Error
	if err != nil {
		return nil, 0, err
	}

	var totalStock int64
	fetchTotalError := r.db.WithContext(ctx).Model(&domain.Stock{}).Count(&totalStock).Error
	if fetchTotalError != nil {
		return nil, 0, fetchTotalError
	}

	if data, err := json.Marshal(stocks); err == nil {
		r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	r.redis.Set(ctx, totalKey, strconv.FormatInt(totalStock, 10), 5*time.Minute)

	return stocks, totalStock, nil
}
