// router package: course api.
package router

import (
	"github.com/labstack/echo/v4"

	"echo-learning/controller"
	"echo-learning/middleware"
)

func DispatchCourseRouters(echoIns *echo.Group) {
	courseGrp := echoIns.Group("/courses")
	// query course list.
	courseGrp.GET("", middleware.WrapResponse(controller.ListCourse))
	// query course by id.
	courseGrp.GET("/:id", middleware.WrapResponse(controller.GetCourse))
	// create course.
	courseGrp.POST("", middleware.WrapResponse(controller.CreateCourse))
	// update course.
	courseGrp.PUT("/:id", middleware.WrapResponse(controller.UpdateCourse))
	// delete course.
	courseGrp.DELETE("/:id", middleware.WrapResponse(controller.DeleteCourse))
}
