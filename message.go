package nsq

import "github.com/nsqio/go-nsq"

type Message struct {
	*nsq.Message

	Channel  string
	Topic    string
	Delegate MessageDelegate
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

func (m *Message) Clone() *Message {
	cloned := *m
	return &cloned
}
