package nscore_extract

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// extractXzWithSystemCommand 解压xz格式文件
func extractXzWithSystemCommand(sourcePath, targetPath string, logger *zap.SugaredLogger) error {
	// 获取原始文件名（去掉.xz扩展名）
	originalName := strings.TrimSuffix(filepath.Base(sourcePath), ".xz")
	outputPath := filepath.Join(targetPath, originalName)

	cmd := exec.Command("xz", "-dc", sourcePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to decompress xz file: %w", err)
	}

	// 写入解压后的文件
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decompressed file: %w", err)
	}

	return nil
}

// extractLzmaWithSystemCommand 解压lzma格式文件
func extractLzmaWithSystemCommand(sourcePath, targetPath string, logger *zap.SugaredLogger) error {
	// 获取原始文件名
	originalName := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	outputPath := filepath.Join(targetPath, originalName)

	var cmd *exec.Cmd
	if _, err := exec.LookPath("xz"); err == nil {
		cmd = exec.Command("xz", "-dc", sourcePath)
	} else if _, err := exec.LookPath("lzma"); err == nil {
		cmd = exec.Command("lzma", "-dc", sourcePath)
	} else {
		return fmt.Errorf("xz or lzma command not found, please install xz-utils package")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to decompress lzma file: %w", err)
	}

	// 写入解压后的文件
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decompressed file: %w", err)
	}

	return nil
}

// extractZWithSystemCommand 解压Z格式文件
func extractZWithSystemCommand(sourcePath, targetPath string, logger *zap.SugaredLogger) error {
	// 获取原始文件名（去掉.Z扩展名）
	originalName := strings.TrimSuffix(filepath.Base(sourcePath), ".Z")
	outputPath := filepath.Join(targetPath, originalName)

	var cmd *exec.Cmd
	if _, err := exec.LookPath("uncompress"); err == nil {
		// 使用uncompress命令
		cmd = exec.Command("uncompress", "-c", sourcePath)
	} else if _, err := exec.LookPath("gzip"); err == nil {
		// gzip也可以处理Z格式
		cmd = exec.Command("gzip", "-dc", sourcePath)
	} else {
		return fmt.Errorf("uncompress or gzip command not found, please install gzip package")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to decompress Z file: %w", err)
	}

	// 写入解压后的文件
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decompressed file: %w", err)
	}

	return nil
}

// 系统命令支持的解压函数

// extractWithSystemCommand 使用系统命令解压文件
func extractWithSystemCommand(sourcePath, targetPath, extension string, logger *zap.SugaredLogger) error {
	var cmd *exec.Cmd

	switch extension {
	case "tar.xz":
		// 使用tar命令直接解压tar.xz
		cmd = exec.Command("tar", "-xJf", sourcePath, "-C", targetPath)
	case "xz":
		// 先解压xz，然后检查是否为tar文件
		return extractXzWithSystemCommand(sourcePath, targetPath, logger)
	case "7z":
		// 使用7z命令
		if _, err := exec.LookPath("7z"); err == nil {
			cmd = exec.Command("7z", "x", sourcePath, "-o"+targetPath, "-y")
		} else if _, err := exec.LookPath("7za"); err == nil {
			cmd = exec.Command("7za", "x", sourcePath, "-o"+targetPath, "-y")
		} else {
			return fmt.Errorf("7z command not found, please install p7zip package")
		}
	case "rar":
		// 使用unrar命令
		if _, err := exec.LookPath("unrar"); err == nil {
			cmd = exec.Command("unrar", "x", sourcePath, targetPath)
		} else {
			return fmt.Errorf("unrar command not found, please install unrar package")
		}
	case "lz", "lzma":
		// 使用xz工具解压lzma格式
		return extractLzmaWithSystemCommand(sourcePath, targetPath, logger)
	case "Z":
		// 使用uncompress或gzip解压Z格式
		return extractZWithSystemCommand(sourcePath, targetPath, logger)
	default:
		return fmt.Errorf("unsupported system command extraction for format: %s", extension)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("System command extraction failed", zap.String("command", cmd.String()), zap.String("output", string(output)))
		return fmt.Errorf("failed to extract %s file: %w, output: %s", extension, err, string(output))
	}

	logger.Debug("System command extraction successful", zap.String("format", extension))
	return nil
}
