package main

import (
    "log"
    "time"
    "xenigo/internal/config"
    "xenigo/internal/discord"
    "xenigo/internal/reddit"
)

func monitorSubreddit(target config.Target, accessToken, userAgent, context string, cache *Cache, sendInitial bool) {
    // Function to fetch and process Reddit data
    fetchAndProcess := func(sendToDiscord bool) {
        log.Printf("Executing monitor check for subreddit: %s", target.Monitor.Subreddit)
        redditResponse, err := reddit.FetchRedditData(target, accessToken, userAgent, context)
        if err != nil {
            log.Printf("Error fetching Reddit data for subreddit %s: %v", target.Monitor.Subreddit, err)
            return
        }

        for _, child := range redditResponse.Data.Children {
            post := child.Data
            // Check if the post has already been processed
            if !cache.IsProcessed(post.Permalink) {
                if sendToDiscord {
                    processAndSendPost(post, target)
                }
                // Mark the post as processed
                cache.AddProcessedPermalink(post.Permalink)
            }
        }
    }

    // Determine if we should send the initial fetch data to Discord
    sendToDiscord := sendInitial

    // Run immediately on start
    fetchAndProcess(sendToDiscord)

    // Set up the ticker for subsequent runs
    ticker := time.NewTicker(time.Duration(target.Options.Interval) * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        fetchAndProcess(true)
    }
}

func processAndSendPost(post reddit.RedditPost, target config.Target) {
    embed := discord.CreateDiscordEmbed(post, target)

    log.Printf("Subreddit: %s\nTitle: %s\nURL: %s\n", target.Monitor.Subreddit, post.Title, post.URL)
    if post.Author != "" {
        log.Printf("Author: %s\n", post.Author)
    }
    if post.Permalink != "" {
        log.Printf("Discussion URL: https://www.reddit.com%s\n", post.Permalink)
    }
    if post.Selftext != "" {
        log.Printf("Text Body: %s\n", post.Selftext)
    }

    if err := discord.SendToDiscord(target.Output.WebhookURL, embed); err != nil {
        log.Printf("Error sending to Discord webhook: %v", err)
    }
}
