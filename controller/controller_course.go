// course api controller
package controller

import (
	"github.com/labstack/echo/v4"

	"echo-learning/logic/course_logic"
	"echo-learning/model"
)

// ListCourse query all courses.
func ListCourse(ctx echo.Context) (interface{}, error) {
	resp, err := course_logic.ListCourse(ctx)
	if err != nil {
		return model.Response{
			Code: model.CodeFailure,
			Msg:  err.Error(),
			Data: nil,
		}, nil
	}

	return model.Response{
		Code: model.CodeSuccess,
		Msg:  model.MsgSuccess,
		Data: resp,
	}, nil
}

// GetCourse query one course by id.
func GetCourse(ctx echo.Context) (interface{}, error) {
	req := &model.GetCourseReq{}
	if err := ParseRequest(ctx, req); err != nil {
		return nil, err
	}

	resp, err := course_logic.GetCourse(ctx, req)
	if err != nil {
		return model.Response{
			Code: model.CodeFailure,
			Msg:  err.Error(),
			Data: nil,
		}, nil
	}

	return model.Response{
		Code: model.CodeSuccess,
		Msg:  model.MsgSuccess,
		Data: resp,
	}, nil
}

// CreateCourse create one course.
func CreateCourse(ctx echo.Context) (interface{}, error) {
	req := &model.CreateCourseReq{}
	if err := ParseRequest(ctx, req); err != nil {
		return nil, err
	}

	resp, err := course_logic.CreateCourse(ctx, req)
	if err != nil {
		return model.Response{
			Code: model.CodeFailure,
			Msg:  err.Error(),
			Data: nil,
		}, nil
	}

	return model.Response{
		Code: model.CodeSuccess,
		Msg:  model.MsgSuccess,
		Data: resp,
	}, nil
}

// UpdateCourse update one course by id.
func UpdateCourse(ctx echo.Context) (interface{}, error) {
	req := &model.UpdateCourseReq{}
	if err := ParseRequest(ctx, req); err != nil {
		return nil, err
	}

	resp, err := course_logic.UpdateCourse(ctx, req)
	if err != nil {
		return model.Response{
			Code: model.CodeFailure,
			Msg:  err.Error(),
			Data: nil,
		}, nil
	}

	return model.Response{
		Code: model.CodeSuccess,
		Msg:  model.MsgSuccess,
		Data: resp,
	}, nil
}

// DeleteCourse delete one course by id.
func DeleteCourse(ctx echo.Context) (interface{}, error) {
	req := &model.DeleteCourseReq{}
	if err := ParseRequest(ctx, req); err != nil {
		return nil, err
	}

	err := course_logic.DeleteCourse(ctx, req)
	if err != nil {
		return model.Response{
			Code: model.CodeFailure,
			Msg:  err.Error(),
			Data: nil,
		}, nil
	}

	return model.Response{
		Code: model.CodeSuccess,
		Msg:  model.MsgSuccess,
		Data: nil,
	}, nil
}
