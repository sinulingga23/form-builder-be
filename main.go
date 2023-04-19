package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinulingga23/form-builder-be/config"
	delivery "github.com/sinulingga23/form-builder-be/delivery/http"
	"github.com/sinulingga23/form-builder-be/implement/repository"
	"github.com/sinulingga23/form-builder-be/implement/usecase"
	"github.com/sinulingga23/form-builder-be/monitoring"
)

var (
	port = "8085"
)

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
	}))

	// repository
	db, errConnectDB := config.ConnectDB()
	if errConnectDB != nil {
		log.Fatalf("errConnectDB: %v", errConnectDB)
	}
	mPartnerRepository := repository.NewMPartnerRepository(db)
	mFieldTypeRepository := repository.NewMFieldTypeRepository(db)
	mFormRepository := repository.NewMFormRepository(db)
	mFormFieldRepository := repository.NewMFormFieldRepository(db)
	mFormFieldChildsRepository := repository.NewMFormFieldChildsRepository(db)

	// metric
	registry := prometheus.NewRegistry()
	metric := monitoring.NewMetric(registry)

	// usecase
	mFormUsecase := usecase.NewMFormUsecase(
		db,
		mPartnerRepository,
		mFieldTypeRepository,
		mFormRepository,
		mFormFieldRepository,
		mFormFieldChildsRepository,
		metric,
	)

	// delivery http
	formHttp := delivery.NewFormHttp(mFormUsecase)
	formHttp.ServeHandler(&r.RouterGroup)

	promhttp.Handler()
	r.GET("/metrics", func(ctx *gin.Context) {
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})

	log.Printf("form-builder-be service served on :%v", port)
	if errListenAndServe := http.ListenAndServe(fmt.Sprintf(":%v", port), r); errListenAndServe != nil {
		log.Fatalf("errListenAndServe: %v", errListenAndServe)
	}
}
