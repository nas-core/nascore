package followStartAndCron

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func execLegoRenewOrGet(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	legoLogFile := nsCfg.ThirdPartyExt.AcmeLego.LEGO_PATH + "/lego_execLegoRenewOrGet.log"
	commandStr := nsCfg.ThirdPartyExt.AcmeLego.Command
	commandStr = strings.ReplaceAll(commandStr, "${BinPath}", nsCfg.ThirdPartyExt.AcmeLego.BinPath)
	commandStr = strings.ReplaceAll(commandStr, "${LEGO_PATH}", nsCfg.ThirdPartyExt.AcmeLego.LEGO_PATH)
	stdoutArr, stderrArr, errArr := excMultiLineCommand_Sequentially(&commandStr, logger, legoLogFile)
	logger.Debug(" execLegoCommand err len", len(errArr), " err ", errArr)
	logger.Debug(" execLegoCommand stdoutArr len ", len(stdoutArr), " stdoutArr", stdoutArr)
	logger.Debug(" execLegoCommand stdoutArr len", len(stderrArr), " stderrArr ", stderrArr)

}

func exeRcloneAutoMount(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	commandStr1 := nsCfg.ThirdPartyExt.Rclone.AutoUnMountCommand
	commandStr2 := nsCfg.ThirdPartyExt.Rclone.AutoMountCommand
	commandStr2 = strings.ReplaceAll(commandStr2, "${BinPath}", nsCfg.ThirdPartyExt.Rclone.BinPath)

	stdoutArr, stderrArr, errArr := excMultiLineCommand_Sequentially(&commandStr1, logger, "")
	logger.Debug(" exeRcloneAutoUnMount err len", len(errArr), " err ", errArr)
	logger.Debug(" exeRcloneAutoUnMount stdoutArr len ", len(stdoutArr), " stdoutArr", stdoutArr)
	logger.Debug(" exeRcloneAutoUnMount stdoutArr len", len(stderrArr), " stderrArr ", stderrArr)
	stdoutArr, stderrArr, errArr = excMultiLineCommand_Sequentially(&commandStr2, logger, "")
	logger.Debug(" exeRcloneAutoMount err len", len(errArr), " err ", errArr)
	logger.Debug(" exeRcloneAutoMount stdoutArr len ", len(stdoutArr), " stdoutArr", stdoutArr)
	logger.Debug(" exeRcloneAutoMount stdoutArr len", len(stderrArr), " stderrArr ", stderrArr)
}

func excMultiLineCommand_Sequentially(commandStr *string, logger *zap.SugaredLogger, logFile string) (stdoutArr []string, stderrArr []string, errArr []error) {
	lines := strings.Split(*commandStr, "\n")
	var envs []string
	for _, line := range lines {
		var cmdName string
		var cmdArgs []string
		line = strings.TrimSpace(line) // 去除行首尾空格
		if line == "" {
			continue
		}
		// 兼容 export 和 windows 下 set 方式设置环境变量
		if envVar, ok := strings.CutPrefix(line, "export "); ok {
			envs = append(envs, envVar)
		} else if envVar, ok := strings.CutPrefix(line, "set "); ok { // 解析 windows 下 set 方式的环境变量
			envs = append(envs, envVar)
		} else {
			parts := strings.Fields(line) // 使用 Fields 分割命令和参数
			if len(parts) > 0 {
				cmdName = parts[0]
				cmdArgs = parts[1:]
			}
		}
		if cmdName != "" {
			if len(cmdArgs) > 0 && cmdArgs[len(cmdArgs)-1] == "&nascore" { // 移除 cmdArgs 中结尾的 &nascore
				cmdArgs = cmdArgs[:len(cmdArgs)-1]
			}
			cmd := exec.Command(cmdName, cmdArgs...)
			cmd.Env = append(os.Environ(), envs...) // 继承当前环境，并添加新的环境变量
			var stdin, stdout, stderr bytes.Buffer
			cmd.Stdin = &stdin
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if strings.HasSuffix(strings.TrimSpace(line), "&nascore") { // 检查是否以 &nascore 结尾
				logger.Debug("MultiLineCommand Asynchronous execution   ", line)
				go func(cmdLine string) { // 使用 goroutine 执行命令
					err := cmd.Run()
					if err != nil {
						errArr = append(errArr, err)
					}
					stdoutArr = append(stdoutArr, stdout.String())
					stderrArr = append(stderrArr, stderr.String())
					writeLogToFile(logFile, cmdLine, stdout.String(), stderr.String())
				}(line)
			} else {
				logger.Debug("MultiLineCommand Sequential execution   ", line)

				err := cmd.Run() // 顺序执行命令
				if err != nil {
					errArr = append(errArr, err)
				}
				stdoutArr = append(stdoutArr, stdout.String())
				stderrArr = append(stderrArr, stderr.String())
				writeLogToFile(logFile, line, stdout.String(), stderr.String())
			}
		}
	}
	return stdoutArr, stderrArr, errArr
}

// writeLogToFile 将命令执行的输出写入到指定日志文件
func writeLogToFile(logFile, cmdLine, outStr, errStr string) {
	if logFile != "" {
		f, ferr := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if ferr == nil {
			defer f.Close()
			f.WriteString("\n[cmd] " + cmdLine + "\n")
			if outStr != "" {
				f.WriteString("[stdout]\n" + outStr + "\n")
			}
			if errStr != "" {
				f.WriteString("[stderr]\n" + errStr + "\n")
			}
		}
	}
}
