package followStartAndCron

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func DdnsSGOFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	binPath := nsCfg.ThirdPartyExt.DdnsGO.BinPath
	configPath := nsCfg.ThirdPartyExt.DdnsGO.ConfigFilePath
	cmd := exec.Command(binPath, "-c", configPath)
	err = cmd.Start()
	if err != nil {
		logger.Error("DdnsSGOFollowStart failed: %v", err)
		return err
	}
	// 生成pid文件路径：配置文件路径+.pid
	pidFile := nsCfg.Server.TempFilePath + "ddnsgo.pid"
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	os.WriteFile(pidFile, []byte(pidStr), 0644)
	go func() {
		cmd.Wait()
		os.Remove(pidFile)
	}()
	logger.Debug("DdnsSGOFollowStart started, pid: %s, pidfile: %s", pidStr, pidFile)
	return nil
}

// ./openlist server --data ./oplist_data
func OpenlistFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	cmd := exec.Command(nsCfg.ThirdPartyExt.Openlist.BinPath, "server", "--data", nsCfg.ThirdPartyExt.Openlist.DataPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Warn("OpenlistFollowStart output: %v", string(output))
		}
		logger.Warn("OpenlistFollowStart failed: %v", err)
		return err
	}
	logger.Debug("OpenlistFollowStart output: %s", string(output))

	return nil
}

func Caddy2FollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	cmd := exec.Command(nsCfg.ThirdPartyExt.Caddy2.BinPath, "run", "--config", nsCfg.ThirdPartyExt.Caddy2.ConfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			logger.Warn("Caddy2FollowStart output: %v", string(output))
		}
		logger.Warn("Caddy2FollowStart failed: %v", err)
		return err
	}
	logger.Debug("Caddy2FollowStart output: %s", string(output))

	return nil
}

// RcloneFollowStart executes rclone mount commands from system configuration
func RcloneFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) (err error) {
	exeRcloneAutoMount(nsCfg, logger)

	return nil
}
