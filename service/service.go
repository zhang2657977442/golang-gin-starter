package service

import (
	db "github.com/zhang2657977442/golang-gin-starter/dao/database"
	// "github.com/zhang2657977442/golang-gin-starter/dao/sid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhang2657977442/golang-gin-starter/constants"
	"github.com/zhang2657977442/golang-gin-starter/utils/log"
)

var serviceInstance *Service = nil

type Service struct {
	Conf     *constants.Config
	DbClient *db.DbClient
	// SidClient *sid.SidClient
}

func NewService() *Service {
	if serviceInstance == nil {
		log.Error("Service Init", "Init", "Firstly invoke initService before get service instance.")
		return nil
	}
	return serviceInstance
}

func InitService() error {
	if err := db.NewDbClient(); err != nil {
		log.Error("Db Init", "Service Init", "Init db client error, error: [%v]", err)
		return err
	}
	serviceInstance = &Service{
		Conf:     constants.Conf(),
		DbClient: db.Db(),
		// SidClient: &sid.SidClient{},
	}
	return nil
}
