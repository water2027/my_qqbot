package dto

type Author struct {
	Id           string `json:"id"`
	Username     string `json:"username,omitempty"`
	Bot          bool   `json:"bot,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	MemberOpenId string `json:"member_openid,omitempty"`
	UnionOpenId  string `json:"union_openid,omitempty"`
}

type Member struct {
	JoinedAt string `json:"joined_at"`
	Nick     string `json:"nick"`
	Roles    []int  `json:"roles"`
}

type Payload struct {
	Id     string `json:"id"`
	Opcode int    `json:"op"` // 0 13
	Data   any    `json:"d"`
	S      int    `json:"s"`
	Type   string `json:"t"` // GROUP_AT_MESSAGE_CREATE | AT_MESSAGE_CREATE
}

