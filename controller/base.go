package controller

import (
	"github.com/labstack/echo/v4"
)

// ParseRequest bind and validate req.
func ParseRequest(ctx echo.Context, req interface{}) error {
	return ctx.Bind(req)
}
