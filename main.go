package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"qqbot/controller"
	"qqbot/dto"

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

		case "AT_MESSAGE_CREATE":
			var atMsg dto.AtMessage
			data, err := json.Marshal(payload.Data)
			if err != nil {
				fmt.Println("Failed to marshal data:", err)
				c.JSON(400, gin.H{"error": "Failed to marshal data"})
				return
			}
			if err := json.Unmarshal(data, &atMsg); err != nil {
				fmt.Println("Failed to unmarshal data:", err)
				c.JSON(400, gin.H{"error": "Failed to parse AT_MESSAGE_CREATE"})
				return
			}
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
	r.Run(":8080")
}
