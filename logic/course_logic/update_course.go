// course_logic package
package course_logic

import (
	"errors"

	"github.com/labstack/echo/v4"

	"echo-learning/dao"
	"echo-learning/model"
)

// UpdateCourse 按 ID 更新课程
func UpdateCourse(ctx echo.Context, req *model.UpdateCourseReq) (*dao.Course, error) {
	if req.Title == "" {
		return nil, errors.New("title is empty")
	}

	// 先查原记录，不存在直接报错
	course, err := dao.CourseTable.QueryByID(ctx, req.ID)
	if err != nil {
		ctx.Logger().Warn("query course failure, err=", err)
		return nil, err
	}
	if course.ID == 0 {
		return nil, errRecordNotFound(req.ID)
	}

	// 用请求字段覆盖后更新
	course.Title = req.Title
	course.Teacher = req.Teacher
	course.Description = req.Description
	course.Price = req.Price
	if err := dao.CourseTable.Update(ctx, course); err != nil {
		ctx.Logger().Warn("update course failure, err=", err)
		return nil, err
	}
	return course, nil
}
