package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-wallet/config"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
)

// SetupServer creates and configures the HTTP server
func SetupServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + config.AppConfig.ServerConfig.Port,
		Handler:      router,
		ReadTimeout:  config.AppConfig.ServerConfig.ReadTimeout,
		WriteTimeout: config.AppConfig.ServerConfig.WriteTimeout,
		IdleTimeout:  config.AppConfig.ServerConfig.IdleTimeout,
	}
}

// StartServer starts the HTTP server in a goroutine
func StartServer(srv *http.Server) {
	go func() {
		utils.LogInfo("Server starting on port "+config.AppConfig.ServerConfig.Port, nil)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogFatal(err, "Failed to start server", nil)
		}
	}()
}

// WaitForShutdown waits for interrupt signal and performs graceful shutdown
func WaitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.LogInfo("Shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError(err, "Server forced to shutdown", nil)
	}

	utils.LogInfo("Server exiting", nil)
}
