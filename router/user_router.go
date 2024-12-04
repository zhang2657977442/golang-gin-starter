package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhang2657977442/golang-gin-starter/entity"
	"github.com/zhang2657977442/golang-gin-starter/service"
	"github.com/zhang2657977442/golang-gin-starter/utils"
	"github.com/zhang2657977442/golang-gin-starter/utils/log"
)

type userRouter struct {
	Router  *gin.Engine
	Service *service.Service
}

func init() {
	RegisterHttpRouter("user", exportUserHandler)
}

func exportUserHandler(router *gin.Engine, service *service.Service) {
	u := &userRouter{Router: router, Service: service}
	api := router.Group("/api/user")
	api.Handle(http.MethodPost, "/login", base60HandlersNoLogin(u.login)...)
}

//		@Description	用户登录
//		@Tags			用户账户管理
//		@Summary		用户登录
//	    @Param request body entity.UserLoginReq true "请求参数"
//		@Success		0 {object} entity.UserLoginRps
//		@Failure		500
//		@Router			/user/login [post]
func (ur *userRouter) login(c *gin.Context) {
	req, err := HandlePostRequest(c, &entity.UserLoginReq{})
	if err != nil {
		log.Error("Requset Parsing Error", "User Login", "Login request parameter parsing error")
		utils.NewRsp(c).Fail(entity.PARSE_PARAM_ERROR)
		return
	}

	reqEntity, _ := req.(*entity.UserLoginReq)
	rps, err := ur.Service.UserLogin(reqEntity)
	if err != nil {
		log.Error("Login Processing", "User Login", "User login process error: [%v]", err)
		utils.NewRsp(c).Fail(entity.ERROR_PASSWORD)
		return
	}
	log.Info("Login Success", "User Login", "User login success, userName: [%v]", reqEntity.Username)
	utils.NewRsp(c).Success(rps)
}
