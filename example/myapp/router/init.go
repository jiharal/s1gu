package router

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/KancioDevelopment/lib-angindai/logging"
	"github.com/asaskevich/govalidator"
	"github.com/gomodule/redigo/redis"
	"github.com/jiharal/s1gu/example/myapp/api"
	"github.com/pkg/errors"
)

type (
	InitOption struct{}
)

var (
	logger    *logging.Logger
	dbPool    *sql.DB
	cachePool *redis.Pool
	cfg       InitOption

	userService *api.UserModule
)

func Init(lg *logging.Logger, db *sql.DB, cache *redis.Pool, opt InitOption) {
	logger = lg
	dbPool = db
	cachePool = cache
	cfg = opt

	userService = api.NewUserModule(dbPool, cachePool)
}

// ParseBodyData parse json-formatted request body into given struct.
func ParseBodyData(ctx context.Context, r *http.Request, data interface{}) error {
	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "read")
	}

	err = json.Unmarshal(bBody, data)
	if err != nil {
		return errors.Wrap(err, "json")
	}

	valid, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errors.Wrap(err, "validate")
	}

	if !valid {
		return errors.New("invalid data")
	}

	return nil
}
