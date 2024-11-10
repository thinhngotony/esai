package ai

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
)

func (c *Client) ProcessImage(ctx context.Context, imagePath string, writer io.Writer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.logger.Info("processing image",
		zap.String("image_path", imagePath),
	)

	data, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image file: %w", err)
	}

	model := c.client.GenerativeModel(c.config.ImageModel)
	prompt := genai.ImageData("png", data)

	// Enable streaming for image processing
	iter := model.GenerateContentStream(ctx, prompt)

	for {
		resp, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			continue
		}

		part := resp.Candidates[0].Content.Parts[0]
		if text, ok := part.(genai.Text); ok {
			if _, err := writer.Write([]byte(text)); err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
		}
	}

	return nil
}
