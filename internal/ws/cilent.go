package ws

import (
	"encoding/json"
	"ginchat/pkg/logger"
	"ginchat/pkg/redis"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		redis.SetOffline(c.UserID)
		Mgr.Remove(c.UserID)
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	redis.SetOnline(c.UserID)

	err := c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.handleMessage(data) // ← 调这里
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 在同一个文件最底部
func (c *Client) handleMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	// 去重
	if msg.MsgID != "" && redis.IsDuplicate(msg.MsgID) {
		return
	}
	logger.Info("收到消息",
		zap.Uint("fromID", c.UserID),
		zap.Uint("toID", msg.ToID),
		zap.String("content", msg.Content),
	)
	msg.FromID = c.UserID
	switch msg.Type {
	case "private":
		// 异步写DB
		SubmitMsg(MsgTask{
			FromID:  c.UserID,
			ToID:    msg.ToID,
			Content: msg.Content,
		})
		newData, _ := json.Marshal(msg)
		Mgr.Send(msg.ToID, newData)
	case "group":
		SubmitMsg(MsgTask{
			Type:    "group",
			FromID:  c.UserID,
			ToID:    msg.ToID, // 群ID
			Content: msg.Content,
		})
		Mgr.SendGroup(msg.ToID, c.UserID, data)
		// 立刻转发

	case "ping":
		pong, _ := json.Marshal(Message{Type: "pong"})
		c.Send <- pong
	}
}
