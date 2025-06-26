package server_http

import (
	"flag"
	"log"
	"nascore_v3/pkgs/toml_export"
	"nascore_v3/system_config"
	"os"
)

func initFlagsAndSysCfg() (sys_config *system_config.SysCfg, configFilePath string, err error) {
	configFilePathFlag := flag.String("c", "", "config file (exp: -c /etc/config.toml)")
	flag.Parse()
	sys_config, configFilePath, err = getConfigAndFilePathOrCreatConfigFile(configFilePathFlag)
	return
}

func getConfigAndFilePathOrCreatConfigFile(configFilePathInput *string) (sys_config *system_config.SysCfg, configFilePath string, errFun error) {
	configFilePath = ""
	isNeedCreate := false
	if *configFilePathInput != "" { // 检查文件是否存在
		if _, err := os.Stat(*configFilePathInput); os.IsNotExist(err) {
			log.Fatal("config file does not exist", *configFilePathInput)
			return nil, "", err
		}
		configFilePath = *configFilePathInput
	} else {
		if _, err := os.Stat("./nascore.toml"); os.IsNotExist(err) {
			if _, err := os.Stat("~/.config/nascore.toml"); os.IsNotExist(err) {
				if _, err := os.Stat("/etc/nascore.toml"); os.IsNotExist(err) {
					log.Println("config file not found: ./nascore.toml ~/.config/nascore.toml /etc/nascore.toml")
					isNeedCreate = true
					configFilePath = "/etc/nascore.toml"
				}
			}
			configFilePath = "~/.config/nascore.toml"
			log.Println("config file is found: ", configFilePath)
		}
		configFilePath = "./nascore.toml"
		log.Println("config file is found: ", configFilePath)
	}
	if isNeedCreate {
		err := toml_export.Export(system_config.NewDefaultConfig(), &configFilePath)
		if err != nil {
			log.Fatal("failed to create config file:", err)
			return nil, "", err
		}
	}
	sys_config, err := system_config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatal("failed to load config file:", err)
		return nil, "", err
	}
	return sys_config, configFilePath, nil
}
