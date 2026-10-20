// router package
package router

import (
	"github.com/labstack/echo/v4"
)

// main router.
func DispatchRouters(echoIns *echo.Echo) {
	rRoot := echoIns.Group("/api")

	// course api.
	DispatchCourseRouters(rRoot)
}
