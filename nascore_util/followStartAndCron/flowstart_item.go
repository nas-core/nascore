package followStartAndCron

import (
	"github.com/nas-core/nascore/nascore_util/exeStart"
	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func DdnsSGOFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	exeStart.KillDDNSGo(nsCfg, logger)
	err = exeStart.StartDDNSGo(nsCfg, logger)
	return err
}

// ./openlist server --data ./oplist_data
func OpenlistFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	exeStart.KillOpenlist(nsCfg, logger)
	err = exeStart.StartOpenlist(nsCfg, logger)
	return err
}

func Caddy2FollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	exeStart.KillCaddy2(nsCfg, logger)
	err = exeStart.StartCaddy2(nsCfg, logger)
	return err
}

// RcloneFollowStart executes rclone mount commands from system configuration
func RcloneFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	exeRcloneAutoMount(nsCfg, logger)

	return nil
}
