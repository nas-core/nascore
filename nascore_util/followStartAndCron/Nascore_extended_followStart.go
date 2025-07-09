// CodeSpace/nascore/nascore_util/followStartAndCron/Nascore_extended_followStart.go
package followStartAndCron

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nas-core/nascore/nascore_util/isDevMode"
	"github.com/nas-core/nascore/nascore_util/system_config"
	"go.uber.org/zap"
)

// Nascore_extended_followStart 扩展的启动跟踪函数
func Nascore_extended_followStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	socketFilePathValue := nsCfg.Server.UnixSocketFilePath
	if len(socketFilePathValue) > 0 && socketFilePathValue[len(socketFilePathValue)-1] != '/' {
		socketFilePathValue += "/"
	}

	// 搜索 对应目录下的二进制文件 并获取路径 包括.exe
	// 其次搜索 同目录下的子目录 extended 下的 二进制可执行程序 和 .exe
	// 然后搜索 环境变量 NASCOTE_EXTENDED_PATH  目录下的 二进制可执行程序 和 .exe
	// 最后搜索 /home/yh/myworkspace/nas-core/CodeSpace/nascore_vod 目录下的 二进制程序（测试环境）
	// 如果可执行文件的名字包含 tv 或者vod 那么执行命令 是 ./文件名 -s /tmp/nascore.socket
	var searchPaths []string
	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("获取执行文件路径失败: %v", err)
		return err
	}
	currentDir := filepath.Dir(executablePath)
	extendedDir := filepath.Join(currentDir, "extended")
	searchPaths = append(searchPaths, currentDir)
	searchPaths = append(searchPaths, extendedDir)

	// 添加环境变量 NASCOTE_EXTENDED_PATH
	extendedPath := os.Getenv("NASCOTE_EXTENDED_PATH")
	if extendedPath != "" {
		searchPaths = append(searchPaths, extendedPath)
	}

	// 添加测试环境目录
	if isDevMode.IsDevMode() {
		searchPaths = append(searchPaths, "/home/yh/myworkspace/nas-core/CodeSpace/nascore_vod")
	}
	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			logger.Errorf("err read: %s, err: %v", path, err)
			continue // 忽略此目录，继续下一个
		}

		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()
				filePath := filepath.Join(path, fileName)
				cmdParams := []string{} // 命令参数
				switch {
				case strings.Contains(strings.ToLower(fileName), "tv"), strings.Contains(strings.ToLower(fileName), "vod"):
					cmdParams = []string{"-s", socketFilePathValue + system_config.NasCoreVodSocketPath, "-githubDownloadMirror", nsCfg.ThirdPartyExt.GitHubDownloadMirror}
					executeIfMatching(filePath, fileName, cmdParams, logger)
				}

			}
		}
	}

	return nil
}

// executeIfMatching 执行匹配的文件
func executeIfMatching(filePath string, fileName string, cmdParams []string, logger *zap.SugaredLogger) {
	// 检查文件是否是可执行文件
	fileInfo, err := os.Stat(filePath) // 使用 os.Stat 获取文件信息
	if err != nil {
		logger.Errorf("获取文件信息失败: %s, 错误: %v", fileName, err)
		return
	}

	// 检查文件权限，判断是否可执行
	if fileInfo.Mode().Perm()&0111 != 0 {
		cmd := exec.Command(filePath, cmdParams...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if output != nil {
				logger.Errorf("executeIfMatching output: %v", string(output))
			}
			logger.Errorf("executeIfMatching 执行失败: %s, 错误: %v", filePath, err)
			return // 忽略此文件，继续下一个
		}
		logger.Infof("executeIfMatching output: %s", string(output))
	} else {
		logger.Warnf("文件不可执行: %s", fileName)
	}
}
