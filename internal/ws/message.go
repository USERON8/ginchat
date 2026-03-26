// internal/ws/message.go
package ws

type Message struct {
	Type    string `json:"type"` // private/ping/pong
	ToID    uint   `json:"toId"`
	Content string `json:"content"`
	FromID  uint   `json:"fromId"` // 服务端转发时填上
}
