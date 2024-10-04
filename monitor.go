package main

import (
	"log"
	"time"
	"xenigo/internal/config"
	"xenigo/internal/notifier"
	"xenigo/internal/reddit"
)

func monitorSubreddit(target config.Target, accessToken, userAgent, context string, cache *Cache, sendInitial bool, devFlags *config.DeveloperFlags) {
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
			if !cache.IsProcessed(post.Permalink) || config.GetFlag(devFlags.IgnoreCache){
				if sendToDiscord {
					notifier.ProcessAndSendPost(post, target, devFlags)
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
