package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/internal/service"
	"simcomm-monolith/util"
	"time"

	_ "simcomm-monolith/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
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

	e.Logger.SetLevel(log.DEBUG)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Debug("debug nih")
	e.Logger.Info("info nih")
	e.Logger.Error("error nih")

	log.Debug("gomon debug nih")
	log.Info("gomon info nih")
	log.Error("gomon error nih")

	// Routes
	e.GET("/health", health)

	cfg := config.GetConfig()
	if cfg.ServerConfig.Env == "dev" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	db := util.GetDB(cfg)
	redisClient := util.GetRedisClient(cfg)
	rabbitMQConnection := util.GetRabbitMQConnection(cfg.RabbitMQConfig)
	tpQueue := repository.NewQueueDeclare(rabbitMQConnection, "transfer_product")

	var queues []repository.Queue
	queues = append(queues, tpQueue)

	userRepo := repository.NewPostgreUserRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient, cfg)
	svc := service.NewUserService(userRepo, redisRepo, cfg)
	RegisterUserHandler(e, svc)

	productRepo := repository.NewPostgreProductRepository(db)
	productSvc := service.NewProductService(productRepo, redisRepo, cfg)
	RegisterProductHandler(e, productSvc)

	warehouseRepo := repository.NewPostgreWarehouseRepository(db)
	warehouseSvc := service.NewWarehouseService(warehouseRepo, redisRepo, cfg)
	tpQueue.AddReceiver(context.Background(), warehouseSvc.ProcessTPQueue)
	RegisterWarehouseHandler(e, warehouseSvc)

	shopRepo := repository.NewPostgreShopRepository(db)
	shopSvc := service.NewShopService(warehouseSvc, shopRepo, redisRepo, tpQueue, cfg)
	RegisterShopHandler(e, shopSvc)

	orderRepo := repository.NewPostgreOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, redisRepo, cfg)
	RegisterOrderHandler(e, orderSvc)

	// Start server
	// e.Logger.Fatal(e.Start(fmt.Sprintf("%v", cfg.ServerConfig.Host) + ":" + fmt.Sprintf("%v", cfg.ServerConfig.Port)))

	queueHandler := QueueHandler{
		queues: queues,
		conn:   rabbitMQConnection,
	}

	HandleServer(e, cfg, queueHandler)
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "Server Up")
}

func HandleServer(e *echo.Echo, cfg *config.Config, qh QueueHandler) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		err := e.Start(fmt.Sprintf("%v", cfg.ServerConfig.Host) + ":" + fmt.Sprintf("%v", cfg.ServerConfig.Port))
		log.Error(err)
		e.Logger.Fatal(err)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	qh.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type QueueHandler struct {
	queues []repository.Queue
	conn   *amqp.Connection
}

func (qh *QueueHandler) Close() {
	for i := 0; i < len(qh.queues); i++ {
		qh.queues[i].Close()
	}
	qh.conn.Close()
}
