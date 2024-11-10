// cmd/main.go
package main

import (
	"context"
	"fmt"
	"main/pkg/ai"
	"main/pkg/config"
	"main/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()
	}()

	// Initialize AI client
	client, err := ai.NewClient(ctx, &ai.Config{
		APIKey:     cfg.APIKey,
		ModelName:  cfg.ModelName,
		ImageModel: cfg.ImageModel,
	}, log)
	if err != nil {
		log.Fatal("Failed to initialize AI client", zap.Error(err))
	}
	defer client.Close()

	// Process text
	if err := client.ProcessText(ctx, "Tell me about AI", os.Stdout); err != nil {
		log.Error("Failed to process text", zap.Error(err))
	}

	// Process image
	if err := client.ProcessImage(ctx, "./test.png", os.Stdout); err != nil {
		log.Error("Failed to process image", zap.Error(err))
	}
}
