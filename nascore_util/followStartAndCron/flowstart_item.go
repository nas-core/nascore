package followStartAndCron

import (
	"os/exec"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func DdnsSGOFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	cmd := exec.Command(nsCfg.ThirdPartyExt.DdnsGO.BinPath, "-c", nsCfg.ThirdPartyExt.DdnsGO.ConfigFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Errorf("DdnsSGOFollowStart output: %v", string(output))
		}
		logger.Errorf("DdnsSGOFollowStart failed: %v", err)
		return err
	}
	logger.Infof("DdnsSGOFollowStart output: %s", string(output))

	return nil
}

// ./openlist server --data ./oplist_data
func OpenlistFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	cmd := exec.Command(nsCfg.ThirdPartyExt.Openlist.BinPath, "server", "--data", nsCfg.ThirdPartyExt.Openlist.DataPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Errorf("OpenlistFollowStart output: %v", string(output))
		}
		logger.Errorf("OpenlistFollowStart failed: %v", err)
		return err
	}
	logger.Infof("OpenlistFollowStart output: %s", string(output))

	return nil
}

func Caddy2FollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	cmd := exec.Command(nsCfg.ThirdPartyExt.Caddy2.BinPath, "run", "--config", nsCfg.ThirdPartyExt.Caddy2.ConfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Errorf("Caddy2FollowStart output: %v", string(output))
		}
		logger.Errorf("Caddy2FollowStart failed: %v", err)
		return err
	}
	logger.Infof("Caddy2FollowStart output: %s", string(output))

	return nil
}

// RcloneFollowStart executes rclone mount commands from system configuration
func RcloneFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	logger.Info("rclone auto mount is enable ...")

	exeRcloneAutoMount(nsCfg, logger)

	return nil
}
