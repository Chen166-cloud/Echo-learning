// dao package: gorm session.
package dao

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"echo-learning/services/mysql"
)

// GetSession 获取 gorm 会话（带上请求 context，超时可随请求取消）
func GetSession(ctx echo.Context) *gorm.DB {
	return mysql.GetDB().WithContext(ctx.Request().Context())
}
