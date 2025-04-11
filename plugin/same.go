package plugin

import (
	"fmt"
	"qqbot/service/message"
)

func init() {
	message.MS.RegisterBeforeSendHook(message.BeforeSendHook{
		Priority: 255,
		Fn: func(msg *message.Message) error {
			success := msg.SetContent(fmt.Sprintf("[%s] %s", "HELLO", msg.GetRawContent()))
			if !success {
				return fmt.Errorf("设置消息内容失败")
			}
			return nil
		},
	})
}