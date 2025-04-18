package main

import (
	"test-wallet/bootstrap"
	"test-wallet/utils"
)

func main() {
	// Initialize application
	if err := bootstrap.InitializeApp(); err != nil {
		utils.LogFatal(err, "Failed to initialize application", nil)
	}

	// Setup router and server
	router := bootstrap.SetupRouter()
	srv := bootstrap.SetupServer(router)

	// Start server and wait for shutdown
	bootstrap.StartServer(srv)
	bootstrap.WaitForShutdown(srv)
}
