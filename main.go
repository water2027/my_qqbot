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
	Type   string `json:"t"` // GROUP_AT_MESSAGE_CREATE | AT_MESSAGE_CREATE
}

type ValidationRequest struct {
	PlainToken string `json:"plain_token"`
	EventTs    string `json:"event_ts"`
}

type ValidationResponse struct {
	PlainToken string `json:"plain_token"`
	Signature  string `json:"signature"`
}

type Author struct {
	Id           string `json:"id"`
	Username     string `json:"username,omitempty"`
	Bot          bool   `json:"bot,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	MemberOpenId string `json:"member_openid,omitempty"`
	UnionOpenId  string `json:"union_openid,omitempty"`
}

// MessageScene defines message context information
type MessageScene struct {
	CallbackData string `json:"callback_data"`
	Source       string `json:"source"`
}

// Member represents information about a guild member
type Member struct {
	JoinedAt string `json:"joined_at"`
	Nick     string `json:"nick"`
	Roles    []int  `json:"roles"`
}

// Mention represents a user mentioned in a message
type Mention struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
	Avatar   string `json:"avatar"`
}

// GroupAtMessage represents messages from GROUP_AT_MESSAGE_CREATE events
type GroupAtMessage struct {
	Author       Author       `json:"author"`
	Content      string       `json:"content"`
	GroupId      string       `json:"group_id"`
	GroupOpenId  string       `json:"group_openid"`
	Id           string       `json:"id"`
	MessageScene MessageScene `json:"message_scene"`
	MessageType  int          `json:"message_type"`
	Timestamp    string       `json:"timestamp"`
}

// AtMessage represents messages from AT_MESSAGE_CREATE events
type AtMessage struct {
	Author       Author    `json:"author"`
	ChannelId    string    `json:"channel_id"`
	Content      string    `json:"content"`
	GuildId      string    `json:"guild_id"`
	Id           string    `json:"id"`
	Member       Member    `json:"member"`
	Mentions     []Mention `json:"mentions"`
	Seq          int       `json:"seq"`
	SeqInChannel int       `json:"seq_in_channel"`
	Timestamp    string    `json:"timestamp"`
}

func webPush(c *gin.Context) {
	var payload Payload
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
			var groupMsg GroupAtMessage
			data, err := json.Marshal(payload.Data)
			if err != nil {
				c.JSON(400, gin.H{"error": "Failed to marshal data"})
				return
			}
			if err := json.Unmarshal(data, &groupMsg); err != nil {
				c.JSON(400, gin.H{"error": "Failed to parse GROUP_AT_MESSAGE_CREATE"})
				return
			}
			fmt.Println("Group message author:", groupMsg.Author.Username)
			fmt.Println("Group message author ID:", groupMsg.Author.Id)
			fmt.Println("Group message author bot:", groupMsg.Author.Bot)
			fmt.Println("Group message author avatar:", groupMsg.Author.Avatar)
			fmt.Println("Group message author member open ID:", groupMsg.Author.MemberOpenId)
			fmt.Println("Group message author union open ID:", groupMsg.Author.UnionOpenId)
			fmt.Println("Group message content:", groupMsg.Content)
			fmt.Println("Group message group ID:", groupMsg.GroupId)
			fmt.Println("Group message timestamp:", groupMsg.Timestamp)
			fmt.Println("Group message ID:", groupMsg.Id)
			fmt.Println("Group message message type:", groupMsg.MessageType)
			fmt.Println("Group message message scene:", groupMsg.MessageScene.CallbackData)
			fmt.Println("Group message message scene source:", groupMsg.MessageScene.Source)

		case "AT_MESSAGE_CREATE":
			var atMsg AtMessage
			data, err := json.Marshal(payload.Data)
			if err != nil {
				c.JSON(400, gin.H{"error": "Failed to marshal data"})
				return
			}
			if err := json.Unmarshal(data, &atMsg); err != nil {
				c.JSON(400, gin.H{"error": "Failed to parse AT_MESSAGE_CREATE"})
				return
			}
			// Process at message
			fmt.Println("At message author:", atMsg.Author.Username)
			fmt.Println("At message author ID:", atMsg.Author.Id)
			fmt.Println("At message author bot:", atMsg.Author.Bot)
			fmt.Println("At message author avatar:", atMsg.Author.Avatar)
			fmt.Println("At message author member open ID:", atMsg.Author.MemberOpenId)
			fmt.Println("At message author union open ID:", atMsg.Author.UnionOpenId)
			fmt.Println("At message content:", atMsg.Content)
			fmt.Println("At message channel ID:", atMsg.ChannelId)
			fmt.Println("At message guild ID:", atMsg.GuildId)
			fmt.Println("At message ID:", atMsg.Id)
			fmt.Println("At message member joined at:", atMsg.Member.JoinedAt)
			fmt.Println("At message member nick:", atMsg.Member.Nick)
			fmt.Println("At message member roles:", atMsg.Member.Roles)
			fmt.Println("At message mentions:", atMsg.Mentions)
			fmt.Println("At message mentions ID:", atMsg.Mentions[0].Id)
			fmt.Println("At message mentions username:", atMsg.Mentions[0].Username)
			fmt.Println("At message mentions bot:", atMsg.Mentions[0].Bot)
			fmt.Println("At message mentions avatar:", atMsg.Mentions[0].Avatar)
			fmt.Println("At message timestamp:", atMsg.Timestamp)
			fmt.Println("At message seq:", atMsg.Seq)
			fmt.Println("At message seq in channel:", atMsg.SeqInChannel)
			
		}
	case 13:
		// 服务器验证
		validate(c, &payload)
	}
}

func validate(c *gin.Context, payload *Payload) {
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
	r.POST("/qqbot", webPush)
	r.Run(":8080")
}
