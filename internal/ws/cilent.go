package ws

import (
	"encoding/json"
	"ginchat/internal/model"
	"ginchat/pkg/database"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

// ReadPump ：一直读，收到消息就处理
func (c *Client) ReadPump() {
	// 函数退出时：关闭连接，从 Manager 移除
	defer func() {
		Mgr.Remove(c.UserID)
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		// 阻塞等待客户端发消息
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			// 客户端断开了（关闭浏览器/网络断了）
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue // 格式不对，跳过
		}

		// 处理消息
		c.handleMessage(msg, data)
	}
}

// WritePump ：一直等 channel，有消息就发出去
func (c *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {

		}
	}(c.Conn)

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				// channel 被关闭了，说明用户下线
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			// 发给客户端
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

// handleMessage：消息处理逻辑
func (c *Client) handleMessage(msg Message, raw []byte) {
	switch msg.Type {
	case "private":
		// 1. 存数据库
		record := model.PrivateMessage{
			FromID:  c.UserID,
			ToID:    msg.ToID,
			Content: msg.Content,
		}
		database.DB.Create(&record)

		// 2. 找到对方，发过去
		Mgr.Send(msg.ToID, raw)

	case "ping":
		// 心跳，回一个 pong
		pong, _ := json.Marshal(Message{Type: "pong"})
		c.Send <- pong
	}
}
