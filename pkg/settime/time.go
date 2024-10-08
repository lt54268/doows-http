package settime

import (
	"time"
)

// 返回当前时间的格式化字符串
func GetCurrentFormattedTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
