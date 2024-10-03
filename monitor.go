package main

import (
	"fmt"
	"log"
	"time"
	"xenigo/internal/config"
	"xenigo/internal/discord"
	"xenigo/internal/reddit"
)

func monitorSubreddit(target config.Target, accessToken, userAgent, context string) {
	// Function to fetch and process Reddit data
	fetchAndProcess := func() {
		redditResponse, err := reddit.FetchRedditData(target, accessToken, userAgent, context)
		if err != nil {
			log.Printf("Error fetching Reddit data for subreddit %s: %v", target.Monitor.Subreddit, err)
			return
		}

		for _, child := range redditResponse.Data.Children {
			post := child.Data
			processAndSendPost(post, target)
		}
	}

	// Run immediately on start
	fetchAndProcess()

	// Set up the ticker for subsequent runs
	ticker := time.NewTicker(time.Duration(target.Options.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fetchAndProcess()
	}
}

func processAndSendPost(post reddit.RedditPost, target config.Target) {
	embed := discord.CreateDiscordEmbed(post, target)

	fmt.Printf("Subreddit: %s\nTitle: %s\nURL: %s\n", target.Monitor.Subreddit, post.Title, post.URL)
	if post.Author != "" {
		fmt.Printf("Author: %s\n", post.Author)
	}
	if post.Permalink != "" {
		fmt.Printf("Discussion URL: https://www.reddit.com%s\n", post.Permalink)
	}
	if post.Selftext != "" {
		fmt.Printf("Text Body: %s\n", post.Selftext)
	}
	fmt.Println()

	if err := discord.SendToDiscord(target.Output.WebhookURL, embed); err != nil {
		log.Printf("Error sending to Discord webhook: %v", err)
	}
}
