package slack

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "xenigo/internal/output"
)

type SlackMessage struct {
    Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
    Title  string       `json:"title"`
    Text   string       `json:"text"`
    Fields []SlackField `json:"fields,omitempty"`
}

type SlackField struct {
    Title string `json:"title"`
    Value string `json:"value"`
}

type SlackSender struct {
    WebhookURL string
}

func (s *SlackSender) SendMessage(embed output.MessageEmbed) error {
    log.Printf("Sending message to Slack: %s", embed.Title) // Log statement

    slackAttachment := SlackAttachment{
        Title:  embed.Title,
        Text:   embed.Description,
        Fields: convertFields(embed.Fields),
    }

    message := SlackMessage{Attachments: []SlackAttachment{slackAttachment}}
    messageBody, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message body: %w", err)
    }

    resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(messageBody))
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
    }

    return nil
}

func convertFields(fields []output.EmbedField) []SlackField {
    var slackFields []SlackField
    for _, field := range fields {
        slackFields = append(slackFields, SlackField{
            Title: field.Name,
            Value: field.Value,
        })
    }
    return slackFields
}