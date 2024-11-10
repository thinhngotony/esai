package ai

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
)

type TextProcessor interface {
	ProcessText(ctx context.Context, input string, writer io.Writer) error
	ProcessImage(ctx context.Context, imagePath string, writer io.Writer) error
}

func (c *Client) ProcessText(ctx context.Context, input string, writer io.Writer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.logger.Info("processing text input",
		zap.String("input_length", fmt.Sprint(len(input))),
		zap.String("model", c.config.ModelName),
	)

	model := c.client.GenerativeModel(c.config.ModelName)

	// Configure the model
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(2048)

	// Set safety settings if needed
	/*
	   model.SafetySettings = []*genai.SafetySetting{
	       {
	           Category:  genai.HarmCategoryHarassment,
	           Threshold: genai.HarmBlockMedium,
	       },
	       // Add other safety settings as needed
	   }
	*/

	var lastError error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Log retry attempt
			c.logger.Info("retrying text processing",
				zap.Int("attempt", attempt+1),
				zap.String("previous_error", lastError.Error()),
			)

			// Wait before retrying with exponential backoff
			backoff := initialBackoff * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(backoff):
			}
		}

		// Enable streaming with proper error handling
		iter := model.GenerateContentStream(ctx, genai.Text(input))

		// Process the stream
		var hadContent bool
		for {
			resp, err := iter.Next()
			if err == io.EOF {
				if hadContent {
					return nil // Successfully processed all content
				}
				break // No content received, will retry if attempts remain
			}
			if err != nil {
				lastError = fmt.Errorf("generation error: %w", err)
				break // Break the inner loop to retry
			}

			// Process candidates
			if len(resp.Candidates) == 0 {
				continue
			}

			candidate := resp.Candidates[0]
			// Check for blocked content
			if candidate.FinishReason == genai.FinishReasonSafety {
				return fmt.Errorf("content blocked by safety settings")
			}

			// Process content parts
			if len(candidate.Content.Parts) == 0 {
				continue
			}

			part := candidate.Content.Parts[0]
			if text, ok := part.(genai.Text); ok {
				if _, err := writer.Write([]byte(text)); err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
				hadContent = true
			}
		}

		if hadContent {
			return nil // Successfully processed some content
		}
	}

	// If we've exhausted all retries
	if lastError != nil {
		return fmt.Errorf("failed to process text after %d attempts: %w", maxRetries, lastError)
	}

	return fmt.Errorf("no content generated after %d attempts", maxRetries)
}
