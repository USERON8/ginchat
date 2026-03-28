// pkg/response/code.go
package response

type Code int

const (
	// 通用
	CodeSuccess     Code = 0
	CodeParamError  Code = 1001
	CodeServerError Code = 1002

	// 用户
	CodeUserNotFound  Code = 2001
	CodeUserExist     Code = 2002
	CodePasswordError Code = 2003
	CodeTokenInvalid  Code = 2004
	CodeUnauthorized  Code = 2005

	// 好友
	CodeFriendExist        Code = 3001
	CodeFriendNotFound     Code = 3002
	CodeFriendNoPermission Code = 3003

	// 消息
	CodeMsgSendFail Code = 4001
	// 用户
	CodeRegisterFail Code = 2006
	CodeUpdateFail   Code = 2007

	// 好友
	CodeCannotAddSelf Code = 3004
	// 补充群组错误码
	CodeGroupNotFound     Code = 5001
	CodeGroupNoPermission Code = 5002
	CodeGroupMemberExist  Code = 5003
	CodeRateLimitExceeded Code = 1003
)

// 错误码对应的默认消息
var msgMap = map[Code]string{
	CodeSuccess:            "success",
	CodeParamError:         "参数错误",
	CodeServerError:        "服务器错误",
	CodeUserNotFound:       "用户不存在",
	CodeUserExist:          "用户名已存在",
	CodePasswordError:      "密码错误",
	CodeTokenInvalid:       "token无效",
	CodeUnauthorized:       "未登录",
	CodeFriendExist:        "已申请或已是好友",
	CodeFriendNotFound:     "申请不存在",
	CodeFriendNoPermission: "无权操作",
	CodeMsgSendFail:        "消息发送失败",
	CodeRegisterFail:       "注册失败",
	CodeUpdateFail:         "更新失败",
	CodeCannotAddSelf:      "不能添加自己",
	// msgMap 补充
	CodeGroupNotFound:     "群组不存在",
	CodeGroupNoPermission: "无权操作",
	CodeGroupMemberExist:  "用户已在群中",
	CodeRateLimitExceeded: "请求过于频繁，请稍后再试",
}

func (c Code) Msg() string {
	if msg, ok := msgMap[c]; ok {
		return msg
	}
	return "未知错误"
}
