package system_config

import (
	"github.com/joyanhui/golang-pkgs/pkgs/exePath"
	"github.com/nas-core/nascore/nascore_util/isDevMode"
)

const (
	PrefixApi         = "/@api/"         // 固定
	PrefixDdnsGo      = "/@ddnsgo/"      // 固定
	PrefixAdguardhome = "/@adguardhome/" // 固定
	PrefixNasCoreTv   = "/@nascore_vod/" // 固定
	PrefixLinks       = "/@links/"       // 固定
	MaxUserLength     = 5
)

var ConfigFilePath string
var DbUserPath = exePath.GetExeDir(isDevMode.IsDevMode()) + "nascore.db" // 固定

// ExtensionStatusMap 用于存储扩展名与其可用状态
var ExtensionStatusMap = make(map[string]bool)

// ExtensionSocketMap 用于存储扩展名与其 socket 路径（为后续多扩展做准备）
var ExtensionSocketMap = map[string]string{
	"nascore_vod": "nascore_tv.socket", // 默认值，后续可由配置文件覆盖
}
