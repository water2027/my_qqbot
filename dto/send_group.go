package dto

type SendGroupMessage struct {
	Content string `json:"content"`
	MsgType int    `json:"msg_type"`
	MsgId   string `json:"msg_id"`
}
