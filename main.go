package main

import (
    "fmt"
    "os"
)

func main() {
    config, err := loadConfig("config.yaml")
    if err != nil {
        fmt.Println("Error reading config:", err)
        os.Exit(1)
    }

    context := determineContext(config)

    var accessToken string
    if context == "elevated" {
        accessToken, err = getAccessToken(config.OAuth)
        if err != nil {
            fmt.Println("Error getting access token:", err)
            os.Exit(1)
        }
    }

    for _, fetchConfig := range config.Fetches {
        redditResponse, err := fetchRedditData(&fetchConfig, accessToken, config.UserAgent, context)
        if err != nil {
            fmt.Println("Error fetching Reddit data:", err)
            continue
        }

        if len(redditResponse.Data.Children) > 0 {
            post := redditResponse.Data.Children[0].Data
            fmt.Printf("Subreddit: %s\nTitle: %s\nURL: %s\n", fetchConfig.Subreddit, post.Title, post.URL)
        }
    }
}