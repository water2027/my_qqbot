package controller

import (
	"encoding/json"
	"fmt"

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
		token := service.AuthHelper.GetToken()
		fmt.Println("获取token: ", token)
		_, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/messages", groupId), dto.SendGroupMessage{
			Content: msg,
			MsgType: 0,
		}, utils.WithToken(token))
		fmt.Println("发送消息", err)
	}()

	c.JSON(200, gin.H{
		"message": "success",
	})
}
