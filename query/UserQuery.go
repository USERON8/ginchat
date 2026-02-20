package query

type UserQuery[T any] interface {
	// GetByID SELECT * FROM user_basic WHERE id=@id
	GetByID(id uint) (T, error)

	// GetByUsername SELECT * FROM user_basic WHERE username=@username
	GetByUsername(username string) (T, error)

	// GetAllOnlineUsers SELECT * FROM user_basic WHERE is_logged_in = 1
	GetAllOnlineUsers() ([]T, error)
}
