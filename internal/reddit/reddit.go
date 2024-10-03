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
)

const (
    tokenURL   = "https://www.reddit.com/api/v1/access_token"
    apiURL     = "https://oauth.reddit.com/r/%s/%s"
    jsonAPIURL = "https://www.reddit.com/r/%s/%s.json"
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