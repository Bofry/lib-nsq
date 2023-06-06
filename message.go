package nsq

import "github.com/nsqio/go-nsq"

type Message struct {
	*nsq.Message

	Topic string
}

func (m *Message) Content() *MessageContent {
	content, _ := DecodeMessageContent(m.Body)
	if content != nil {
		return content
	}
	return &MessageContent{
		Body: m.Body,
	}
}
