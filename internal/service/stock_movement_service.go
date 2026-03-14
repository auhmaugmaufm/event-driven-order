package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/google/uuid"
)

type StockMovementService struct {
	repo      domain.StockMovementRepository
	txm       domain.TxRepository
	stockRepo domain.StockRepository
}

func NewStockMovementService(
	repo domain.StockMovementRepository,
	txm domain.TxRepository,
	stockRepo domain.StockRepository) *StockMovementService {
	return &StockMovementService{repo: repo, txm: txm, stockRepo: stockRepo}
}

func (s *StockMovementService) Create(ctx context.Context, req *dto.StockMovementRequest) error {
	stock, err := s.stockRepo.GetProductStock(ctx, req.ProductID)
	if err != nil {
		return err
	}
	return s.txm.ExecStockMovementTx(ctx, func(repo domain.StockMovementRepository, stockRepo domain.StockRepository) error {
		stockMovement := &domain.StockMovement{
			StockID:      stock.ID,
			MovementType: req.MovementType,
			Quantity:     req.Quantity,
		}
		if err := s.repo.Create(ctx, stockMovement); err != nil {
			return err
		}
		// stock := &domain.Stock{
		// 	ProductID: req.ProductID,
		// 	Quantity:  req.Quantity,
		// } --> useless
		if req.MovementType == "IN" {
			return s.stockRepo.IncreaseStockWithTx(ctx, stock.ID, req.Quantity)
		} else {
			return s.stockRepo.DecreaseStockWithTx(ctx, stock.ID, req.Quantity)
		}
	})
}

func (s *StockMovementService) GetByMovementID(ctx context.Context, id uuid.UUID) (*dto.StockMovementResponse, error) {
	movement, err := s.GetByMovementID(ctx, id)
	if err != nil {
		return nil, errors.New("Stock not found")
	}
	return &dto.StockMovementResponse{
		ID:           movement.ID,
		ProductID:    movement.ID,
		MovementType: movement.MovementType,
		Quantity:     movement.Quantity,
		CreatedAt:    movement.CreatedAt,
		UpdatedAt:    movement.UpdatedAt,
	}, nil
}

func (s *StockMovementService) GetAllMovement(ctx context.Context) ([]dto.StockMovementResponse, error) {
	movements, err := s.repo.GetStockMovement(ctx)
	if err != nil {
		return nil, errors.New("Movements not found")
	}
	res := make([]dto.StockMovementResponse, 0, len(movements))
	for _, movement := range movements {
		res = append(res, dto.StockMovementResponse{
			ID:           movement.ID,
			ProductID:    movement.ID,
			MovementType: movement.MovementType,
			Quantity:     movement.Quantity,
			CreatedAt:    movement.CreatedAt,
			UpdatedAt:    movement.UpdatedAt,
		})
	}
	return res, nil
}
