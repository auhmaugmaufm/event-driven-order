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

type orderRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client) domain.OrderRepository {
	return &orderRepository{db: db, redis: redis}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "orders:*")
}

func (r *orderRepository) GetAll(ctx context.Context, pagination *domain.Pagination) ([]domain.Order, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("orders:page:%d:limit:%d", page, limit)
	totalKey := "orders:total"

	cached, err := r.redis.Get(ctx, totalKey).Result()
	if err == nil {
		var orders []domain.Order
		if err := json.Unmarshal([]byte(cached), &orders); err == nil {
			totalStr, err := r.redis.Get(ctx, totalKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(totalStr, 10, 64)
				return orders, total, nil
			}
		}
	}

	var orders []domain.Order
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Preload("Items").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	var totalOrder int64
	if err := r.db.WithContext(ctx).Model(&domain.Order{}).Count(&totalOrder).Error; err != nil {
		return nil, 0, err
	}

	if data, err := json.Marshal(orders); err == nil {
		r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	r.redis.Set(ctx, totalKey, strconv.FormatInt(totalOrder, 10), 5*time.Minute)

	return orders, totalOrder, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order *domain.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
