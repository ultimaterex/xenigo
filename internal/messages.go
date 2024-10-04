
package messages

type MessageSender interface {
    SendMessage(embed MessageEmbed) error
}

type MessageEmbed struct {
    Title       string
    Description string
    URL         string
    Author      string
    Fields      []EmbedField
}

type EmbedField struct {
    Name  string
    Value string
}