package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendMessageCallsTelegramAPI(t *testing.T) {
	var gotRequest sendMessageRequest

	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", req.Method, http.MethodPost)
			}
			if req.URL.String() != "https://telegram.test/bottoken/sendMessage" {
				t.Fatalf("url = %q, want %q", req.URL.String(), "https://telegram.test/bottoken/sendMessage")
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("Content-Type = %q, want %q", req.Header.Get("Content-Type"), "application/json")
			}
			if err := json.NewDecoder(req.Body).Decode(&gotRequest); err != nil {
				t.Fatalf("decode request body: %v", err)
			}

			return jsonResponse(http.StatusOK, `{"ok":true}`), nil
		}),
	}

	client := NewClient("token", httpClient)
	client.apiBaseURL = "https://telegram.test"

	err := client.SendMessage(context.Background(), 123, "hello")
	if err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}
	if gotRequest.ChatID != 123 {
		t.Fatalf("chatID = %d, want 123", gotRequest.ChatID)
	}
	if gotRequest.Text != "hello" {
		t.Fatalf("text = %q, want %q", gotRequest.Text, "hello")
	}
}

func TestSendMessageReturnsTelegramAPIError(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"ok":false,"description":"chat not found"}`), nil
		}),
	}

	client := NewClient("token", httpClient)
	client.apiBaseURL = "https://telegram.test"

	err := client.SendMessage(context.Background(), 123, "hello")
	if err == nil {
		t.Fatal("SendMessage returned nil error, want error")
	}
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
}
