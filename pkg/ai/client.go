package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Client struct {
	client *genai.Client
	logger *zap.Logger
	config *Config
	mu     sync.RWMutex
}

type Config struct {
	APIKey     string
	ModelName  string
	ImageModel string
}

func NewClient(ctx context.Context, cfg *Config, logger *zap.Logger) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		client: client,
		logger: logger,
		config: cfg,
	}, nil
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client != nil {
		c.client.Close()
	}
}
