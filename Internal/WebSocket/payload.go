package websocket

type MessagePayload struct {
	MatchID    uint   `json:"matchID"`
	FromUserID uint   `json:"fromUserID"`
	ToUserID   uint   `json:"toUserID"`
	Message    string `json:"message"`
}
