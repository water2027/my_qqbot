package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"qqbot/controller"
	"qqbot/dto"
	"qqbot/service"

	"github.com/joho/godotenv"
)

func webPush(c *gin.Context) {
	var payload dto.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	switch payload.Opcode {
	case 0:
		// 服务器推送信息过来了
		fmt.Println("Received type:", payload.Type)
		fmt.Println("Received message:", payload.Data)
		switch payload.Type {
		case "GROUP_AT_MESSAGE_CREATE":
			controller.HandleGroupAtCreate(c, &payload)
		default:
			if strings.Contains(payload.Type, "AT_MESSAGE_CREATE") {
				controller.HandleAtCreate(c, &payload)
			}
	case 13:
		// 服务器验证
		controller.Validate(c, &payload)
	}
}

func main() {
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		godotenv.Load()
	}

	r := gin.Default()
	r.POST("/qqbot", webPush)
	service.Init()
	r.Run(":8080")
}
