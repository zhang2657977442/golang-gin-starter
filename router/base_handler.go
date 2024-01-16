package router

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	_ "github.com/AutoML_Group/omniForce-Backend/docs"
	"github.com/AutoML_Group/omniForce-Backend/entity"
	"github.com/AutoML_Group/omniForce-Backend/service"
	"github.com/AutoML_Group/omniForce-Backend/utils"
	"github.com/AutoML_Group/omniForce-Backend/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	Timeout    = "timeout"
	CancelFunc = "__cancelFunc"
	Ctx        = "__ctx"
	Start      = "__start_time"
	UserId     = "__user_id"
)

// default 60s to timout
var base60 = newBaseHandler(60)

type baseHandler struct {
	timeout int64 //default timeout Second
}

func newBaseHandler(timeout int64) *baseHandler {
	return &baseHandler{timeout: timeout}
}

func HandlePostRequest(c *gin.Context, reqEntity interface{}) (interface{}, error) {
	req := reqEntity
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("Request Param Handler", c.Request.URL.String(), "Parse Post Body Error")
		return nil, entity.PARSE_PARAM_ERROR
	}

	err = json.Unmarshal(bytes, req)
	if err != nil {
		log.Error("Request Param Handler", c.Request.URL.String(), "Unmarshal Param Error")
		return nil, entity.PARSE_PARAM_ERROR
	}

	return req, nil
}

func (b *baseHandler) Timeout() int64 {
	return b.timeout
}

func (b *baseHandler) LoginHandler(c *gin.Context) {

	//Token放在Header的Authorization中
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    entity.MISSING_TOKEN.Code,
			"message": "缺少token信息",
		})
		c.Abort()
		return
	}

	// 我们使用之前定义好的解析JWT的函数来解析它
	mc, err := utils.ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    entity.TOKEN_INVALID_ERROR.Code,
			"message": "无效的Token",
		})
		c.Abort()
		return
	}
	// 将当前请求的userid信息保存到请求的上下文c上
	c.Set(UserId, mc.UserId)
}

func (b *baseHandler) PaincHandler(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var msg string
				switch r.(type) {
				case error:
					msg = r.(error).Error()
				default:
					if str, err := cast.ToStringE(r); err != nil {
						msg = " Server internal error "
					} else {
						msg = str
					}
				}
				log.Error("Base Handler", "Panic", msg)
				buf := make([]byte, 10240)
				runtime.Stack(buf, false)
				log.Error("Base Handler", "Panic", string(buf))
				c.JSON(http.StatusInternalServerError, map[string]string{"message": msg})
			}
		}()

		handlerFunc(c)
	}
}

func (b *baseHandler) TimeOutHandler(c *gin.Context) {
	param := c.Query(Timeout)

	// add start time for monitoring
	c.Set(Start, time.Now())

	ctx := context.Background()

	if param != "" {
		if v, err := cast.ToInt64E(param); err != nil {
			log.Error("Base Handler", "Panic", "parse timeout err , it must int value")
		} else {
			ctx, cancelFunc := context.WithTimeout(ctx, time.Duration(v*int64(time.Second)))
			c.Set(Ctx, ctx)
			c.Set(CancelFunc, cancelFunc)
			return
		}
	}

	if b.timeout > 0 {
		ctx, cancelFunc := context.WithTimeout(ctx, time.Duration(b.timeout*int64(time.Second)))
		c.Set(Ctx, ctx)
		c.Set(CancelFunc, cancelFunc)
		return
	}

	c.Set(Ctx, ctx)
}

func (b *baseHandler) TimeOutEndHandler(c *gin.Context) {
	if value, exists := c.Get(CancelFunc); exists && value != nil {
		value.(context.CancelFunc)()
	}
}

func (b baseHandler) GetHandlers(businessFunc gin.HandlerFunc) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		b.TimeOutHandler,
		b.LoginHandler,
		b.PaincHandler(businessFunc),
		b.TimeOutEndHandler,
	}
}

func (b baseHandler) GetHandlersNoLogin(businessFunc gin.HandlerFunc) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		b.TimeOutHandler,
		b.PaincHandler(businessFunc),
		b.TimeOutEndHandler,
	}
}

func base60Handlers(businessFunc gin.HandlerFunc) []gin.HandlerFunc {
	return base60.GetHandlers(businessFunc)
}

func base60HandlersNoLogin(businessFunc gin.HandlerFunc) []gin.HandlerFunc {
	return base60.GetHandlersNoLogin(businessFunc)
}

var HttpRouters = make(map[string]func(router *gin.Engine, service *service.Service), 9)

func RegisterHttpRouter(name string, handlerFunc func(router *gin.Engine, service *service.Service)) {
	HttpRouters[name] = handlerFunc
}

type baseRouter struct {
	Router  *gin.Engine
	Service *service.Service
}

func init() {
	RegisterHttpRouter("base", exportBaseHandler)
}

func initSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func exportBaseHandler(router *gin.Engine, service *service.Service) {
	initSwagger(router)
	b := &baseRouter{Router: router, Service: service}
	router.Handle(http.MethodGet, "/", base60Handlers(b.helloWorld)...)

}

func (br *baseRouter) helloWorld(c *gin.Context) {
	log.Info("OmniForce Test", "Base Router", "This is OmniForce Backend")
}
