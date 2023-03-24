package controller

import (
	"context"

	"github.com/go-redis/redis"
)

var RdbLikeVideoId *redis.Client //key:VideoId,value:userId
var Ctx = context.Background()

var RdbLikeUserId *redis.Client //key:userId,value:VideoId

func InitRedis() {
	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       5, //  选择将点赞视频id信息存入 DB5.
	})
}
