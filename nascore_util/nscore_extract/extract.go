package nscore_extract

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Go内置支持的解压函数

// ExtractTar 解压 tar 格式文件
func ExtractTar(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	return extractTarContent(tarReader, targetPath)
}

// ExtractTarGz 解压 tar.gz 格式文件
func ExtractTarGz(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	return extractTarContent(tarReader, targetPath)
}

// ExtractTarBz2 解压 tar.bz2 格式文件
func ExtractTarBz2(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)
	tarReader := tar.NewReader(bz2Reader)
	return extractTarContent(tarReader, targetPath)
}

// ExtractZip 解压 zip 格式文件
func ExtractZip(sourcePath, targetPath string) error {
	reader, err := zip.OpenReader(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// 构建目标文件路径
		targetFilePath := filepath.Join(targetPath, file.Name)

		// 安全检查：确保目标路径在目标目录内（防止目录遍历攻击）
		if !strings.HasPrefix(targetFilePath, filepath.Clean(targetPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in archive: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			// 创建目录
			if err := os.MkdirAll(targetFilePath, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetFilePath, err)
			}
			continue
		}

		// 创建文件
		if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", targetFilePath, err)
		}

		// 打开zip文件中的文件
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		// 创建目标文件
		outFile, err := os.Create(targetFilePath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
		}

		// 复制文件内容
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", targetFilePath, err)
		}

		// 设置文件权限
		if err := os.Chmod(targetFilePath, file.FileInfo().Mode()); err != nil {
			return fmt.Errorf("failed to set permissions for %s: %w", targetFilePath, err)
		}
	}

	return nil
}

// ExtractGz 解压 gz 格式文件
func ExtractGz(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// 获取原始文件名（去掉.gz扩展名）
	originalName := strings.TrimSuffix(filepath.Base(sourcePath), ".gz")
	outputPath := filepath.Join(targetPath, originalName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, gzReader)
	return err
}

// extractBz2 解压 bz2 格式文件
func extractBz2(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)

	// 获取原始文件名（去掉.bz2扩展名）
	originalName := strings.TrimSuffix(filepath.Base(sourcePath), ".bz2")
	outputPath := filepath.Join(targetPath, originalName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, bz2Reader)
	return err
}

// extractTarZ 解压 tar.z (tar+zlib) 格式文件
func ExtractTarZ(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建zlib reader
	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer zlibReader.Close()

	// 创建tar reader
	tarReader := tar.NewReader(zlibReader)

	return extractTarContent(tarReader, targetPath)
}

// extractTarContent 提取tar内容的通用函数
func extractTarContent(tarReader *tar.Reader, targetPath string) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// 构建目标文件路径
		targetFilePath := filepath.Join(targetPath, header.Name)

		// 安全检查：确保目标路径在目标目录内（防止目录遍历攻击）
		if !strings.HasPrefix(targetFilePath, filepath.Clean(targetPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(targetFilePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetFilePath, err)
			}
		case tar.TypeReg:
			// 创建文件
			if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetFilePath, err)
			}

			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to extract file %s: %w", targetFilePath, err)
			}

			outFile.Close()

			// 设置文件权限
			if err := os.Chmod(targetFilePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set permissions for %s: %w", targetFilePath, err)
			}
		default:
			// 跳过其他类型的文件（符号链接等）
			continue
		}
	}

	return nil
}
