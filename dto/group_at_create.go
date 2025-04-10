package dto

type MessageScene struct {
	CallbackData string `json:"callback_data"`
	Source       string `json:"source"`
}

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