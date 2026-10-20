// model package: course request define.
package model

// CreateCourseReq 新增课程请求体
type CreateCourseReq struct {
	Title       string  `json:"title"`
	Teacher     string  `json:"teacher"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// UpdateCourseReq 更新课程请求体，ID 从路径参数绑定
type UpdateCourseReq struct {
	ID          int64   `json:"id" param:"id"`
	Title       string  `json:"title"`
	Teacher     string  `json:"teacher"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// GetCourseReq 查询课程详情请求，ID 从路径参数绑定
type GetCourseReq struct {
	ID int64 `param:"id"`
}

// DeleteCourseReq 删除课程请求，ID 从路径参数绑定
type DeleteCourseReq struct {
	ID int64 `param:"id"`
}
