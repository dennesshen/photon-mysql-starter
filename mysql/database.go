package mysql

import (
	"context"

	"gorm.io/gorm"
)

var masterDb, slaveDb *gorm.DB

// 取得 master database
func Master(ctx context.Context) *gorm.DB {
	return masterDb.WithContext(ctx)
}

// 取得 slave database
func Slave(ctx context.Context) *gorm.DB {
	return slaveDb.WithContext(ctx)
}
