package router

import (
	"net/http"

	"github.com/AutoML_Group/omniForce-Backend/entity"
	"github.com/AutoML_Group/omniForce-Backend/service"
	"github.com/AutoML_Group/omniForce-Backend/utils"
	"github.com/gin-gonic/gin"
)

type fileRouter struct {
	Router  *gin.Engine
	Service *service.Service
}

func init() {
	RegisterHttpRouter("file", exportFileHandler)
}

func exportFileHandler(router *gin.Engine, service *service.Service) {
	u := &fileRouter{Router: router, Service: service}
	api := router.Group("/api/file")
	api.Handle(http.MethodPost, "/uploadFile", base60Handlers(u.uploadFile)...)

}

// @Description	本地文件上传
// @Tags			文件管理
// @Accept mpfd
// @Produce json
// @Summary		本地文件上传
// @Param file formData file true "File"
// @Success		0 {object} entity.UploadedFileRps
// @Failure		500
// @Security		ApiKeyAuth
// @Router			/file/uploadFile [post]
func (ur *fileRouter) uploadFile(c *gin.Context) {
	userId := c.GetString(UserId)
	file, err := c.FormFile("file")
	if err != nil {
		utils.NewRsp(c).Fail(entity.RSPONSE_ERROR)
		return
	}

	reqEntity := entity.UploadedFileReq{File: file}
	rps, err := ur.Service.UploadFile(userId, &reqEntity)
	if err != nil {
		utils.NewRsp(c).Fail(entity.RSPONSE_ERROR)
		return
	}
	utils.NewRsp(c).Success(rps)
}
