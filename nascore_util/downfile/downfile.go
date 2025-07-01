package downfile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile 下载文件到指定目录
func DownloadFile(url string, saveDir string, saveName string) error {
	// 创建保存目录如果不存在
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		err := os.MkdirAll(saveDir, 0755)
		if err != nil {
			return fmt.Errorf("mkdir err: %w", err)
		}
	}

	// 创建保存文件
	fileSavePath := filepath.Join(saveDir, saveName)
	out, err := os.Create(fileSavePath)
	if err != nil {
		return fmt.Errorf("mkfile err: %w", err)
	}
	defer out.Close()

	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		log.Println("url ", url)
		return fmt.Errorf("down file GET err: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		log.Println("url ", url)
		return fmt.Errorf("down file HTTP err core: %d", resp.StatusCode)
	}

	// 将响应体的数据写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("down file write err: %w", err)
	}
	// 查看文件信息
	fileInfo, err := os.Stat(fileSavePath)
	if err != nil {
		return fmt.Errorf("down file get file stat err: %w", err)
	}
	fmt.Printf("down file size: %d bit\n", fileInfo.Size())
	// 真实路径
	absPath, err := filepath.Abs(fileSavePath)
	if err != nil {
		return fmt.Errorf("down file get file abs path err: %w", err)
	}
	fmt.Printf("down file abs path: %s\n", absPath)
	return nil
}
