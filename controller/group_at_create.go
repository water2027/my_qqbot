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

func HandleGroupAtCreate(c *gin.Context, payload *dto.Payload) {
	data := dto.GroupAtMessage{}
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
		fmt.Println("收到群@消息", data.GroupOpenId, data.Content)
		groupId := data.GroupOpenId
		msg := data.Content
		msgId := data.Id
		token := service.AuthHelper.GetToken()
		fmt.Println("获取token: ", token)
		resp, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/messages", groupId), dto.SendGroupMessage{
			Content: msg,
			MsgType: 0,
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
