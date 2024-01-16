package service

import (
	"github.com/AutoML_Group/omniForce-Backend/entity"
	"github.com/AutoML_Group/omniForce-Backend/utils"
)

func (s *Service) UserLogin(req *entity.UserLoginReq) (*entity.UserLoginRps, error) {
	// uData, err := s.DbClient.GetUserData(req.Username)
	// if err != nil {
	// 	log.Error("Get Salt and Passowrd", "User Login", "Db query error, username: [%v]", req.Username)
	// 	return nil, err
	// }
	// salt, _ := base64.StdEncoding.DecodeString(uData.Salt)
	// hashedInputPassword := utils.EncryptWithSalt(req.Password, salt)

	// //判断输入密码是否正确
	// if base64.StdEncoding.EncodeToString(hashedInputPassword) != uData.Password {
	// 	//生成token
	// 	log.Error("Compare Password", "User Login", "Password error, username: [%v]", req.Username)
	// 	return nil, entity.ERROR_PASSWORD
	// }

	tokenString, _ := utils.GenToken("userId", 24*7)
	rsp := &entity.UserLoginRps{
		Token: tokenString,
	}
	return rsp, nil
}
