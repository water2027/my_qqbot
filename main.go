package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"

	"crypto/ed25519"
)

type Payload struct {
	Id     string `json:"id"`
	Opcode int    `json:"op"` // 0 12 13
	Data   any    `json:"d"`
	S      int    `json:"s"`
	Type   string `json:"t"`
}

type ValidationRequest struct {
	PlainToken string `json:"plain_token"`
	EventTs    string `json:"event_ts"`
}

type ValidationResponse struct {
	PlainToken string `json:"plain_token"`
	Signature  string `json:"signature"`
}

func validate(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	validationPayload := &ValidationRequest{}
    dataBytes, err := json.Marshal(payload.Data)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Then unmarshal those bytes into the validationPayload
    if err := json.Unmarshal(dataBytes, validationPayload); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }


	seed := os.Getenv("APP_SECRET")
	if seed == "" {
		c.JSON(500, gin.H{"error": "APP_SECRET not set"})
		return
	}

	for len(seed) < ed25519.SeedSize {
		seed = strings.Repeat(seed, 2)
	}
	seed = seed[:ed25519.SeedSize]
	reader := strings.NewReader(seed)

	_, privateKey, err := ed25519.GenerateKey(reader)

	if err != nil {
		fmt.Println("ed25519 generate key failed:", err)
		return
	}

	var msg bytes.Buffer
	msg.WriteString(validationPayload.EventTs)
	msg.WriteString(validationPayload.PlainToken)

	signature := hex.EncodeToString(ed25519.Sign(privateKey, msg.Bytes()))

	c.JSON(200, ValidationResponse{
		PlainToken: validationPayload.PlainToken,
		Signature:  signature,
	})

}

func main() {
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		godotenv.Load()
	}

	r := gin.Default()
	r.POST("/qqbot", validate)
	r.Run(":8080")
}
