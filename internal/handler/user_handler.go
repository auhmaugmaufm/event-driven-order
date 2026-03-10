package handler

import (
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
	cfg     *config.Config
}

func NewUserHandler(svc *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{service: svc, cfg: cfg}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "bad_request",
			Message: "invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	if err := h.service.Create(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "internal_error",
			Message: "failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "registered successfully",
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{Error: "bad_request", Message: "invalid body"})
	}
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(dto.ErrorResponse{Error: "unauthorized", Message: err.Error()})
	}
	return c.JSON(fiber.Map{"token": token})
}

// func (h *UserHandler) Register(c *fiber.Ctx) error {
// 	var req dto.UserRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error:   "bad_request",
// 			Message: "invalid request body",
// 		})
// 	}

// 	if req.Email == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error:   "validation_error",
// 			Message: "Email is required",
// 		})
// 	}

// 	if req.Password == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error:   "validation_error",
// 			Message: "Password is required",
// 		})
// 	}

// 	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	user := &domain.User{
// 		Email:        req.Email,
// 		PasswordHash: string(bytes),
// 	}
// 	if err := h.service.Create(user); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
// 			Error:   "internal_error",
// 			Message: "failed to create user",
// 		})
// 	}
// 	return nil
// }
