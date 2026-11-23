package main

import (
	"flag"
	"fmt"
	"os"

	"booking_system/booking_service/pkg/server"

	"go.uber.org/zap"
)

func main() {
	// Define command-line flags
	env := flag.String("env", "dev", "Environment (dev, qc, uat, prod)")
	flag.Parse()

	// Initialize application
	app, err := server.InitializeApp(*env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	// Start server
	if err := server.RunServer(app); err != nil {
		app.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
