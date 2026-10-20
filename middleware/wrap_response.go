// middleware package
package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"echo-learning/model"
)

// WrapResponse 包装 controller 的返回值，统一按 Response 结构输出 JSON。
// controller 的签名固定为 (ctx echo.Context) (interface{}, error)：
// 返回 error（如参数绑定失败）也按统一格式输出，不中断请求。
func WrapResponse(handler func(echo.Context) (interface{}, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		result, err := handler(ctx)
		if err != nil {
			return ctx.JSON(http.StatusOK, model.Response{
				Code: model.CodeFailure,
				Msg:  err.Error(),
				Data: nil,
			})
		}
		return ctx.JSON(http.StatusOK, result)
	}
}
