package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type OrderService struct {
	repo        domain.OrderRepository
	productRepo domain.ProductRepository
	txm         domain.TxRepository
	stockRepo   domain.StockRepository
}

func NewOrderService(repo domain.OrderRepository, productRepo domain.ProductRepository, txm domain.TxRepository, stockRepo domain.StockRepository) *OrderService {
	return &OrderService{repo: repo, productRepo: productRepo, txm: txm, stockRepo: stockRepo}
}

func (s *OrderService) Create(ctx context.Context, req *dto.OrderRequest) error {
	productsIDs := make([]uuid.UUID, len(req.Items))
	for i, item := range req.Items {
		productsIDs[i] = item.ProductID
	}
	products, err := s.productRepo.GetByIDs(ctx, productsIDs)
	if err != nil {
		return err
	}

	productMap := make(map[uuid.UUID]*domain.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	var totalAmount uint
	items := make([]domain.OrderItem, len(req.Items))

	for i, item := range req.Items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return errors.New("Products not found")
		}

		totalAmount += product.ProductPrice * uint(item.Quantity)
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.ProductPrice,
		}
	}

	return s.txm.ExecOrderTx(ctx, func(
		orderRepo domain.OrderRepository,
		stockMovementRepo domain.StockMovementRepository,
		stockRepo domain.StockRepository) error {
		order := &domain.Order{
			UserID:      req.UserID,
			TotalAmount: totalAmount,
			Items:       items,
		}
		if err := orderRepo.Create(ctx, order); err != nil {
			return err
		}

		stockAdjustments := make([]domain.StockAdjustment, len(order.Items))
		stockMovements := make([]domain.StockMovement, len(order.Items))

		for i, item := range order.Items {
			stock, err := stockRepo.GetProductStock(ctx, item.ProductID)
			if err != nil {
				return err
			}

			stockAdjustments[i] = domain.StockAdjustment{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}

			stockMovements[i] = domain.StockMovement{
				StockID:      stock.ID,
				MovementType: "OUT",
				Quantity:     item.Quantity,
			}
		}

		if err := stockRepo.DecreaseStockBulkWithTx(ctx, stockAdjustments); err != nil {
			return err
		}

		if err := stockMovementRepo.CreateBulk(ctx, stockMovements); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetAll(ctx context.Context, req *domain.Pagination) ([]domain.Order, int64, error) {
	pagination := &domain.Pagination{
		Limit: req.Limit,
		Page:  req.Page,
	}
	orders, total, err := s.repo.GetAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
