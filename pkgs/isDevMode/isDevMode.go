package isDevMode

import "os"

// 从环境变量检查开发模式
/*
开发模式下 禁用 webui的 嵌入式fs 方便刷新webui测试
开发模式下 日志打印会更多
*/
func IsDevMode() bool {
	return os.Getenv("nascore_DEV_MODE") == "1" || os.Getenv("nascore_DEV_MODE") == "true"
}
