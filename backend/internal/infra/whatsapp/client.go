package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	phoneNumberID string
	accessToken   string
	apiVersion    string
	httpClient    *http.Client
}

type TextMessage struct {
	Body string `json:"body"`
}

type WhatsAppTextRequest struct {
	MessagingProduct string     `json:"messaging_product"`
	RecipientType    string     `json:"recipient_type"`
	To               string     `json:"to"`
	Type             string     `json:"type"`
	Text             TextMessage `json:"text"`
}

type WhatsAppTemplateRequest struct {
	MessagingProduct string              `json:"messaging_product"`
	RecipientType    string              `json:"recipient_type"`
	To               string              `json:"to"`
	Type             string              `json:"type"`
	Template         TemplateComponents  `json:"template"`
}

type TemplateComponents struct {
	Name       string              `json:"name"`
	Language   Language            `json:"language"`
	Components []TemplateComponent `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters"`
}

type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type WhatsAppResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Error *WhatsAppError `json:"error,omitempty"`
}

type WhatsAppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewClient(phoneNumberID string, accessToken string, apiVersion string) *Client {
	if apiVersion == "" {
		apiVersion = "v21.0"
	}

	return &Client{
		phoneNumberID: phoneNumberID,
		accessToken:   accessToken,
		apiVersion:    apiVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SendText(ctx context.Context, to string, message string) error {
	reqBody := WhatsAppTextRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: TextMessage{
			Body: message,
		},
	}

	return c.sendMessage(ctx, reqBody)
}

func (c *Client) SendTemplate(ctx context.Context, to string, templateName string, parameters map[string]string) error {
	components := []TemplateComponent{}
	if len(parameters) > 0 {
		params := make([]TemplateParameter, 0)
		for _, value := range parameters {
			params = append(params, TemplateParameter{
				Type: "text",
				Text: value,
			})
		}
		components = append(components, TemplateComponent{
			Type:       "body",
			Parameters: params,
		})
	}

	reqBody := WhatsAppTemplateRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template: TemplateComponents{
			Name:     templateName,
			Language: Language{Code: "es"},
			Components: components,
		},
	}

	return c.sendMessage(ctx, reqBody)
}

func (c *Client) VerifyWebhook(ctx context.Context, token string, challenge string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("verification token is required")
	}

	return challenge, nil
}

func (c *Client) sendMessage(ctx context.Context, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", c.apiVersion, c.phoneNumberID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	var result WhatsAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return fmt.Errorf("whatsapp API error: %s (code %d)", result.Error.Message, result.Error.Code)
	}

	return nil
}
