
	package router

	import (
		"database/sql"
		"encoding/json"
		"net/http"

		"github.com/AuthScureDevelopment/lib-arjuna/logging"
		"github.com/gomodule/redigo/redis"
	)

	var (
		Logger    *logging.Logger
		DbPool    *sql.DB
		CachePool *redis.Pool
	)

	func Init(db *sql.DB, cachePool *redis.Pool, logger *logging.Logger) {
		DbPool = db
		CachePool = cachePool
		Logger = logger
	}

	type ErrorMethod struct {
		Errors interface{} `json:"errors"`
	}

	type ResponseMethod struct {
		Data interface{} `json:"data"`
	}

	func respondWithError(w http.ResponseWriter, code int, message interface{}) {
		response, _ := json.Marshal(ErrorMethod{Errors: message})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}

	func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
		response, _ := json.Marshal(ResponseMethod{Data: payload})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
	