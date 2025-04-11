package ai

import (
	"io"
	"os"
	"encoding/json"
	"qqbot/service/message"
	"qqbot/utils"
)

type GLMResponse struct {
	Created   int64    `json:"created"`
	ID        string   `json:"id"`
	Model     string   `json:"model"`
	RequestID string   `json:"request_id"`
	Choices   []Choice `json:"choices"`
	Usage     Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func getAiResponse(msg *message.Message) error {
	if !msg.CanBeSet() {
		return nil
	}

	apiUrl := os.Getenv("AI_API_URL")
	apiKey := os.Getenv("AI_API_KEY")
	if apiUrl == "" || apiKey == "" {
		msg.SetContent("AI API URL or API Key is not set")
		return nil
	}
	resp, err := utils.NetHelper.POST(apiUrl, map[string]interface{}{
		"model": "glm-4-flash",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": msg.GetRawContent(),
			},
		},
	}, utils.WithToken(apiKey))
	if err != nil {
		msg.SetContent("AI API request failed: " + err.Error())
		return nil
	}
	defer resp.Body.Close()
	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		msg.SetContent("AI API response read failed: " + err.Error())
		return nil
	}

	var glmResponse GLMResponse
	err = json.Unmarshal(bytesData, &glmResponse)
	if err != nil {
		msg.SetContent("AI API response parse failed: " + err.Error())
		return nil
	}

	if len(glmResponse.Choices) == 0 {
		msg.SetContent("AI API response is empty")
		return nil
	}

	msg.SetContent(glmResponse.Choices[0].Message.Content)
	return nil
}

func init() {
	message.MS.RegisterBeforeSendHook(message.BeforeSendHook{
		Priority: 50,
		Fn:       getAiResponse,
	})
}
