package main

import (
	"fmt"
	"os"
	
	"io"

	"github.com/gin-gonic/gin"

	"qqbot/controller"
	"qqbot/dto"
	"qqbot/service"

	"github.com/joho/godotenv"
)

func webPush(c *gin.Context) {
	bytesData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("读取请求体失败", err)
	}
	fmt.Println("请求体", string(bytesData))
	var payload dto.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	switch payload.Opcode {
	case 0:
		// 服务器推送信息过来了
		switch payload.Type {
		case "GROUP_AT_MESSAGE_CREATE":
			controller.HandleGroupAtCreate(c, &payload)
		default:
			fmt.Println("未知消息类型", payload.Type)
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
