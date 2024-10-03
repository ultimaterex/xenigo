package discord

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type DiscordWebhook struct {
    Content string `json:"content"`
}

func NotifyDiscord(webhookURL string, content string) error {
    webhook := DiscordWebhook{Content: content}
    webhookBody, err := json.Marshal(webhook)
    if err != nil {
        return fmt.Errorf("failed to marshal webhook body: %w", err)
    }

    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(webhookBody))
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("received non-204 response code: %d", resp.StatusCode)
    }

    return nil
}