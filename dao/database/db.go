package db

import (
	"database/sql"
	"fmt"

	"github.com/AutoML_Group/omniForce-Backend/dao/database/client"
	"github.com/AutoML_Group/omniForce-Backend/utils/log"
)

const (
	TABLE_FILE_INFO = "file_info"
)

var singleDbClient *DbClient

type DbClient struct {
	MYSQL *sql.DB
}

func Db() *DbClient {
	return singleDbClient
}

func NewDbClient() error {

	mysqlClient, err := client.NewMysqlClient()
	if err != nil {
		log.Error("Client Init", "DB Client", "Init mysql client error, error: [%v]", err)
		return fmt.Errorf("Init mysql client error")
	}
	log.Info("Client Init", "DB Client", "Init mysql client success")

	singleDbClient = &DbClient{
		MYSQL: mysqlClient,
	}
	return nil
}
