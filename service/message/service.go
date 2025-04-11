package message

import (
	"fmt"
	"io"

	"qqbot/service/auth"
	"qqbot/utils"
)

type MessageService struct {
	beforeSend []BeforeSendHook
	msgChan    chan Message
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
		token := auth.AuthHelper.GetToken()
		resp, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/messages", msg.routeId), msg.ToStruct(), utils.WithToken(token))
		if err != nil {
			fmt.Println("发送消息失败", err)
		}
		defer resp.Body.Close()
		bytesData, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取数据失败", err)
		}
		fmt.Println("发送消息成功", string(bytesData))
	default:
		return
	}
}

func (ms *MessageService) Run() error {
	for msg := range ms.msgChan {
		if msg.msgId == "" || msg.routeId == "" {
			continue
		}
		ms.sendMessage(msg)
	}
	return nil
}

func (ms *MessageService) RegisterBeforeSendHook(hook BeforeSendHook) {
	// 查找插入位置
	insertIndex := 0
	for i, existingHook := range ms.beforeSend {
		if hook.Priority < existingHook.Priority {
			insertIndex = i + 1
		} else {
			break
		}
	}

	// 在指定位置插入新的 hook
	if insertIndex == len(ms.beforeSend) {
		// 如果应该添加到末尾，直接追加
		ms.beforeSend = append(ms.beforeSend, hook)
	} else {
		// 否则，在中间插入
		ms.beforeSend = append(ms.beforeSend[:insertIndex+1], ms.beforeSend[insertIndex:]...)
		ms.beforeSend[insertIndex] = hook
	}
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
