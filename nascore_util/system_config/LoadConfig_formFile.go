package system_config

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*SysCfg, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("viper.ReadInConfig file failed: ", err)
	}
	config := NewDefaultConfig() // 初始化 config 为指针类型
	err := viper.Unmarshal(config)
	if err != nil {
		log.Println("viper.Unmarshal failed: ", err)
	}

	// 统一补全目录路径结尾
	config.Server.TempFilePath = EnsureDirPathSuffix(config.Server.TempFilePath)
	config.ThirdPartyExt.Openlist.DataPath = EnsureDirPathSuffix(config.ThirdPartyExt.Openlist.DataPath)
	config.ThirdPartyExt.AcmeLego.LEGO_PATH = EnsureDirPathSuffix(config.ThirdPartyExt.AcmeLego.LEGO_PATH)

	if config.Secret.JwtSecret == "" {
		config.Secret.JwtSecret = GenerateStr(1)
		log.Println("config.Secret.JwtSecret is empty set :", config.Secret.JwtSecret)
	}
	if config.Secret.Sha256HashSalt == "" {
		config.Secret.Sha256HashSalt = GenerateStr(2)
		log.Println("config.Secret.Sha256HashSalt is empty set :", config.Secret.Sha256HashSalt)
	}
	if config.Secret.AESkey == "" {
		config.Secret.AESkey = GenerateStr(3)
		log.Println("config.Secret.AESkey is empty set :", config.Secret.AESkey)
	}

	return config, err
}
