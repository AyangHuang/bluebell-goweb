package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodePostIDNotExist
	CodeNeedLogin
	CodeInvalidToken
	CodePostExpired
	CodePostZero
	CodePostVoteRepeat
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodePostIDNotExist:  "帖子不存在",
	// 一般不暴露服务器内部错误，对外统一暴露“服务繁忙”
	CodeServerBusy: "服务繁忙",

	CodeNeedLogin:      "需要登录",
	CodeInvalidToken:   "无效的token",
	CodePostExpired:    "帖子已过期，无法投票",
	CodePostZero:       "未投票，无法取消",
	CodePostVoteRepeat: "请无重复投票",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
