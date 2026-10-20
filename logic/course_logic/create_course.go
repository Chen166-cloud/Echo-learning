// course_logic package
package course_logic

import (
	"errors"

	"github.com/labstack/echo/v4"

	"echo-learning/dao"
	"echo-learning/model"
)

// CreateCourse 新增课程
func CreateCourse(ctx echo.Context, req *model.CreateCourseReq) (*dao.Course, error) {
	if req.Title == "" {
		return nil, errors.New("title is empty")
	}

	course := &dao.Course{
		Title:       req.Title,
		Teacher:     req.Teacher,
		Description: req.Description,
		Price:       req.Price,
	}
	if err := dao.CourseTable.Insert(ctx, course); err != nil {
		ctx.Logger().Warn("insert course failure, err=", err)
		return nil, err
	}
	return course, nil
}
