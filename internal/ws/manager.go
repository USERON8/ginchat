// internal/ws/manager.go
package ws

import "sync"

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
