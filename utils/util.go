package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zhang2657977442/golang-gin-starter/entity"
	"github.com/zhang2657977442/golang-gin-starter/utils/log"
)

var MySecret = []byte("ADMINCHATRHINO") // TODO 需要替换 jwt 加密的key

type MyClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func IsValidStruct(data map[string]bool, obj interface{}) bool {
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		return false
	}

	for key := range data {
		if err := validate.Var(data[key], fmt.Sprintf("required,eqfield=%s", key)); err != nil {
			return false
		}
	}
	return true
}

func GenerateSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	_, _ = rand.Read(salt)
	return salt
}

func EncryptWithSalt(password string, salt []byte) []byte {
	// 将密码与盐拼接在一起
	passwordWithSalt := append([]byte(password), salt...)

	// 使用哈希算法进行加密
	h := sha256.New()
	h.Write(passwordWithSalt)
	hashedPassword := h.Sum(nil)

	return hashedPassword
}

func FormatFileSize(size interface{}) (result string) {
	size64, _ := size.(float32)
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := 0

	for size64 >= 1024 && i+1 < len(units) {
		size64 /= 1024
		i++
	}

	result = fmt.Sprintf("%.2f %s", size64, units[i])
	return result
}

func GenerateUUIDs(count int) []string {
	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		uuids[i] = uuid.New().String()
	}
	return uuids
}

func TimeStampToUnixTime(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
}

func GenToken(userid string, hour int) (string, error) {
	c := MyClaims{
		userid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Issuer:    "admin-chatRhino-Backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, entity.TOKEN_INVALID_ERROR
}

func GetKeyByValue(m map[string]int, value int) (string, bool) {
	reverseMap := make(map[int]string)
	for key, val := range m {
		reverseMap[val] = key
	}

	key, found := reverseMap[value]
	if !found {
		log.Error("Get Dict Key Handler", "", "Assert Type Failed")
		return "", false
	}
	return key, found
}

func DictAssertType(value interface{}) (map[string]int, error) {
	resultType, ok := value.(map[string]int)
	if !ok {
		err := errors.New("Assert Type Failed")
		log.Error("Assert Type Handler", "", "Assert Type Failed")
		return nil, err
	}
	return resultType, nil
}
