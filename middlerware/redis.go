package middlerware

import (
	"github.com/go-redis/redis"

	"../config"
)

var rc *redis.Client

func init() {
	dsn := config.AppConfig.Redis.Host + ":" + config.AppConfig.Redis.Port
	rc = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: config.AppConfig.Redis.Password,
	})
	Cont.Register("redis", rc)
}

func GetRedis(r interface{}) *redis.Client {
	rc, ok := r.(*redis.Client)
	if !ok {
		return nil
	}
	return rc
}
