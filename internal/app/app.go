package app

import (
	"fmt"
	"nextclan/transaction-gateway/transaction-receive-service/config"
	v1 "nextclan/transaction-gateway/transaction-receive-service/internal/controller/http/v1"
	"nextclan/transaction-gateway/transaction-receive-service/internal/controller/httprpc"
	usecase "nextclan/transaction-gateway/transaction-receive-service/internal/usecase"
	rpc "nextclan/transaction-gateway/transaction-receive-service/pkg/httprpc"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/httpserver"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/loaffinity"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"
	messaging "nextclan/transaction-gateway/transaction-receive-service/pkg/rabbitmq"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	fmt.Println("Starting App...")

	loaffinity := loaffinity.NewLoaffinityClient(cfg.Loaffinity.URL, cfg.Loaffinity.Username, cfg.Loaffinity.Password, l)
	//Init Redis
	redisCache := redis.NewRedisClient(redis.Addr(cfg.Redis.Addr), redis.Password(cfg.Redis.Password))

	// Use case
	receiveRawTransactionUseCase := usecase.NewReceiveRawTransaction(l, redisCache)
	receiveValidatedTransactionUseCase := usecase.NewReceiveValidatedTransaction(l)
	getAddressUTXOlmUseCase := usecase.NewGetAddressUTXO(l, loaffinity)

	// HTTP Server
	httpServer := initializeRPC(l, receiveRawTransactionUseCase, receiveValidatedTransactionUseCase, getAddressUTXOlmUseCase, cfg)

	//Init client
	initializeMessaging(cfg)

	// Shutdown
	ShutdownApplicationHandler(l, httpServer)
}

func initializeHttp(l *logger.Logger, rrt *usecase.ReceiveRawTransactionUsecase, rvt *usecase.ReceiveValidatedTransactionUseCase, cfg *config.Config) *httpserver.Server {
	handler := gin.New()
	v1.NewRouter(handler, rrt, rvt, l)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	return httpServer
}

func initializeRPC(l *logger.Logger, rrt *usecase.ReceiveRawTransactionUsecase, rvt *usecase.ReceiveValidatedTransactionUseCase, gau *usecase.GetAddressUTXOUseCase, cfg *config.Config) *rpc.Server {
	handler := gin.New()
	httprpc.NewRouter(handler, l, rrt, rvt, gau)
	httpServer := rpc.New(handler, rpc.Port(cfg.HTTP.Port))
	return httpServer
}

func initializeMessaging(cfg *config.Config) {
	//TODO dependency injection for usecase scope
	usecase.MessagingClient = &messaging.MessagingClient{}
	usecase.MessagingClient.Connect(cfg.RMQ.URL)
}

func ShutdownApplicationHandler(l *logger.Logger, httpServer *rpc.Server) {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	}

	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = usecase.MessagingClient.Close()
}
