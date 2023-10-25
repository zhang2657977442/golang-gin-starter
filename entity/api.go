package entity

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    int64       `json:"time"`
	Data    interface{} `json:"data"`
}
