package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
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
                Title string `json:"title"`
                URL   string `json:"url"`
            } `json:"data"`
        } `json:"children"`
    } `json:"data"`
}

func getAccessToken(oauthConfig *OAuthConfig) (string, error) {
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

func fetchRedditData(fetchConfig *FetchConfig, accessToken string, userAgent string, context string) (*RedditResponse, error) {
    client := &http.Client{}
    var req *http.Request
    var err error

    if context == "elevated" {
        req, err = http.NewRequest("GET", fmt.Sprintf(apiURL, fetchConfig.Subreddit, fetchConfig.Sorting), nil)
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    } else {
        req, err = http.NewRequest("GET", fmt.Sprintf(jsonAPIURL, fetchConfig.Subreddit, fetchConfig.Sorting), nil)
    }

    if err != nil {
        return nil, err
    }

    req.Header.Set("User-Agent", userAgent)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response code")
    }

    var redditResponse RedditResponse
    if err := json.NewDecoder(resp.Body).Decode(&redditResponse); err != nil {
        return nil, err
    }

    return &redditResponse, nil
}