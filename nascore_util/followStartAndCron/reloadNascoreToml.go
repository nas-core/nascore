package followStartAndCron

import (
	"log"

	"github.com/nas-core/nascore/nascore_util/system_config"
)

func reloadNascoreToml(nsCfg *system_config.SysCfg) {
	tmpNsCfg, err := system_config.LoadConfig(system_config.ConfigFilePath)
	if err == nil {
		*nsCfg = *tmpNsCfg
	} else {
		log.Println("hot reload nascore toml file  err", err.Error())
	}
}
