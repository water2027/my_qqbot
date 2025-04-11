package service

import (
	"fmt"
	"sort"

	"qqbot/utils"
)

type BeforeSendFunc func(msg *Message) error

type BeforeSendHook struct {
	Fn BeforeSendFunc
	// 优先级
	Priority int
}

type MessageService struct {
	beforeSend []BeforeSendHook
	msgChan    chan Message
}

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

func (m *Message) SetMedia(media *MediaObject) bool {
	if m.hasSet {
		return false
	}
	m.msgType = MediaMessage
	m.media = media
	m.hasSet = true
	return true
}

func (m *Message) GetRawContent() string {
	return m.rawContent
}

type MediaObject struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
}

func (m *MediaObject) ToStruct() interface{} {
	return map[string]interface{}{
		"file_uuid": m.FileUUID,
		"file_info": m.FileInfo,
		"ttl":       m.TTL,
	}
}

func NewMessage(msgId, routeId string, routeType MessageRouteType, content string) *Message {
	return &Message{
		msgId:     msgId,
		routeId:   routeId,
		routeType: routeType,
		rawContent: content,
	}
}

func (ms *MessageService) ReceiveMessage(msg Message) {
	ms.msgChan <- msg
}

func (ms *MessageService) sendMessage(msg Message) {
	for _, hook := range ms.beforeSend {
		if err := hook.Fn(&msg); err != nil {
			fmt.Println("发送消息前处理失败", err)
			return
		}
	}
	switch msg.routeType {
	case Group:
		token := AuthHelper.GetToken()
		// 默认情况
		_, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/messages", msg.routeId), msg.ToStruct(), utils.WithToken(token))
		if err != nil {
			fmt.Println("发送消息失败", err)
		}
	default:
		return
	}
}

func (ms *MessageService) Run() error {
	// 给hooks排序
	sort.Slice(ms.beforeSend, func(i, j int) bool {
		return ms.beforeSend[i].Priority < ms.beforeSend[j].Priority
	})
	for msg := range ms.msgChan {
		if msg.msgId == "" || msg.routeId == "" {
			continue
		}
		ms.sendMessage(msg)
	}
	return nil
}

func (ms *MessageService) RegisterBeforeSendHook(hook BeforeSendHook) {
	ms.beforeSend = append(ms.beforeSend, hook)
}

var MS *MessageService

func init() {
	MS = &MessageService{
		msgChan:    make(chan Message, 100),
		beforeSend: make([]BeforeSendHook, 0),
	}
	go MS.Run()
	fmt.Println("消息服务启动成功")
}
