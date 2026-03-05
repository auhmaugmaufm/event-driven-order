package repository

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductReposity(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	err := r.db.Create(product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
