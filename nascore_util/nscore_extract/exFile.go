package nscore_extract

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// ExtractFile 根据文件类型执行解压操作
func ExtractFile(sourcePath string, targetPath string, logger *zap.SugaredLogger) error {
	sourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	if sourcePath == "" {
		return fmt.Errorf("source path is empty")
	}
	extension := GetFileExtension(sourcePath)
	switch extension {
	// Go内置支持的格式
	case "tarz", "tar.z":
		return ExtractTarZ(sourcePath, targetPath)
	case "tar":
		return ExtractTar(sourcePath, targetPath)
	case "tar.gz":
		return ExtractTarGz(sourcePath, targetPath)
	case "tar.bz2":
		return ExtractTarBz2(sourcePath, targetPath)
	case "zip":
		return ExtractZip(sourcePath, targetPath)
	case "gz":
		return ExtractGz(sourcePath, targetPath)

	// 需要系统命令支持的格式
	case "tar.xz", "xz":
		return extractWithSystemCommand(sourcePath, targetPath, extension, logger)
	case "7z":
		return extractWithSystemCommand(sourcePath, targetPath, extension, logger)
	case "rar":
		return extractWithSystemCommand(sourcePath, targetPath, extension, logger)
	case "bz2":
		return extractBz2(sourcePath, targetPath)
	case "lz", "lzma":
		return extractWithSystemCommand(sourcePath, targetPath, extension, logger)
	case "Z":
		return extractWithSystemCommand(sourcePath, targetPath, extension, logger)

	default:
		return fmt.Errorf("unsupported file format: %s", extension)
	}
}

// GetFileExtension 获取文件扩展名，处理多重扩展名
func GetFileExtension(filename string) string {
	filename = strings.ToLower(filename)

	// 处理多重扩展名
	if strings.HasSuffix(filename, ".tar.gz") {
		return "tar.gz"
	}
	if strings.HasSuffix(filename, ".tar.xz") {
		return "tar.xz"
	}
	if strings.HasSuffix(filename, ".tar.bz2") {
		return "tar.bz2"
	}
	if strings.HasSuffix(filename, ".tar.z") || strings.HasSuffix(filename, ".tarz") {
		return "tarz"
	}

	// 单一扩展名
	ext := filepath.Ext(filename)
	if ext != "" {
		ext = ext[1:] // 去掉开头的点
		// 检查特殊格式
		switch ext {
		case "zip", "7z", "gz", "tar", "xz", "bz2", "rar", "lz", "lzma", "Z":
			return ext
		default:
			return ext
		}
	}

	return ""
}
