package downfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// getFileNameFromHeader 从 Content-Disposition 提取文件名
func getFileNameFromHeader(header string) string {
	re := regexp.MustCompile(`filename="?([^";]+)"?`)
	matches := re.FindStringSubmatch(header)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// getFileNameFromURL 从 URL 提取文件名
func getFileNameFromURL(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "tmp_download_file"
	}
	path := u.Path
	segments := strings.Split(path, "/")
	name := segments[len(segments)-1]
	if name == "" {
		return "tmp_download_file"
	}
	return name
}

// DownloadFile 下载文件到指定目录，返回实际保存的文件完整路径
func DownloadFile(urlStr string, saveDir string, saveName string) (string, error) {
	// 创建保存目录如果不存在
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		err := os.MkdirAll(saveDir, 0755)
		if err != nil {
			return "", fmt.Errorf("mkdir err: %w", err)
		}
	}

	// 发起 HTTP GET 请求
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("down file GET err: %w", err)
	}
	defer resp.Body.Close()

	// 优先用 Content-Disposition
	if saveName == "" {
		cd := resp.Header.Get("Content-Disposition")
		if cd != "" {
			saveName = getFileNameFromHeader(cd)
		}
		if saveName == "" {
			saveName = getFileNameFromURL(urlStr)
		}
	}

	if saveName == "" {
		return "", fmt.Errorf("down file get file name err")
	}

	// 创建保存文件
	fileSavePath := filepath.Join(saveDir, saveName)
	out, err := os.Create(fileSavePath)
	if err != nil {
		return "", fmt.Errorf("mkfile err: %w", err)
	}
	defer out.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("down file HTTP err core: %d", resp.StatusCode)
	}

	// 将响应体的数据写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("down file write err: %w", err)
	}

	// 查看文件信息
	fileInfo, err := os.Stat(fileSavePath)
	if err != nil {
		return "", fmt.Errorf("down file get file stat err: %w", err)
	}
	fmt.Printf("down file size: %d bit\n", fileInfo.Size())
	// 真实路径
	absPath, err := filepath.Abs(fileSavePath)
	if err != nil {
		return "", fmt.Errorf("down file get file abs path err: %w", err)
	}
	fmt.Printf("down file abs path: %s\n", absPath)
	return absPath, nil
}
