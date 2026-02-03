package infrastructure

import (
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
)

func CreateRedis() *redis.Client {
	redis, err := redisPkg.NewClient("")
	common.HandleError(err)

	return redis
}
