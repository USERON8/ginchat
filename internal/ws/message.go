// internal/ws/message.go
package ws

type Message struct {
	MsgID   string `json:"msgId"` // 客户端生成，用于去重
	Type    string `json:"type"`  // private/ping/pong
	ToID    uint   `json:"toId"`
	Content string `json:"content"`
	FromID  uint   `json:"fromId"` // 服务端填
}
