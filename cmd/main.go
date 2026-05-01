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

	rdb := database.NewRedisCache(cfg)
	defer rdb.Close()

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpireHour)

	txManager := repository.NewTxManager(db)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtManager)
	userHandler := handler.NewUserHandler(userService, cfg)

	stockRepository := repository.NewStockRepository(db, rdb)
	stockService := service.NewStockService(stockRepository)
	stockHandler := handler.NewStockHandler(stockService, cfg)

	productRepository := repository.NewProductRepository(db, rdb)
	productService := service.NewProductService(productRepository, txManager, stockRepository)
	productHandler := handler.NewProductHandler(productService, cfg)

	stockMovementRepository := repository.NewStockMovementRepository(db, rdb)
	stockMovementService := service.NewStockMovementService(stockMovementRepository, txManager, stockRepository)
	stockMovementHandler := handler.NewStockMovementHandler(stockMovementService, cfg)

	orderRepository := repository.NewOrderRepository(db, rdb)
	orderService := service.NewOrderService(orderRepository, productRepository, txManager, stockRepository)
	orderHandler := handler.NewOrderHandler(orderService, cfg)

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

	stockMovement := protected.Group("/stock-movement")
	stockMovement.Post("", stockMovementHandler.Create)
	stockMovement.Get("", stockMovementHandler.GetAllMovement)
	stockMovement.Get("/:id", stockMovementHandler.GetMovementByID)
	stockMovement.Get("", stockMovementHandler.GetAllMovementType)

	order := protected.Group("/order")
	order.Get("", orderHandler.GetAll)
	order.Post("/create", orderHandler.Create)
	order.Get("/:id", orderHandler.GetByID)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(app.Listen(addr))
}
