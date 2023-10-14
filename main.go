package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const socketPath = "/tmp/test.sock"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := NewRouter()

	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		slog.Error("Unable to listen", slog.String("error", err.Error()))
		return
	}

	httpServer := &http.Server{
		Handler: router,
	}

	go func() {
		if err := httpServer.Serve(socket); err != nil && err != http.ErrServerClosed {
			slog.Error("Unable to start server", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: ", slog.String("error", err.Error()))
	}

	slog.Info("Server exiting")
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", welcome)
	r.GET("/health", health)

	return r
}

func welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the API",
	})
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
