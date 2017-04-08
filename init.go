package smservice

import (
	"database/sql"
)

func init() {

	// 加载配置文件
	LoadConfig()

	// 加载redis
	redisClient = initRedis()

	// 加载mysql
	dataSource := config.MysqlConf["user"] + ":" + config.MysqlConf["password"] + "@tcp(" + config.MysqlConf["host"] + ":" + config.MysqlConf["port"] + ")/" + config.MysqlConf["name"] + "?charset=utf8"

	db, err = sql.Open(DBTYPE, dataSource)

	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}
