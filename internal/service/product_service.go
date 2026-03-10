package service

import (
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product *domain.Product) error {
	stock := &domain.Stock{
		Quantity: 0,
	}
	return s.repo.CreateWithStock(product, stock)
}

func (s *ProductService) GetByID(id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("Product not found")
	}
	return &dto.ProductResponse{
		ID:           product.ID,
		ProductName:  product.ProductName,
		ProductPrice: product.ProductPrice,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}, nil
}

func (s *ProductService) GetAll() ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("Products not found")
	}
	resp := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.ProductResponse{
			ID:           product.ID,
			ProductName:  product.ProductName,
			ProductPrice: product.ProductPrice,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		})
	}
	return resp, nil
}
