package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AutoML_Group/omniForce-Backend/constants"
	"github.com/AutoML_Group/omniForce-Backend/router"
	"github.com/AutoML_Group/omniForce-Backend/service"
	"github.com/AutoML_Group/omniForce-Backend/utils/log"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	confPath string
)

func init() {
	//-config config.toml
	flag.StringVar(&confPath, "conf", getDefaultConfigFile(), "chatrhino config path")
}

func main() {
	flag.Parse()

	if err := constants.InitConfig(confPath); err != nil {
		log.Error("Init Config", "Init", "config Init error, erroe: [%v], err")
		panic(err)
	}

	if err := initLogger(constants.Conf().Log); err != nil {
		log.Error("Init Logger", "Init", "Logger Init error, erroe: [%v], err")
		panic(err)
	}
	defer log.Flush()

	if err := service.InitService(); err != nil {
		log.Error("Init Service", "Init", "Service Init error, error [%v]", err)
		panic(err)
	}

	startHttpServer()
}

func getGinLogger() gin.HandlerFunc {
	c := constants.Conf().Log
	writer := &lumberjack.Logger{
		Filename:   filepath.Join(c.Path, "gin.log"),
		MaxAge:     c.MaxRetainDays,
		MaxSize:    c.MaxFileSize,
		MaxBackups: c.MaxRetainCount,
		LocalTime:  true,
		Compress:   false,
	}

	return gin.LoggerWithWriter(writer)
}

// @title					NextGPT接口文档
// @version					1.0
// @description				NextGPT Backend Server
// @BasePath					/api
//
// @securityDefinitions.basic	BasicAuth
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for what is this security definition being used
func startHttpServer() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(getGinLogger(), gin.Recovery())

	service := service.NewService()
	if service == nil {
		panic("service is not ready!")
	}

	//register handler
	for _, handler := range router.HttpRouters {
		handler(engine, service)
	}

	if err := engine.Run(":" + fmt.Sprint(constants.Conf().Global.Port)); err != nil {
		log.Error("Strat Backend", "Init", "exit, err: %v, ", err)
		panic(err)
	}

}

func initLogger(c *constants.LogConfig) error {
	zapCfg := &log.LogrusConfig{
		Filename:       filepath.Join(c.Path, fmt.Sprintf("%s.log", filepath.Base(os.Args[0]))),
		Level:          c.Level,
		MaxFileSizeMB:  c.MaxFileSize,
		MaxRetainDays:  c.MaxRetainDays,
		MaxRetainFiles: c.MaxRetainCount,
	}

	logger, err := log.NewLogrusLogger(zapCfg)
	if err != nil {
		return err
	}

	log.SetDefaultLogger(logger)
	return nil
}

func getDefaultConfigFile() (defaultConfigFile string) {
	if currentExePath, err := getCurrentPath(); err == nil {
		path := currentExePath + "config.toml"
		if ok, err := pathExists(path); ok {
			return path
		} else if err != nil {
			log.Error("Strat Backend", "Init", "check path:%s err : %s", path, err.Error())
		}
	}

	if sourceCodeFileName, err := getCurrentSourceCodePath(); nil == err {
		lastIndex := strings.LastIndex(sourceCodeFileName, "/")
		path := sourceCodeFileName[0:lastIndex+1] + "config.toml"
		if ok, err := pathExists(path); ok {
			return path
		} else if err != nil {
			log.Error("Strat Backend", "Init", "check path:%s err : %s", path, err.Error())
		}
	}
	return
}

func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\"`)
	}
	return path[0 : i+1], nil
}

func getCurrentSourceCodePath() (fileName string, err error) {
	_, fileName, _, ok := callerName(2)
	if !ok {
		err = errors.New("error: Can't get the current source code path")
	}
	return fileName, err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func callerName(skip int) (funcName, file string, line int, ok bool) {
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip); !ok {
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	return
}
