// dao package define: entity / interface / table instance.
package dao

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

// ErrNotFound record not exist.
var ErrNotFound = errors.New("record not found")

// 表名常量
const (
	TableNameCourse = "courses"
)

// Course courses 表实体 —— 本教学项目唯一的一张表
type Course struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"` // 主键，自增
	Title       string    `gorm:"column:title;size:128;not null" json:"title"`
	Teacher     string    `gorm:"column:teacher;size:64" json:"teacher"`
	Description string    `gorm:"column:description;size:512" json:"description"`
	Price       float64   `gorm:"column:price;default:0" json:"price"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"` // GORM 自动维护
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"` // GORM 自动维护
}

// CourseDao course 表的数据访问接口
type CourseDao interface {
	QueryByID(ctx echo.Context, ID int64) (*Course, error)
	QueryAll(ctx echo.Context) ([]*Course, error)
	Insert(ctx echo.Context, course *Course) error
	Update(ctx echo.Context, course *Course) error
	Delete(ctx echo.Context, ID int64) error
}

// CourseTable 全局表实例，logic 层直接通过它访问数据库
var CourseTable CourseDao = &courseDaoImpl{tableName: TableNameCourse}
