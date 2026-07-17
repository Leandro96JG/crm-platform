package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatAPIMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type ChatAPIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewClient(apiKey string, model string) *Client {
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &Client{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) ChatCompletion(ctx context.Context, systemPrompt string, messages []domain.ChatMessage) (string, error) {
	apiMessages := []ChatAPIMessage{
		{Role: "system", Content: systemPrompt},
	}

	for _, msg := range messages {
		apiMessages = append(apiMessages, ChatAPIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	reqBody := ChatCompletionRequest{
		Model:       c.model,
		Messages:    apiMessages,
		Temperature: 0.7,
		MaxTokens:   1024,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}
