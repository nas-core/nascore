package followStartAndCron

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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

	var searchPaths []string
	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("[nascore] Failed to get executable file path, error: %v", err)
		return err
	}
	currentDir := filepath.Dir(executablePath)
	extendedDir := filepath.Join(currentDir, "extended")
	searchPaths = append(searchPaths, currentDir)
	searchPaths = append(searchPaths, extendedDir)

	extendedPath := os.Getenv("NASCOTE_EXTENDED_PATH")
	if extendedPath != "" {
		searchPaths = append(searchPaths, extendedPath)
	}

	if isDevMode.IsDevMode() {
		searchPaths = append(searchPaths, "/home/yh/myworkspace/nas-core/CodeSpace/nascore_vod")
	}

	for _, path := range searchPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			fileName := file.Name()

			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(path, fileName)

			isExecutableCandidate := false
			fileExtension := strings.ToLower(filepath.Ext(fileName))

			switch runtime.GOOS {
			case "windows":
				switch fileExtension {
				case ".exe", ".bat", ".cmd", ".ps1":
					isExecutableCandidate = true
				}
			case "linux", "freebsd", "darwin": // Unix-like systems
				switch fileExtension {
				case "", ".bin", ".sh", ".command": // "" for no extension (common for binaries)
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						logger.Warnf("[nascore] Failed to get file information: %s, error: %v", fileName, err)
						continue
					}
					// 检查是否具有可执行权限
					if fileInfo.Mode().Perm()&0111 != 0 {
						isExecutableCandidate = true
					} else {
						logger.Warnf("[nascore] File does not have execute permission: %s", fileName)
					}
				}
			default:
				logger.Debug("[nascore] Skipping file %s, unsupported operating system: %s", fileName, runtime.GOOS)
				continue
			}

			if !isExecutableCandidate {
				logger.Debug("[nascore] Skipping file %s, does not meet executable file rules for the OS", fileName)
				continue
			}

			switch {
			case strings.Contains(strings.ToLower(fileName), "tv"), strings.Contains(strings.ToLower(fileName), "vod"):
				cmdParams := []string{"-s", socketFilePathValue + system_config.ExtensionSocketMap["nascore_vod"], "-githubDownloadMirror", nsCfg.ThirdPartyExt.GitHubDownloadMirror}
				logger.Debug("[nascore] 🔹Starting execution: %s, parameters: %v", filePath, cmdParams)
				executeIfMatching(filePath, fileName, cmdParams, logger)
			default:
				logger.Debug("[nascore] Skipping file %s, does not match keyword (tv/vod)", fileName)
			}
		}
	}

	return nil
}

// executeIfMatching 执行匹配的文件
func executeIfMatching(filePath string, fileName string, cmdParams []string, logger *zap.SugaredLogger) {
	// 针对 Unix-like 系统再次检查执行权限，防止在某些极端情况下被跳过
	if runtime.GOOS != "windows" {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logger.Errorf("[nascore] 🔸 Failed to get file information, checking before execution: %s, error: %v", fileName, err)
			return
		}
		if fileInfo.Mode().Perm()&0111 == 0 {
			logger.Warnf("[nascore] 🔸 File does not have execute permission (non-Windows system): %s", fileName)
			return
		}
	}

	cmd := exec.Command(filePath, cmdParams...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Errorf("[nascore] 🔸 Execution output: %v", string(output))
		}
		logger.Errorf("[nascore] 🔸 Execution failed: %s, error: %v", filePath, err)
		return
	}
	logger.Debug("[nascore] 🔹 Execution output: %s", string(output))
}

func CheckAllExtensionStatusOnce(nsCfg *system_config.SysCfg) {
	for extName, socketFile := range system_config.ExtensionSocketMap {
		socketPath := nsCfg.Server.UnixSocketFilePath
		if len(socketPath) > 0 && socketPath[len(socketPath)-1] != '/' && socketPath[len(socketPath)-1] != '\\' {
			socketPath += "/"
		}
		socketPath += socketFile

		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
			Timeout: 2 * time.Second,
		}
		req, _ := http.NewRequest("GET", "http://unix/ping", nil)
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			system_config.ExtensionStatusMap[extName] = true
		} else {
			system_config.ExtensionStatusMap[extName] = false
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}
