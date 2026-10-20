// course_logic package
package course_logic

import (
	"github.com/labstack/echo/v4"

	"echo-learning/dao"
	"echo-learning/model"
)

// GetCourse 按 ID 查询课程详情
func GetCourse(ctx echo.Context, req *model.GetCourseReq) (*dao.Course, error) {
	course, err := dao.CourseTable.QueryByID(ctx, req.ID)
	if err != nil {
		ctx.Logger().Warn("query course failure, err=", err)
		return nil, err
	}
	if course.ID == 0 {
		return nil, errRecordNotFound(req.ID)
	}
	return course, nil
}
