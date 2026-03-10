package main

import (
	"fmt"
	"log"

	"github.com/auhmaugmaufm/event-driven-order/internal/auth"
	"github.com/auhmaugmaufm/event-driven-order/internal/handler"
	"github.com/auhmaugmaufm/event-driven-order/internal/middleware"
	"github.com/auhmaugmaufm/event-driven-order/internal/repository"
	"github.com/auhmaugmaufm/event-driven-order/internal/service"
	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/auhmaugmaufm/event-driven-order/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config.Load()
	cfg := config.Get()

	database.RunMigrations(cfg)
	db := database.NewPostgresDB(cfg)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpireHour)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtManager)
	userHandler := handler.NewUserHandler(userService, cfg)

	stockReposity := repository.NewStockRepository(db)
	stockService := service.NewStockService(stockReposity)
	stockHandler := handler.NewStockHandler(stockService, cfg)

	productRepository := repository.NewProductReposity(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProeductHandler(productService, cfg)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "internal_error",
				"message": err.Error(),
			})
		},
	})

	api := app.Group("/api/v1")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	protected := api.Group("", middleware.AuthMiddleware(jwtManager))

	product := protected.Group("/product")
	product.Get("", productHandler.GetAllProducts)
	product.Post("/create", productHandler.Create)
	product.Get("/:id", productHandler.GetProductByID)

	stock := protected.Group("/stock")
	stock.Get("", stockHandler.GetAllProductStocks)
	stock.Get("/:id", stockHandler.GetProductStock)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(app.Listen(addr))
}
