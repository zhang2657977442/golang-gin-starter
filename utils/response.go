// Copyright 2019 The JIMDB Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-gin-starter/entity"
	"github.com/golang-gin-starter/utils/log"
)

type Response struct {
	ginContext *gin.Context
	method     string
	err        error
}

func NewRsp(ginContext *gin.Context) *Response {
	return &Response{
		ginContext: ginContext,
	}
}

func (resp *Response) Success(data interface{}) {
	log.Info("Request Result Handler ", resp.ginContext.Request.URL.String(), "Request Success")
	resp.ginContext.JSON(http.StatusOK, &entity.Response{
		Code:    entity.RSPONSE_OK.Code,
		Message: entity.RSPONSE_OK.Msg,
		Time:    time.Now().Unix(),
		Data:    data,
	})
}

func (resp *Response) Fail(err error) {
	log.Error("Request Result Handler", resp.ginContext.Request.URL.String(), "Request Failed, Error: [%v]", err)
	rspErr, _ := err.(*entity.ChatRhinoError)
	resp.ginContext.JSON(http.StatusOK, &entity.Response{
		Code:    rspErr.Code,
		Message: rspErr.Msg,
		Time:    time.Now().Unix(),
		Data:    nil,
	})
}

func (resp *Response) Error() error {
	return resp.err
}
