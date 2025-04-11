package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"qqbot/dto"
	"qqbot/service/message"
)

func HandleGroupAtCreate(c *gin.Context, payload *dto.Payload) {
	msg := dto.GroupAtMessage{}
	payloadBytes, err := json.Marshal(payload.Data)
	if err != nil {
		fmt.Println("解析请求体失败:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(payloadBytes, &msg); err != nil {
		fmt.Println("解析请求体失败:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	go func() {
		message.MS.ReceiveMessage(*message.NewMessage(msg.Id, msg.GroupOpenId, message.Group, msg.Content))
	}()

	c.JSON(200, gin.H{
		"message": "success",
	})
}
