package notifier

import (
    "fmt"
    "log"
    "xenigo/internal/config"
    "xenigo/internal/discord"
    "xenigo/internal/reddit"
    "xenigo/internal/slack"
    "xenigo/internal/output"
)

func ProcessAndSendPost(post reddit.RedditPost, target config.Target, devFlags *config.DeveloperFlags) {
    if config.GetFlag(devFlags.NotifyMute) {
        log.Printf("Notifications are muted for target: %s", target.Name)
        return
    }

    embed := output.MessageEmbed{
        Title:       post.Title,
        Description: post.Selftext,
        URL:         post.URL,
        Author:      post.Author,
        Fields: []output.EmbedField{
            {Name: "Subreddit", Value: target.Monitor.Subreddit},
            {Name: "Discussion URL", Value: fmt.Sprintf("https://www.reddit.com%s", post.Permalink)},
        },
    }

    var sender output.MessageSender
    log.Printf("Processing target with output type: %s", target.Output.Type) // Add this line for debugging
    switch target.Output.Type {
    case config.OutputTypeDiscord:
        sender = &discord.DiscordSender{WebhookURL: target.Output.WebhookURL}
    case config.OutputTypeSlack:
        sender = &slack.SlackSender{WebhookURL: target.Output.WebhookURL}
    default:
        log.Printf("Unsupported output type: %s", target.Output.Type)
        return
    }

    if err := sender.SendMessage(embed); err != nil {
        log.Printf("Error sending message: %v", err)
    }
}