package db

import (
	"buycar/pkg/constants"
	"buycar/pkg/errno"
	"buycar/pkg/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("数据库初始化获取DSN失败: %v", err))
	}

	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,  // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
			SkipDefaultTransaction: false, // 不禁用默认事务(即单个创建、更新、删除时使用事务)
			TranslateError:         true,  // 允许翻译错误
			Logger:                 logger.Default.LogMode(logger.Warn),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
		})
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL 连接数据库失败: %v", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL 获取数据库句柄失败: %v", err))
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // 最大闲置连接数
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // 最大连接数
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // 最大可复用时间
	sqlDB.SetConnMaxIdleTime(constants.ConnMaxIdleTime) // 最长保持空闲状态时间

	if err = sqlDB.Ping(); err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("数据库连通性检查失败: %v", err))
	}

	return nil
}
