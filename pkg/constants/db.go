package constants

import "time"

const (
	MaxConnections  = 1000             // (DB) 最大连接数
	MaxIdleConns    = 10               // (DB) 最大空闲连接数
	ConnMaxLifetime = 10 * time.Second // (DB) 最大可复用时间
	ConnMaxIdleTime = 5 * time.Minute  // (DB) 最长保持空闲状态时间
)

const (
	UserTableName = "users"
	ScoreTable    = "score_transactions"
	GiftTable     = "gifts"
)
