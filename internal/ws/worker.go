// internal/ws/worker.go
package ws

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
)

// 消息任务
type MsgTask struct {
	FromID  uint
	ToID    uint
	Content string
	Type    string
}

// 全局消息队列
var msgQueue = make(chan MsgTask, 1024) // 能缓冲 1024 条

// 启动 Worker Pool（开 N 个 goroutine 消费队列）
func StartWorkers(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for task := range msgQueue {
				switch task.Type {
				case "private":
					database.DB.Create(&model.PrivateMessage{
						FromID:  task.FromID,
						ToID:    task.ToID,
						Content: task.Content,
					})
				case "group":
					database.DB.Create(&model.GroupMessage{
						GroupID: task.ToID,
						FromID:  task.FromID,
						Content: task.Content,
					})
				}
			}
		}()
	}
}

// 提交任务（不阻塞）
func SubmitMsg(task MsgTask) {
	select {
	case msgQueue <- task:
		// 成功塞进队列
	default:
		// 队列满了，降级：同步写
		database.DB.Create(&model.PrivateMessage{
			FromID:  task.FromID,
			ToID:    task.ToID,
			Content: task.Content,
		})
	}
}
