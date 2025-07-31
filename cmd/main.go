package main

import (
	"context"
	"goCodetest/internal/routers"
	"goCodetest/pkg/logger"
	"goCodetest/pkg/product_store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var pStore *product_store.ProductStore

func init() {
	pStore = product_store.NewProductStore()
}

func main() {
	logger.Info("start app...")

	app := gin.New()
	routers.Setup(app, pStore)

	server := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("run server err %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	logger.Info("Signal %s received, shutting down server...", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %s", err.Error())
	} else {
		logger.Info("Server stopped")
	}
}
