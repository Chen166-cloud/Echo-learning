// course_logic package
package course_logic

import (
	"github.com/labstack/echo/v4"

	"echo-learning/dao"
)

// ListCourse 查询全部课程
func ListCourse(ctx echo.Context) ([]*dao.Course, error) {
	courses, err := dao.CourseTable.QueryAll(ctx)
	if err != nil {
		ctx.Logger().Warn("query course list failure, err=", err)
		return nil, err
	}
	return courses, nil
}
