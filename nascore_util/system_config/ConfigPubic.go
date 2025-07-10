package system_config

const (
	PrefixApi           = "/@api/"         // 固定
	PrefixDdnsGo        = "/@ddnsgo/"      // 固定
	PrefixAdguardhome   = "/@adguardhome/" // 固定
	PrefixNasCoreTv     = "/@nascore_vod/" // 固定
	MaxUserLength       = 5
	NasCoreTvSocketFile = "nascore_tv.socket" // 固定 需要和跟随启动配合
)

var ConfigFilePath string
