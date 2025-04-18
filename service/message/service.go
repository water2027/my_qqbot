package message

import (
	"context"
	"fmt"
	"io"
	"time"

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
	// 4分钟+30秒的超时，如果超过这个时间还没有发送完，就直接返回
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute+30*time.Second)
	defer cancel()
	errChan := make(chan error, 1)
	go func() {
		for _, hook := range ms.beforeSend {
			if err := hook.Fn(&msg); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- nil
	}()
	select {
	case <-ctx.Done():
		// 如果设置过了，那就继续发送
		// 没有就设置超时
		if !msg.hasSet {
			msg.SetContent("消息准备超时")
		}
	case err := <-errChan:
		if err != nil {
			fmt.Println("消息发送失败", err, msg.routeId, msg.content, msg.msgType, msg.msgId)
			return
		}
	}
	switch msg.routeType {
	case Group:
		token := auth.AuthHelper.GetToken()
		fmt.Println("发送消息", msg.routeId, msg.content, msg.msgType, msg.msgId)
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
