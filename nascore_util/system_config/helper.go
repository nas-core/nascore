package system_config

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

func GenerateStr(typeInt int) string {
	baseStr, err := os.Hostname()
	tmpHash := "nascore.eu.org"

	if err != nil {
		baseStr = tmpHash
	}
	h := md5.New()
	switch typeInt {
	case 1:
		io.WriteString(h, baseStr)
		baseStr = fmt.Sprintf("%x", h.Sum(nil))
	case 2:
		io.WriteString(h, baseStr+tmpHash)
		baseStr = fmt.Sprintf("%x", h.Sum(nil))
	case 3:
		io.WriteString(h, baseStr+tmpHash+"https://nascore.eu.org/api/")
		baseStr = fmt.Sprintf("%x", h.Sum(nil))
	}
	return baseStr
}

// EnsureDirPathSuffix 补全目录路径结尾的 / 或 \
func EnsureDirPathSuffix(path string) string {
	if path == "" {
		return path
	}
	var suf string
	if runtime.GOOS == "windows" || (len(path) > 0 && path[len(path)-1] == '\\') {
		suf = `\\`
	} else {
		suf = `/`
	}
	if !strings.HasSuffix(path, suf) {
		return path + suf
	}
	return path
}

// ensureWindowsExeExtension 保留
func EnsureWindowsExeExtension(sys_cfg *SysCfg) {
	if runtime.GOOS == "windows" {
		// 先判断对应的 BinPath 字段非空再拼接 .exe，避免空指针或重复拼接
		if sys_cfg.ThirdPartyExt.Rclone.BinPath != "" && !strings.HasSuffix(sys_cfg.ThirdPartyExt.Rclone.BinPath, ".exe") {
			sys_cfg.ThirdPartyExt.Rclone.BinPath = sys_cfg.ThirdPartyExt.Rclone.BinPath + ".exe"
		}
		if sys_cfg.ThirdPartyExt.DdnsGO.BinPath != "" && !strings.HasSuffix(sys_cfg.ThirdPartyExt.DdnsGO.BinPath, ".exe") {
			sys_cfg.ThirdPartyExt.DdnsGO.BinPath = sys_cfg.ThirdPartyExt.DdnsGO.BinPath + ".exe"
		}
		if sys_cfg.ThirdPartyExt.AcmeLego.BinPath != "" && !strings.HasSuffix(sys_cfg.ThirdPartyExt.AcmeLego.BinPath, ".exe") {
			sys_cfg.ThirdPartyExt.AcmeLego.BinPath = sys_cfg.ThirdPartyExt.AcmeLego.BinPath + ".exe"
		}
		if sys_cfg.ThirdPartyExt.Caddy2.BinPath != "" && !strings.HasSuffix(sys_cfg.ThirdPartyExt.Caddy2.BinPath, ".exe") {
			sys_cfg.ThirdPartyExt.Caddy2.BinPath = sys_cfg.ThirdPartyExt.Caddy2.BinPath + ".exe"
		}
		if sys_cfg.ThirdPartyExt.Openlist.BinPath != "" && !strings.HasSuffix(sys_cfg.ThirdPartyExt.Openlist.BinPath, ".exe") {
			sys_cfg.ThirdPartyExt.Openlist.BinPath = sys_cfg.ThirdPartyExt.Openlist.BinPath + ".exe"
		}
	}
}
