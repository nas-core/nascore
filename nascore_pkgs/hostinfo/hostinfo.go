package hostinfo

import (
	"os"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/host"
)

type HostSystemInfo struct {
	Hostid   string `json:"hostid"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"` // amd64
	GoVer    string `json:"gover"`
	Platform string `json:"platform"`
}

// GetHostSystemInfo 获取主机系统信息
func GetHostSystemInfo() (HostSystemInfo, error) {
	// 强制设置环境变量，优先使用 /etc/machine-id 避免hostid重复的问题
	/*
	 GOPSUTIL_HOST_ID_LINUX_SOURCES=productuuid,machineid,bootid would be the default and would offer the current behavior,
	 GOPSUTIL_HOST_ID_LINUX_SOURCES=machineid,bootid would mean "do not use the /sys/class/dmi/id/product_uuid path",
	 GOPSUTIL_HOST_ID_LINUX_SOURCES=productuuid would mean "only use the /sys/class/dmi/id/product_uuid path".

	*/
	os.Setenv("GOPSUTIL_HOST_ID_LINUX_SOURCES", "machineid,bootid")

	hostInfo, err := host.Info()
	if err != nil {
		return HostSystemInfo{}, err
	}

	systemInfo := HostSystemInfo{
		Hostid:   hostInfo.HostID,
		Hostname: hostInfo.Hostname,
		OS:       strings.ToLower(runtime.GOOS),
		Platform: strings.ToLower(hostInfo.Platform),
		Arch:     strings.ToLower(runtime.GOARCH),
		GoVer:    strings.ToLower(runtime.Version()),
	}
	return systemInfo, nil
}
