package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"xenigo/internal/output"
)

type DiscordWebhook struct {
    Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
    Title       string       `json:"title"`
    Description string       `json:"description"`
    URL         string       `json:"url,omitempty"`
    Author      EmbedAuthor  `json:"author,omitempty"`
    Fields      []EmbedField `json:"fields,omitempty"`
}

type EmbedAuthor struct {
    Name string `json:"name"`
}

type EmbedField struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type DiscordSender struct {
    WebhookURL string
}

func (d *DiscordSender) SendMessage(embed output.MessageEmbed) error {
    discordEmbed := DiscordEmbed{
        Title:       embed.Title,
        Description: embed.Description,
        URL:         embed.URL,
        Author:      EmbedAuthor{Name: embed.Author},
        Fields:      convertFields(embed.Fields),
    }

    webhook := DiscordWebhook{Embeds: []DiscordEmbed{discordEmbed}}
    webhookBody, err := json.Marshal(webhook)
    if err != nil {
        return fmt.Errorf("failed to marshal webhook body: %w", err)
    }

    resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(webhookBody))
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("received non-204 response code: %d", resp.StatusCode)
    }

    return nil
}

func convertFields(fields []output.EmbedField) []EmbedField {
    var embedFields []EmbedField
    for _, field := range fields {
        embedFields = append(embedFields, EmbedField{
            Name:  field.Name,
            Value: field.Value,
        })
    }
    return embedFields
}