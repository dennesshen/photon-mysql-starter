package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dennesshen/photon-core-starter/configuration"
	"github.com/dennesshen/photon-core-starter/log"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbAction func(ctx context.Context, db *gorm.DB) (err error)

var customAction = []DbAction{}

func RegisterDbCustomize(action DbAction) {
	customAction = append(customAction, action)
}

func Start(ctx context.Context) (err error) {
	log.Logger().Info(ctx, "init database")
	config, err = configuration.Get[Config](ctx)
	if err != nil {
		log.Logger().Error(ctx, "failed to get database config", "error", err)
		return
	}
	
	if masterDb, err = connect(ctx, config.Database.Master); err != nil {
		log.Logger().Error(ctx, "fail to connect master database", "error", err, "config", config)
		return
	}
	
	if slaveDb, err = connect(ctx, config.Database.Slave); err != nil {
		log.Logger().Error(ctx, "fail to connect slave database", "error", err, "config", config)
		return
	}
	
	for _, action := range customAction {
		if err = action(ctx, masterDb); err != nil {
			log.Logger().Error(ctx, "failed to customize master database", "error", err)
			return
		}
		if err = action(ctx, slaveDb); err != nil {
			log.Logger().Error(ctx, "failed to customize slave database", "error", err)
			return
		}
	}
	return
}

// 連線資料庫
func connect(ctx context.Context, connectData ConnectData) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		connectData.Username,
		connectData.Password,
		connectData.Host,
		connectData.Port,
		config.Database.Database)
	if db, err = gorm.Open(mysql.Open(dsn)); err != nil {
		log.Logger().Error(ctx, "failed to connect database", "error", err, "dsn", dsn)
		return
	}
	
	err = setConnectPool(ctx, db)
	return
}

// 設定 connection pool
func setConnectPool(ctx context.Context, db *gorm.DB) (err error) {
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		log.Logger().Error(ctx, "failed to use slave database with opentelemetery", "error", err)
		return
	}
	
	maxIdleConns := config.Database.Connection.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	
	maxOpenConns := config.Database.Connection.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 50
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	
	maxLifetime := config.Database.Connection.MaxLifetimeSecond
	if maxLifetime == 0 {
		maxLifetime = 600
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))
	return
}
