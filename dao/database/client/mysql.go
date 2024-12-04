package client

import (
	"database/sql"
	"fmt"

	"github.com/golang-gin-starter/constants"
	"github.com/golang-gin-starter/utils/log"
)

func NewMysqlClient() (*sql.DB, error) {

	// 数据库配置信息
	config := &constants.MysqlConfig{}

	mysqlAddress := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	mysqlClient, err := sql.Open("mysql", mysqlAddress)
	if err != nil {
		log.Info("Client Init", "Mysql", "Connecting to database error, error: [%v], mysql config", err, mysqlAddress)
		return nil, err
	}
	mysqlClient.SetMaxOpenConns(20)
	mysqlClient.SetMaxIdleConns(0)
	return mysqlClient, nil
}
