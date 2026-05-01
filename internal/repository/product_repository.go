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

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redis *redis.Client) domain.ProductRepository {
	return &productRepository{db: db, redis: redis}
}

func NewProductRepositoryWithTx(tx *gorm.DB) domain.ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return err
	}
	return utils.InvalidateCache(ctx, r.redis, "products:*")
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, pagination *domain.Pagination) ([]domain.Product, int64, error) {
	page := pagination.Page
	limit := pagination.Limit

	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("products:page:%d:limit:%d", page, limit)
	totalKey := "products:total"

	cached, err := r.redis.Get(ctx, totalKey).Result()
	if err == nil {
		var products []domain.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			totalStr, err := r.redis.Get(ctx, totalKey).Result()
			if err == nil {
				total, _ := strconv.ParseInt(totalStr, 10, 64)
				return products, total, nil
			}
		}
	}

	var products []domain.Product
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	var totalProduct int64
	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Count(&totalProduct).Error; err != nil {
		return nil, 0, err
	}

	if data, err := json.Marshal(products); err == nil {
		r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	r.redis.Set(ctx, totalKey, strconv.FormatInt(totalProduct, 10), 5*time.Minute)

	return products, totalProduct, nil
}

func (r *productRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
