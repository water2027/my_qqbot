package plugin

import (
	"fmt"
	"qqbot/service"
)

func init() {
	service.MS.RegisterBeforeSendHook(service.BeforeSendHook{
		Priority: 255,
		Fn: func(msg *service.Message) error {
			success := msg.SetContent(fmt.Sprintf("[%s] %s", "HELLO", msg.GetRawContent()))
			if !success {
				return fmt.Errorf("设置消息内容失败")
			}
			return nil
		},
	})
}