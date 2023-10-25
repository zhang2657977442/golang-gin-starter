package db

import (
	"fmt"

	"github.com/AutoML_Group/omniForce-Backend/entity"
	"github.com/AutoML_Group/omniForce-Backend/utils/log"
)

func (c *DbClient) UploadFile(userId string, req *entity.UploadedFileReq) (*entity.UploadedFileRps, error) {
	insertSql := fmt.Sprintf(`insert into %s set where user_id=?`, TABLE_FILE_INFO)
	if _, err := c.MYSQL.Exec(insertSql, userId); err != nil {
		log.Error("Database Operate Handler", "/api/file/uploadFile", "Database Operate Failed, Error: [%v], SQL: [%v]", err, insertSql)
		return &entity.UploadedFileRps{}, entity.DB_ERROR
	}
	return &entity.UploadedFileRps{}, nil
}
