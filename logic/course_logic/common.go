// course_logic package: common define.
package course_logic

import (
	"fmt"

	"echo-learning/dao"
)

// errRecordNotFound 生成「记录不存在」错误描述，包装 dao.ErrNotFound 便于 errors.Is 判断
func errRecordNotFound(ID int64) error {
	return fmt.Errorf("%w: course id=%d not found", dao.ErrNotFound, ID)
}
