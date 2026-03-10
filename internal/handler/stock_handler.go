package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StockHandler struct {
	service *service.StockService
	cfg     *config.Config
}

func NewStockHandler(svc *service.StockService, cfg *config.Config) *StockHandler {
	return &StockHandler{service: svc, cfg: cfg}
}

func (h *StockHandler) IncreaseStock(c *fiber.Ctx) error {
	var req dto.UpdateStockReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})
	}

	err := h.service.IncreaseStock(req.ProductID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *StockHandler) DecreaseStock(c *fiber.Ctx) error {
	var req dto.UpdateStockReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})
	}

	err := h.service.DecreaseStock(req.ProductID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *StockHandler) GetProductStock(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid product id",
		})
	}

	res, err := h.service.GetProductStock(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   res,
		Status: fiber.StatusOK,
	})

}

func (h *StockHandler) GetAllProductStocks(c *fiber.Ctx) error {
	res, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   res,
		Status: fiber.StatusOK,
	})

}
