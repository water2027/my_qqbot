package dto

type ValidationRequest struct {
	PlainToken string `json:"plain_token"`
	EventTs    string `json:"event_ts"`
}

type ValidationResponse struct {
	PlainToken string `json:"plain_token"`
	Signature  string `json:"signature"`
}