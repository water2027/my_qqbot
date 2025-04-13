package homeworkddl

import (
	"fmt"
	"io"
	"os"
	"qqbot/service/message"
	"qqbot/utils"
	"strings"
)

func getHomeworkDDL(msg *message.Message) error {
	if !msg.CanBeSet() {
		return nil
	}
	rawContent := msg.GetRawContent()
	cmd, found := strings.CutPrefix(rawContent, " /作业")
	if !found {
		return nil
	}
	if cmd == "" {
		// 如果cmd为空，那么是查询作业
		return nil
	}

	// 非空，那么是设置作业ddl
	// 设置作业格式：
	// 2025-4-13-22 作业
	ddl, work, found := strings.Cut(cmd, " ")
	if !found {
		msg.SetContent("格式错误！格式为2025-4-13-22 作业")
		return nil
	}

	times := strings.Split(ddl, "-")
	if len(times) < 4 {
		msg.SetContent("格式错误！格式为2025-4-13-22 作业")
		return nil
	}

	year := times[0]
	month := times[1]
	day := times[2]
	hour := times[3]
	
	webhook := os.Getenv("WEBHOOK")
	resp, err := utils.NetHelper.POST(webhook, map[string]interface{}{
		"year":year,
		"month":month,
		"day":day,
		"hour":hour,
		"content":work,
	})

	if err != nil {
		msg.SetContent("设置失败！网络错误")
		return nil
	}

	defer resp.Body.Close()
	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		msg.SetContent("设置失败！网络错误")
		return nil
	}
	if string(bytesData) != "success" {
		msg.SetContent("设置失败！")
		return nil
	}

	// 发给我的企微群机器人
	msg.SetContent("设置成功！")

	return nil
}

func init() {
	message.MS.RegisterBeforeSendHook(message.BeforeSendHook{
		Priority: 200,
		Fn: func(msg *message.Message) error {
			if !msg.CanBeSet() {
				return nil
			}
			success := msg.SetContent(fmt.Sprintf("[%s] %s", "HELLO", msg.GetRawContent()))
			if !success {
				return fmt.Errorf("设置消息内容失败")
			}
			return nil
		},
	})
}