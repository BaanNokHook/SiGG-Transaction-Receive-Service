package httprpc

import (
	"net/http"
	usecase "nextclan/transaction-gateway/transaction-receive-service/internal/usecase"
	"nextclan/transaction-gateway/transaction-receive-service/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, rrt usecase.ReceiveRawTransaction, rvt usecase.ReceiveValidatedTransaction, gau usecase.GetAddressUTXO) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	handler.Use(cors.New(corsConfig))

	// K8s probe
	//how well is the http server running
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	newTransactionRoutes(handler, rrt, rvt, gau, l)

}
