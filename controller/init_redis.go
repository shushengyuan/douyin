package controller

import "github.com/go-redis/redis"

var RdbLikeVideoId *redis.Client //key:VideoId,value:userId
