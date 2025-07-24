package exeStart

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/nas-core/nascore/nascore_util/system_config"
	"go.uber.org/zap"
)

var (
	PidFileCaddy2   = "nascore_caddy2.pid"
	PidFileOpenlist = "nascore_openlist.pid"
	PidFileDDNSGo   = "nascore_ddnsgo.pid"
)

// 通用杀死进程
func killByPidFile(pidFile string, logger *zap.SugaredLogger) {
	pidData, err := os.ReadFile(pidFile)
	if err == nil {
		pid, _ := strconv.Atoi(strings.TrimSpace(string(pidData)))
		if pid > 0 {
			if runtime.GOOS == "windows" {
				exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid)).Run()
			} else {
				exec.Command("kill", fmt.Sprintf("%d", pid)).Run()
			}
		}
	} else {
		logger.Warn("[killByPidFile]", err.Error())
	}
	err = os.Remove(pidFile)
	if err != nil {
		logger.Warn("[killByPidFile] [os.remove] err", err.Error())
	}

}

func KillCaddy2(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	killByPidFile(nsCfg.Server.TempFilePath+PidFileCaddy2, logger)
}
func KillOpenlist(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	killByPidFile(nsCfg.Server.TempFilePath+PidFileOpenlist, logger)
}
func KillDDNSGo(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	killByPidFile(nsCfg.Server.TempFilePath+PidFileDDNSGo, logger)
}

// 启动Caddy2
func StartCaddy2(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) error {
	pidFile := nsCfg.Server.TempFilePath + PidFileCaddy2
	cmd := exec.Command(nsCfg.ThirdPartyExt.Caddy2.BinPath, "run", "--config", nsCfg.ThirdPartyExt.Caddy2.ConfigPath)
	err := cmd.Start()
	if err != nil {
		os.Remove(pidFile)
		logger.Warn("[StartCaddy2] failed: %v", err)
		return err
	}
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(pidStr), 0644)
	if err != nil {
		logger.Warn("[StartCaddy2] [os.writeFile] err", err.Error())
	}
	go func() {
		cmd.Wait()
		os.Remove(pidFile)
	}()
	logger.Debug("[StartCaddy2] started, pid: %s, pidfile: %s", pidStr, pidFile)
	return nil
}

// 启动Openlist
func StartOpenlist(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) error {
	pidFile := nsCfg.Server.TempFilePath + PidFileOpenlist
	cmd := exec.Command(nsCfg.ThirdPartyExt.Openlist.BinPath, "server", "--data", nsCfg.ThirdPartyExt.Openlist.DataPath)
	err := cmd.Start()
	if err != nil {
		os.Remove(pidFile)
		logger.Warn("[StartOpenlist] failed: %v", err)
		return err
	}
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(pidStr), 0644)
	if err != nil {
		logger.Warn("[StartOpenlist] [os.writeFile] err", err.Error())
	}
	go func() {
		cmd.Wait()
		os.Remove(pidFile)
	}()
	logger.Debug("[StartOpenlist] started, pid: %s, pidfile: %s", pidStr, pidFile)
	return nil
}

// 启动DDNSGo
func StartDDNSGo(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) error {
	pidFile := nsCfg.Server.TempFilePath + PidFileDDNSGo
	cmd := exec.Command(nsCfg.ThirdPartyExt.DdnsGO.BinPath, "-c", nsCfg.ThirdPartyExt.DdnsGO.ConfigFilePath)
	err := cmd.Start()
	if err != nil {
		os.Remove(pidFile)
		logger.Warn("[StartDDNSGo] failed: %v", err)
		return err
	}
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(pidStr), 0644)
	if err != nil {
		logger.Warn("[StartDDNSGo] [os.writeFile] err", err.Error())
	}
	go func() {
		cmd.Wait()
		os.Remove(pidFile)
	}()
	logger.Debug("[StartDDNSGo] started, pid: %s, pidfile: %s", pidStr, pidFile)
	return nil
}
