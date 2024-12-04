package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zhang2657977442/golang-gin-starter/entity"
)

func (s *Service) UploadFile(userId string, req *entity.UploadedFileReq) (*entity.UploadedFileRps, error) {

	if req == nil || req.File == nil {
		return nil, errors.New("invalid file")
	}

	// 读取文件
	src, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	content, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// 生成文件ID
	fileId := md5.Sum(content)
	fileIdStr := hex.EncodeToString(fileId[:])
	fileType := strings.Split(req.File.Filename, ".")[1]

	// 保存文件
	dir := s.Conf.Global.FileDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	filePath := dir + "/" + fileIdStr + "." + fileType
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = dst.Write(content)
	if err != nil {
		return nil, err
	}

	// 进行数据库操作
	// data, err := s.DbClient.UploadFile(userId, req)

	rps := &entity.UploadedFileRps{
		FileId: fileIdStr,
		Path:   filePath,
		Size:   len(content),
	}

	return rps, nil
}
