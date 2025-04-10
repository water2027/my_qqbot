package controller

import (
	"qqbot/dto"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"crypto/ed25519"
	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context, payload *dto.Payload) {
	validationPayload := &dto.ValidationRequest{}
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

	c.JSON(200, dto.ValidationResponse{
		PlainToken: validationPayload.PlainToken,
		Signature:  signature,
	})
}