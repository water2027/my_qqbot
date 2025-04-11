package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"qqbot/dto"
	"qqbot/service"
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
		fmt.Println("收到群@消息", msg.GroupOpenId, msg.Content)
		service.MS.ReceiveMessage(*service.NewMessage(msg.Id, msg.GroupOpenId, service.Group, msg.Content))
	}()

	c.JSON(200, gin.H{
		"message": "success",
	})
}
