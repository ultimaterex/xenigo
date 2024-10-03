package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"xenigo/internal/config"
	"xenigo/internal/reddit"
)

type DiscordEmbed struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Author      struct {
        Name string `json:"name"`
    } `json:"author"`
    Fields []struct {
        Name  string `json:"name"`
        Value string `json:"value"`
    } `json:"fields"`
}

type DiscordWebhook struct {
    Embeds []DiscordEmbed `json:"embeds"`
}

func SendToDiscord(webhookURL string, embed DiscordEmbed) error {
    webhook := DiscordWebhook{Embeds: []DiscordEmbed{embed}}
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

func CreateDiscordEmbed(post reddit.RedditPost, target config.Target) DiscordEmbed {
    embed := DiscordEmbed{
        Title:       post.Title,
        Description: post.Selftext,
    }
    if *target.Output.Format.URL {
        embed.URL = post.URL
    }
    if *target.Output.Format.Author {
        embed.Author.Name = post.Author
    }
    if *target.Output.Format.Subreddit {
        embed.Fields = append(embed.Fields, struct {
            Name  string `json:"name"`
            Value string `json:"value"`
        }{
            Name:  "Subreddit",
            Value: target.Monitor.Subreddit,
        })
    }
    if *target.Output.Format.DiscussionURL && post.Permalink != "" {
        embed.Fields = append(embed.Fields, struct {
            Name  string `json:"name"`
            Value string `json:"value"`
        }{
            Name:  "Discussion URL",
            Value: fmt.Sprintf("https://www.reddit.com%s", post.Permalink),
        })
    }
    return embed
}