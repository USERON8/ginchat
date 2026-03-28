// internal/ws/manager.go
package ws

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	"sync"
)

type Manager struct {
	clients map[uint]*Client
	mu      sync.RWMutex
}

var Mgr = &Manager{clients: make(map[uint]*Client)}

func (m *Manager) Add(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.UserID] = c
}

func (m *Manager) Remove(userID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[userID]; ok {
		close(c.Send) // 关闭 channel，触发 WritePump 退出
		delete(m.clients, userID)
	}
}

func (m *Manager) Send(toUserID uint, msg []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.clients[toUserID]; ok {
		// 对方在线，塞进他的 Send channel
		c.Send <- msg
	}
	// 对方不在线，消息已经存 DB 了，没问题
}
func (m *Manager) SendGroup(groupID uint, senderID uint, msg []byte) {
	// 查群成员
	var members []model.GroupMember
	database.DB.Where("group_id = ?", groupID).Find(&members)

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, member := range members {
		if member.UserID == senderID {
			continue // 不发给自己
		}
		if client, ok := m.clients[member.UserID]; ok {
			client.Send <- msg // 在线就推
		}
		// 不在线：消息已存 DB，查历史记录可以看到
	}
}
