package loader

import (
	"github.com/pkg/errors"

	"go-pulgin/internal/datasource/db"
	"go-pulgin/internal/infra/loader"
	"go-pulgin/internal/infra/plugin"
	"go-pulgin/logger"
)

const (
	redisKeyTitle        = "key"
	redisDbNumTitle      = "db_num"
	redisDbNumInvalidSig = -1
	redisDataTypeTitle   = "data_type"
	redisDataTypeString  = "strings"
	redisDataTypeHash    = "hashes"
)

type RedisLoader struct {
	key      string
	dataType string
	dbNum    int
}

func (rl *RedisLoader) Install() error {
	if rl.key == "" || rl.dbNum < 0 {
		return errors.Errorf("[RedisLoader] context invalid. key: %v, dbNum: %v, dataType: %v",
			rl.key, rl.dbNum, rl.dataType)
	}
	return nil
}

func (rl *RedisLoader) Uninstall() {
}

func (rl *RedisLoader) SetContext(ctx plugin.Context) {
	if key, ok := ctx.GetString(redisKeyTitle); ok {
		rl.key = key
	}
	dbNum := ctx.GetIntOrDefault(redisDbNumTitle, redisDbNumInvalidSig)
	rl.dbNum = dbNum
	dataType := ctx.GetStringOrDefault(redisDataTypeTitle, redisDataTypeString)
	rl.dataType = dataType
}

func (rl *RedisLoader) Load(event *plugin.Event) error {
	payload := event.Payload()
	switch rl.dataType {
	case redisDataTypeString:
		if err := db.SetRedisValue(rl.key, payload, rl.dbNum); err != nil {
			logger.Errorf("[RedisLoader] SetRedisValue key: %v, error: %v", rl.key, err)
			return err
		}
	case redisDataTypeHash:
		valMap, ok := payload.(map[string]interface{})
		if !ok {
			logger.Errorf("[RedisLoader] data type is hash but payload is not a map")
			return errors.New("data type is hash but payload is not a map")
		}
		for field, val := range valMap {
			if err := db.SetRedisHashValue(rl.key, field, val, rl.dbNum); err != nil {
				logger.Errorf("[RedisLoader] SetRedisHashValue key: %v, field: %v, error: %v", rl.key, field, err)
				continue
			}
		}
	default:
		return errors.Errorf("[RedisLoader] no supported redis data type: %v", rl.dataType)
	}
	return nil
}

func init() {
	loader.Add("redis_loader", func() loader.Plugin {
		return &RedisLoader{}
	})
}
