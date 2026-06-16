package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultAPIBaseURL = "https://api.telegram.org"

type Client struct {
	apiBaseURL string
	token      string
	httpClient *http.Client
}

func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		apiBaseURL: defaultAPIBaseURL,
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	body, err := json.Marshal(sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal telegram sendMessage request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.methodURL("sendMessage"), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create telegram sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call telegram sendMessage: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read telegram sendMessage response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram sendMessage returned %s: %s", resp.Status, truncate(respBody, 500))
	}

	var parsed telegramResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return fmt.Errorf("decode telegram sendMessage response: %w", err)
	}
	if !parsed.OK {
		return fmt.Errorf("telegram sendMessage failed: %s", parsed.Description)
	}

	return nil
}

func (c *Client) methodURL(method string) string {
	return c.apiBaseURL + "/bot" + c.token + "/" + method
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func truncate(body []byte, limit int) string {
	text := string(body)
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}
