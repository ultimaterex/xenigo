package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"xenigo/internal/config"
	"xenigo/internal/discord"
)

const (
	tokenURL   = "https://www.reddit.com/api/v1/access_token"
	apiURL     = "https://oauth.reddit.com/r/%s/%s"
	jsonAPIURL = "https://www.reddit.com/r/%s/%s.json"
)

type RedditResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Title       string `json:"title"`
				URL         string `json:"url"`
				Author      string `json:"author"`
				Permalink   string `json:"permalink"`
				Selftext    string `json:"selftext"`
				Stickied    bool   `json:"stickied"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetAccessToken(oauthConfig *config.OAuthConfig) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", oauthConfig.Username)
	data.Set("password", oauthConfig.Password)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)
	req.Header.Set("User-Agent", "xenigo")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

func FetchRedditData(target config.Target, accessToken string, userAgent string, context string) (*RedditResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf(jsonAPIURL, target.Monitor.Subreddit, target.Monitor.Sorting)
	if context == "elevated" {
		url = fmt.Sprintf(apiURL, target.Monitor.Subreddit, target.Monitor.Sorting)
	}

	// Add limit parameter to the URL
	url = fmt.Sprintf("%s?limit=%d", url, target.Options.Limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if context == "elevated" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Error: received non-200 response code: %d, body: %s", resp.StatusCode, bodyString)
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var redditResponse RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
		return nil, err
	}

	// Filter out pinned modposts
	filteredChildren := []struct {
		Data struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Author      string `json:"author"`
			Permalink   string `json:"permalink"`
			Selftext    string `json:"selftext"`
			Stickied    bool   `json:"stickied"`
		} `json:"data"`
	}{}

	for _, child := range redditResponse.Data.Children {
		if !child.Data.Stickied {
			filteredChildren = append(filteredChildren, child)
		}
	}

	redditResponse.Data.Children = filteredChildren

	return &redditResponse, nil
}

func FetchData(appConfig *config.AppConfig) {
	var accessToken string
	var err error

	if appConfig.Context == config.ContextElevated {
		accessToken, err = GetAccessToken(appConfig.Config.OAuth)
		if err != nil {
			log.Fatalf("Error getting access token: %v", err)
		}
	}

	for _, target := range appConfig.Config.Targets {
		redditResponse, err := FetchRedditData(target, accessToken, appConfig.Config.UserAgent, string(appConfig.Context))
		if err != nil {
			log.Printf("Error fetching Reddit data for subreddit %s: %v", target.Monitor.Subreddit, err)
			continue
		}

		for _, child := range redditResponse.Data.Children {
			post := child.Data
			message := fmt.Sprintf("Subreddit: %s\nTitle: %s\nURL: %s\n", target.Monitor.Subreddit, post.Title, post.URL)
			
			if post.Author != "" {
				message += fmt.Sprintf("Author: %s\n", post.Author)
			}
			if post.Permalink != "" {
				message += fmt.Sprintf("Discussion URL: https://www.reddit.com%s\n", post.Permalink)
			}
			if post.Selftext != "" {
				message += fmt.Sprintf("Text Body: %s\n", post.Selftext)
			}
			message += "\n"

			fmt.Print(message)

			if err := discord.NotifyDiscord(target.Output.WebhookURL, message); err != nil {
				log.Printf("Error sending to Discord webhook: %v", err)
			}
		}
	}
}
