package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StockMovementHandler struct {
	service *service.StockMovementService
	cfg     *config.Config
}

func NewStockMovementHandler(svc *service.StockMovementService, cfg *config.Config) *StockMovementHandler {
	return &StockMovementHandler{service: svc, cfg: cfg}
}

func (h *StockMovementHandler) Create(c *fiber.Ctx) error {
	var req dto.StockMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}
	if err := h.service.Create(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to create product",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "createad successfully",
	})
}

func (h *StockMovementHandler) GetMovementByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid product id",
		})
	}

	res, err := h.service.GetByMovementID(c.Context(), id)
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

func (h *StockMovementHandler) GetAllMovement(c *fiber.Ctx) error {
	res, err := h.service.GetAllMovement(c.Context())
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
