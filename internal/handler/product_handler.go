package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service *service.ProductService
	cfg     *config.Config
}

func NewProeductHandler(svc *service.ProductService, cfg *config.Config) *ProductHandler {
	return &ProductHandler{service: svc, cfg: cfg}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid request body",
		})
	}
	if req.ProductName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Product name is required",
		})
	}
	if req.ProductPrice <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Product price is required",
		})
	}

	product := &domain.Product{
		ProductName:  req.ProductName,
		ProductPrice: req.ProductPrice,
	}
	if err := h.service.Create(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to create product",
		})
	}
	return nil
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid product id",
		})
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   resp,
		Status: fiber.StatusOK,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	resp, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Data:   resp,
		Status: fiber.StatusOK,
	})
}
