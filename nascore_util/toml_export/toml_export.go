package toml_export

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"github.com/pelletier/go-toml/v2"
)

func Export(cfg *system_config.SysCfg, exportConfigPath *string) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.SetIndentTables(true)
	err := encoder.Encode(cfg)
	if err != nil {
		log.Println("toml Encode err", err)
		return err
	}

	if exportConfigPath == nil || *exportConfigPath == "" {
		log.Println("exportConfigPath is empty or nil, cannot write to file.")
		return err
	}

	err = os.WriteFile(*exportConfigPath, buf.Bytes(), 0644) // 0644 是文件权限
	if err != nil {
		log.Printf("The first time to write to %s failed: %v", *exportConfigPath, err)
		return err
	}

	content, err := os.ReadFile(*exportConfigPath)
	if err != nil {
		log.Printf("re read toml %s failed: %v", *exportConfigPath, err)
		return err
	}

	re := regexp.MustCompile(`(?m)^(\s*[a-zA-Z_]+\s*=\s*)"((?:[^"\\]|\\.)*\\n(?:[^"\\]|\\.)*)"(\s*)$`)
	content = re.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := strings.SplitN(string(match), "=", 2)
		if len(parts) != 2 {
			return match
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 去掉两边的双引号
		value = strings.Trim(value, `"`)

		// 替换转义字符
		value = strings.ReplaceAll(value, `\"`, `"`)
		value = strings.ReplaceAll(value, `\n`, "\n")

		// 重新包裹为 '''
		return []byte(key + " = '''\n" + value + "\n'''")
	})

	if strings.Contains(string(content), "[Secret]") {
		content = []byte(strings.Replace(string(content), "[Secret]", "# Some keys used for encryption will be automatically generated if they are empty, but this may cause the login status or password to become invalid after a restart\n[Secret]", 1))
	}

	if strings.Contains(string(content), "[[Users]]") {
		content = []byte(strings.Replace(string(content), "[[Users]]", "# Permissions Is Json string\n# RUCc R(ReadFile) U(UpdateFile) C(CreateFile) c(CreateDir) D(DeleteFile)\n[[Users]]", 1))
	}

	// 重新写入文件
	err = os.WriteFile(*exportConfigPath, content, 0644)
	if err != nil {
		log.Printf("update TOML file %s failed: %v", *exportConfigPath, err)
		return err
	}

	return nil
}
