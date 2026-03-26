// internal/handler/ws.go
package handler

import (
	"ginchat/internal/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	ws.Mgr.Add(client)

	// 两个 goroutine 分别负责读和写
	go client.WritePump() // 新开一个，不阻塞
	client.ReadPump()     // 当前 goroutine 阻塞在这里读
}
