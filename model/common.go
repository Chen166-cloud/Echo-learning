// model package: common request/response define.
package model

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 响应码
const (
	CodeSuccess = 0
	CodeFailure = 1
)

// 响应消息
const (
	MsgSuccess = "success"
)
