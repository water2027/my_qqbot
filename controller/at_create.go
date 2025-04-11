package controller

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"qqbot/dto"
	"qqbot/service"
	"qqbot/utils"
)

func HandleAtCreate(c *gin.Context, payload *dto.Payload) {
	data := dto.AtMessage{}
	payloadBytes, err := json.Marshal(payload.Data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	go func() {
		fmt.Println("收到频道@消息", data.Id, data.Content)
		msgId := data.Id
		msg := data.Content
		token := service.AuthHelper.GetToken()
		fmt.Println("获取token: ", token)
		resp, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/channels/%s/messages", msgId), dto.SendChannelMessage{
			Content: msg,
			MsgId:   msgId,
		}, utils.WithToken(token))
		if err != nil {
			fmt.Println("发送消息失败", err)
			return
		}
		res, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("读取响应失败", err)
			return
		}

		fmt.Println("发送消息", string(res), err)
	}()

	c.JSON(200, gin.H{
		"message": "success",
	})
}
