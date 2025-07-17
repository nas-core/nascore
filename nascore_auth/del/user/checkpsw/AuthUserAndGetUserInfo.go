package checkpsw

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nas-core/nascore/nascore_pkgs/getip"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

var (
	failedLoginAttempts = make(map[string]int)
	failedLoginMutex    sync.Mutex
)

func AuthUserAndGetUserInfo(r *http.Request, logger *zap.SugaredLogger, sys_cfg *system_config.SysCfg, username string, passwd string) (ok bool, homepath string, err error) {
	defer logger.Sync()
	user_ip := getip.GetClientIP(r)
	var sleep_time time.Duration
	ok = false
	homepath = ""
	for _, user := range sys_cfg.Users {
		if user["username"] == username {
			if strings.HasPrefix(user["passwd"], "sha256:") {
				hashedPasswd := ComputePasswdSHA256Hash(passwd, sys_cfg.Secret.Sha256HashSalt)
				if user["passwd"] == "sha256:"+hashedPasswd {
					ok = true
				}
			} else {
				if user["passwd"] == passwd {
					ok = true
				}
			}
			if ok {
				logger.Info("user login success")
				failedLoginMutex.Lock()
				if _, hasRecordIP := failedLoginAttempts[user_ip]; hasRecordIP {
					sleep_time = get_sleep_time(user_ip, sys_cfg.Limit.MaxFailedLoginSleepTimeSec)
					delete(failedLoginAttempts, user_ip)
				}
				failedLoginMutex.Unlock()
				logger.Warn("user login success,but need sleep:", sleep_time)
				// 如果 sleep_time >0
				if sleep_time > 0 {
					time.Sleep(sleep_time)
				}
				return true, user["home"], nil
			}
		}
	}
	logger.Warnln("nascore user login failed")
	recordFailedLogin(user_ip, sys_cfg.Limit.MaxFailedLoginsIpMap)
	failedLoginMutex.Lock()
	sleep_time = get_sleep_time(user_ip, sys_cfg.Limit.MaxFailedLoginSleepTimeSec)
	failedLoginMutex.Unlock()
	if sleep_time > 0 {
		time.Sleep(sleep_time)
	}
	return false, "", nil
}
