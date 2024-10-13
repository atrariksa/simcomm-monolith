package handler

import (
	"fmt"
	"net/http"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/internal/service"
	"simcomm-monolith/util"

	_ "simcomm-monolith/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger
// @version 1.0
// @description Swagger for User Service
// @host localhost:6022
// @BasePath /
func SetupServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", health)

	cfg := config.GetConfig()
	if cfg.ServerConfig.Env == "dev" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	userRepo := repository.NewPostgreUserRepository(util.GetDB(cfg))
	redisRepo := repository.NewRedisRepository(util.GetRedisClient(cfg), cfg)
	svc := service.NewUserService(userRepo, redisRepo, cfg)
	RegisterUserHandler(e, svc)

	productRepo := repository.NewPostgreProductRepository(util.GetDB(cfg))
	productSvc := service.NewProductService(productRepo, redisRepo, cfg)
	RegisterProductHandler(e, productSvc)

	warehouseRepo := repository.NewPostgreWarehouseRepository(util.GetDB(cfg))
	warehouseSvc := service.NewWarehouseService(warehouseRepo, redisRepo, cfg)
	RegisterWarehouseHandler(e, warehouseSvc)

	shopRepo := repository.NewPostgreShopRepository(util.GetDB(cfg))
	shopSvc := service.NewShopService(shopRepo, redisRepo, cfg)
	RegisterShopHandler(e, shopSvc)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%v", cfg.ServerConfig.Host) + ":" + fmt.Sprintf("%v", cfg.ServerConfig.Port)))
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "Server Up")
}
