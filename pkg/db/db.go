package db

import (
	"gim/config"
	"gim/pkg/logger"
	"gim/pkg/util"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB        *gorm.DB
	RedisCli  *redis.Client
	RedisUtil *util.RedisUtil
)

func init() {
	InitMysql(config.Config.MySQL)
	InitRedis(config.Config.RedisHost, config.Config.RedisPassword)
}

// InitMysql 初始化MySQL
func InitMysql(dataSource string) {
	logger.Logger.Info("init mysql")
	var err error
	DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	DB.SingularTable(true)
	DB.LogMode(true)
	logger.Logger.Info("init mysql ok")
}

// InitRedis 初始化Redis
func InitRedis(addr, password string) {
	logger.Logger.Info("init redis")
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		panic(err)
	}

	RedisUtil = util.NewRedisUtil(RedisCli)
	logger.Logger.Info("init redis ok")
}
