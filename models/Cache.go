package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"time"
)

var redisClient cache.Cache
var enableRedis, _ = beego.AppConfig.Bool("enableRedis")
var redisTime, _ = beego.AppConfig.Int("redisTime")
var YzmClient cache.Cache

func init() {
	if enableRedis {
		config := map[string]string{
			"key":      beego.AppConfig.String("redisKey"),
			"conn":     beego.AppConfig.String("redisConn"),
			"dbNum":    beego.AppConfig.String("redisDbNum"),
			"password": beego.AppConfig.String("redisPwd"),
		}
		bytes, _ := json.Marshal(config)

		redisClient, err = cache.NewCache("redis", string(bytes))
		YzmClient, _ = cache.NewCache("redis", string(bytes))
		if err != nil {
			logs.Error("连接Redis数据库失败")
		} else {
			logs.Info("连接Redis数据库成功")
		}
	}
}

type cacheDb struct{}

var CacheDb = &cacheDb{}

// Set 写入数据的方法
func (c cacheDb) Set(key string, value interface{}) {
	if enableRedis {
		bytes, _ := json.Marshal(value)
		redisClient.Put(key, string(bytes), time.Second*time.Duration(redisTime))
	}
}

// Get 接受数据的方法
func (c cacheDb) Get(key string, obj interface{}) bool {
	if enableRedis {
		if redisStr := redisClient.Get(key); redisStr != nil {
			fmt.Println("在Redis里面读取数据")
			redisValue, ok := redisStr.([]uint8)
			if ok {
				json.Unmarshal([]byte(redisValue), obj)
				return true
			} else {
				fmt.Println("获取Redis数据失败")
			}
		}
	}
	return false
}
