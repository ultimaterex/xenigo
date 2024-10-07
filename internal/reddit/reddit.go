package reddit

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
    "sync"
    "time"
    "xenigo/internal/config"
)

const (
    tokenURL        = "https://www.reddit.com/api/v1/access_token"
    apiURL          = "https://oauth.reddit.com/r/%s/%s"
    jsonAPIURL      = "https://www.reddit.com/r/%s/%s.json"
    defaultRetryCount = 3
    defaultRetryInterval = 2
    defaultRetryIntervalSeconds = 2 * time.Second
)

type RedditResponse struct {
    Data struct {
        Children []struct {
            Data RedditPost `json:"data"`
        } `json:"children"`
    } `json:"data"`
}

type RedditPost struct {
    Title       string `json:"title"`
    URL         string `json:"url"`
    Author      string `json:"author"`
    Permalink   string `json:"permalink"`
    Selftext    string `json:"selftext"`
    Stickied    bool   `json:"stickied"`
}

var (
    tokenMutex      sync.Mutex
    lastRefreshTime time.Time
    refreshInterval = 1 * time.Minute // Set the minimum interval between token refreshes
)

func GetAccessToken(oauthConfig *config.OAuthConfig) (string, error) {
    data := url.Values{}
    data.Set("grant_type", "password")
    data.Set("username", oauthConfig.Username)
    data.Set("password", oauthConfig.Password)

    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.SetBasicAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)
    req.Header.Set("User-Agent", "xenigo")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    token, ok := result["access_token"].(string)
    if !ok {
        return "", fmt.Errorf("access token not found in response")
    }

    return token, nil
}

func refreshAccessToken(oauthConfig *config.OAuthConfig) (string, error) {
    log.Println("Refreshing access token...")
    newToken, err := GetAccessToken(oauthConfig)
    if err != nil {
        return "", fmt.Errorf("failed to refresh access token: %w", err)
    }
    return newToken, nil
}

func FetchRedditData(target config.Target, accessToken string, userAgent string, context string, oauthConfig *config.OAuthConfig) (*RedditResponse, error) {
    client := &http.Client{
        Timeout: 10 * time.Second, // Set a timeout for the HTTP client
    }
    url := fmt.Sprintf(jsonAPIURL, target.Monitor.Subreddit, target.Monitor.Sorting)
    if context == "elevated" {
        url = fmt.Sprintf(apiURL, target.Monitor.Subreddit, target.Monitor.Sorting)
    }
    // Add limit parameter to the URL
    url = fmt.Sprintf("%s?limit=%d", url, target.Options.Limit)
    var redditResponse RedditResponse
    retries := target.Options.RetryCount
    if retries == 0 {
        retries = defaultRetryCount // Default retry count
    }
    retryInterval := target.Options.RetryInterval
    if retryInterval == 0 {
        retryInterval = defaultRetryInterval // Default retry interval
    }
    for i := 0; i < retries; i++ {
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            return nil, fmt.Errorf("failed to create request: %w", err)
        }
        if context == "elevated" {
            req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
        }
        req.Header.Set("User-Agent", userAgent)
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Attempt %d: Error fetching Reddit data: %v", i+1, err)
            time.Sleep(defaultRetryIntervalSeconds) // Wait before retrying
            continue
        }
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusUnauthorized {
            // Refresh the access token and retry
            tokenMutex.Lock()
            if time.Since(lastRefreshTime) < refreshInterval {
                tokenMutex.Unlock()
                log.Println("Token was recently refreshed, waiting before retrying...")
                time.Sleep(refreshInterval - time.Since(lastRefreshTime))
                continue
            }
            accessToken, err = refreshAccessToken(oauthConfig)
            lastRefreshTime = time.Now()
            tokenMutex.Unlock()
            if err != nil {
                return nil, err
            }
            continue
        }
        if resp.StatusCode != http.StatusOK {
            bodyBytes, _ := io.ReadAll(resp.Body)
            bodyString := string(bodyBytes)
            log.Printf("Error: received non-200 response code: %d, body: %s", resp.StatusCode, bodyString)
            return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
        }
        if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
            return nil, fmt.Errorf("failed to decode response: %w", err)
        }
        // Filter out pinned modposts
        filteredChildren := []struct {
            Data RedditPost `json:"data"`
        }{}
        for _, child := range redditResponse.Data.Children {
            if !child.Data.Stickied {
                filteredChildren = append(filteredChildren, child)
            }
        }
        redditResponse.Data.Children = filteredChildren
        return &redditResponse, nil
    }
    return nil, fmt.Errorf("failed to fetch Reddit data after %d attempts", retries)
}