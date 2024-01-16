package entity

var (
	RSPONSE_OK             = &ChatRhinoError{0, "操作成功"}
	DB_ERROR               = &ChatRhinoError{1, "数据库操作出错"}
	PARSE_PARAM_ERROR      = &ChatRhinoError{2, "请求解析错误"}
	TOKEN_INVALID_ERROR    = &ChatRhinoError{3, "登录信息无效"}
	ERROR_PASSWORD         = &ChatRhinoError{4, "密码错误"}
	USER_NOT_EXIST         = &ChatRhinoError{5, "账户不存在"}
	HTTP_REQUEST_ERROR     = &ChatRhinoError{6, "Http请求出错"}
	INTERNAL_ERROR         = &ChatRhinoError{7, "服务器繁忙"}
	USERNAME_ALREADY_EXIST = &ChatRhinoError{9, "用户名已存在"}
	EMAIL_ALREADY_EXIST    = &ChatRhinoError{10, "邮件已被中使用"}
	RSPONSE_ERROR          = &ChatRhinoError{11, "操作失败"}
	MISSING_TOKEN          = &ChatRhinoError{401, "缺少token信息"}
)

type ChatRhinoError struct {
	Code int
	Msg  string
}

func (e *ChatRhinoError) Error() string {
	return e.Msg
}
