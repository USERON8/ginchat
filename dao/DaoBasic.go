package dao

import (
	"context"
	"ginchat/global"
)

// BaseDao 通用数据访问层
type BaseDao[T any] struct{}

func NewBaseDao[T any]() *BaseDao[T] {
	return &BaseDao[T]{}
}

// Create 创建
func (d *BaseDao[T]) Create(ctx context.Context, data *T) error {
	return global.DB.WithContext(ctx).Create(data).Error
}

// GetByID 根据ID查询
func (d *BaseDao[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var data T
	err := global.DB.WithContext(ctx).First(&data, id).Error
	return &data, err
}

// List 列表查询
func (d *BaseDao[T]) List(ctx context.Context, where interface{}, args ...interface{}) ([]*T, error) {
	var list []*T
	err := global.DB.WithContext(ctx).Where(where, args...).Find(&list).Error
	return list, err
}
