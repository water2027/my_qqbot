package dto

type AtMessage struct {
	Author       Author    `json:"author"`
	ChannelId    string    `json:"channel_id"`
	Content      string    `json:"content"`
	GuildId      string    `json:"guild_id"`
	Id           string    `json:"id"`
	Member       Member    `json:"member"`
	Seq          int       `json:"seq"`
	Timestamp    string    `json:"timestamp"`
}