package message

import (
	"fmt"
)

type MessageType int

const (
	TextMessage     MessageType = 0
	MarkdownMessage MessageType = iota + 1
	ArkMessage
	EmbedMessage
	MediaMessage MessageType = 7
)

type MessageRouteType int

const (
	User MessageRouteType = iota
	Group
	Channel
	ChannelPrivate
)

type Message struct {
	msgId      string
	content    string
	rawContent string
	msgType    MessageType
	media      *MediaObject
	routeType  MessageRouteType // 通过哪一个渠道发送的
	routeId    string           // id
	hasSet     bool             // 是否设置过
}

func NewMessage(msgId, routeId string, routeType MessageRouteType, content string) *Message {
	fmt.Println("创建消息", msgId, routeId, routeType, content)
	return &Message{
		msgId:      msgId,
		routeId:    routeId,
		routeType:  routeType,
		rawContent: content,
	}
}

func (m *Message) ToStruct() interface{} {
	switch m.msgType {
	case TextMessage:
		return map[string]interface{}{
			"msg_id":   m.msgId,
			"content":  m.content,
			"msg_type": m.msgType,
		}
	case MediaMessage:
		return map[string]interface{}{
			"msg_id":   m.msgId,
			"msg_type": m.msgType,
			"media":    m.media.ToStruct(),
			"content":  m.content,
		}
	}

	return nil
}

func (m *Message) SetContent(content string) bool {
	if m.hasSet {
		return false
	}
	m.msgType = TextMessage
	m.content = content
	m.hasSet = true
	return true
}

func (m *Message) SetMedia(media *MediaObject, content string) bool {
	if m.hasSet {
		return false
	}
	m.msgType = MediaMessage
	m.media = media
	m.content = content
	m.hasSet = true
	return true
}

func (m *Message) GetRawContent() string {
	return m.rawContent
}

func (m *Message) GetRouteId() string {
	return m.routeId
}

func (m *Message) CanBeSet() bool {
	return m.hasSet
}
