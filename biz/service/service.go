package service

import (
	"buycar/pkg/constants"
	"buycar/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetUidFormContext 从请求上下文中获取用户ID
func GetUidFormContext(c *app.RequestContext) int64 {
	uid, _ := c.Get(constants.ContextUid)
	userid, err := convertToInt64(uid)
	if err != nil {
		panic(err)
	}

	return userid
}

// convertToInt64 将各种数值类型转换为int64
func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, errno.NewErrNo(errno.InternalServiceErrorCode, "无法转换为int64类型")
	}
}
