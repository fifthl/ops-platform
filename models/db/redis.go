/*
 * @Author: nevin
 * @Date: 2023-11-13 11:26:14
 * @LastEditTime: 2023-11-13 15:42:58
 * @LastEditors: nevin
 * @Description: redis
 */
package db

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var RedisDb *redis.Client // 声明一个全局的redisDb变量

// 根据redis配置初始化一个客户端
func initClient() (err error) {
	dbnum, err := beego.AppConfig.Int("redis_db")
	if err != nil {
		panic(err)
	}

	RedisDb = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.String("redis_addr"),     // redis地址
		Password: beego.AppConfig.String("redis_password"), // redis密码，没有则留空
		DB:       dbnum,                                    // 默认数据库，默认是0
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err = RedisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	err := initClient()
	if err != nil {
		//redis连接错误
		panic(err)
	}
	log.Println("Redis连接成功")
}

// TODO: 设置
func SetEX(key string, values interface{}, endtime time.Duration) bool {
	err := RedisDb.Set(key, values, endtime).Err()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// TODO: 获取
func GetEX() {
	val2, err := RedisDb.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func GetRedis(key string) (value string, err error) {
	result := RedisDb.Get(key)
	return result.Result()
}
