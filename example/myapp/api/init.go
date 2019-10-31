package api

import (
	"database/sql"

	"github.com/KancioDevelopment/lib-angindai/logging"
	"github.com/gomodule/redigo/redis"
)

type (
	InitOption struct {
		SessionExpire string
	}
)

var (
	logger    *logging.Logger
	dbPool    *sql.DB
	cachePool *redis.Pool
	cfg       InitOption
)

func Init(lg *logging.Logger, db *sql.DB, cache *redis.Pool, opt InitOption) {
	logger = lg
	dbPool = db
	cachePool = cache
	cfg = opt
}
