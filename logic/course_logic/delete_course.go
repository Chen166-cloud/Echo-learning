// course_logic package
package course_logic

import (
	"github.com/labstack/echo/v4"

	"echo-learning/dao"
	"echo-learning/model"
)

// DeleteCourse 按 ID 删除课程
func DeleteCourse(ctx echo.Context, req *model.DeleteCourseReq) error {
	// 先确认记录存在
	course, err := dao.CourseTable.QueryByID(ctx, req.ID)
	if err != nil {
		ctx.Logger().Warn("query course failure, err=", err)
		return err
	}
	if course.ID == 0 {
		return errRecordNotFound(req.ID)
	}

	if err := dao.CourseTable.Delete(ctx, req.ID); err != nil {
		ctx.Logger().Warn("delete course failure, err=", err)
		return err
	}
	return nil
}
