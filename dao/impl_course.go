// dao package: course table implement.
package dao

import (
	"github.com/labstack/echo/v4"
)

// courseDaoImpl courses 表的数据访问实现
type courseDaoImpl struct {
	tableName string
}

// QueryByID 按主键查询单条记录，不存在时返回零值对象（ID == 0），由 logic 层判断
func (i *courseDaoImpl) QueryByID(ctx echo.Context, ID int64) (*Course, error) {
	course := &Course{}
	result := GetSession(ctx).Table(i.tableName).Where("id = ?", ID).Find(course)
	if result.Error != nil {
		return nil, result.Error
	}
	return course, nil
}

// QueryAll 查询全部课程，按 ID 倒序（新的在前）
func (i *courseDaoImpl) QueryAll(ctx echo.Context) ([]*Course, error) {
	courses := make([]*Course, 0)
	result := GetSession(ctx).Table(i.tableName).Order("id DESC").Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}

// Insert 新增一条课程记录
func (i *courseDaoImpl) Insert(ctx echo.Context, course *Course) error {
	return GetSession(ctx).Table(i.tableName).Create(course).Error
}

// Update 按主键更新整条记录
func (i *courseDaoImpl) Update(ctx echo.Context, course *Course) error {
	return GetSession(ctx).Table(i.tableName).Save(course).Error
}

// Delete 按主键删除一条记录
func (i *courseDaoImpl) Delete(ctx echo.Context, ID int64) error {
	return GetSession(ctx).Table(i.tableName).Where("id = ?", ID).Delete(&Course{}).Error
}
