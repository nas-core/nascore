package followStartAndCron

import (
	"sync/atomic"
	"time"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

var (
	isLoopOneSecondrun           int32
	isCheckingCron               int32
	lastExecADGuardsGetRulesTime int64
	lastExecLegoRenewOrGetTime   int64
	isReloadingNascoreToml       int32

	isRcloneMountFollowStart int32
	isDdnsSGOFollowStart     int32
	isCaddy2FollowStart      int32
	isOpenlistFollowStart    int32

	isExtProgramFollowStart int32
)

// 防止独立部署的情况 没有请求的时候 自动任务不执行
func FollowStartAndCronMainLoop_forMachine(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	for {
		FollowStartAndCronMain_forStateless(nsCfg, logger)
		time.Sleep(time.Second)
	}
}

// 插入到每一个路由前面 方便兼容无状态服务器
func FollowStartAndCronMain_forStateless(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	if atomic.LoadInt32(&isLoopOneSecondrun) == 0 { // 避免循环启动
		if nsCfg.Server.IsRunInServerLess {
			loopCheckFollowStart(nsCfg, logger)
		} else {
			go loopCheckFollowStart(nsCfg, logger)
		}
	}
	if atomic.LoadInt32(&isCheckingCron) == 0 { // 如果已经在检查 那么不检查
		if nsCfg.Server.IsRunInServerLess {
			cronFunc(nsCfg, logger)
		} else {
			go cronFunc(nsCfg, logger)
		}
	}
}

/**
 * 计划任务主函数
 */
func cronFunc(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	atomic.StoreInt32(&isCheckingCron, 1)
	nowTimeInt64 := time.Now().Unix() // 秒
	if nsCfg.ThirdPartyExt.AdGuard.AutoUpdateRulesEnable {
		if nowTimeInt64-lastExecADGuardsGetRulesTime > int64(nsCfg.ThirdPartyExt.AdGuard.AutoUpdateRulesInterval*3600) {
			execADGuardsGetRules(nsCfg, logger)
			atomic.StoreInt64(&lastExecADGuardsGetRulesTime, nowTimeInt64)
		}
	}

	if nsCfg.ThirdPartyExt.AcmeLego.IsLegoAutoRenew {
		if nowTimeInt64-lastExecLegoRenewOrGetTime > int64(nsCfg.ThirdPartyExt.AcmeLego.AutoUpdateCheckInterval*3600) {
			execLegoRenewOrGet(nsCfg, logger)
			atomic.StoreInt64(&lastExecLegoRenewOrGetTime, nowTimeInt64)
		}
	}
	if atomic.LoadInt32(&isReloadingNascoreToml) == 0 {
		atomic.StoreInt32(&isReloadingNascoreToml, 1)
		reloadNascoreToml(nsCfg)
		atomic.StoreInt32(&isReloadingNascoreToml, 0)
	}
	atomic.StoreInt32(&isCheckingCron, 0)
}

func loopCheckFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	for {
		atomic.StoreInt32(&isLoopOneSecondrun, 1)
		if nsCfg.ThirdPartyExt.Rclone.AutoMountEnable && atomic.LoadInt32(&isRcloneMountFollowStart) == 0 {
			atomic.StoreInt32(&isRcloneMountFollowStart, 1)
			go RcloneFollowStart(nsCfg, logger)

		}
		if nsCfg.ThirdPartyExt.DdnsGO.AutoStartEnable && atomic.LoadInt32(&isDdnsSGOFollowStart) == 0 {
			atomic.StoreInt32(&isDdnsSGOFollowStart, 1)
			go DdnsSGOFollowStart(nsCfg, logger)
		}
		// caddy2
		if nsCfg.ThirdPartyExt.Caddy2.AutoStartEnable && atomic.LoadInt32(&isCaddy2FollowStart) == 0 {
			atomic.StoreInt32(&isCaddy2FollowStart, 1)
			go Caddy2FollowStart(nsCfg, logger)
		}
		// openlost
		if nsCfg.ThirdPartyExt.Openlist.AutoStartEnable && atomic.LoadInt32(&isOpenlistFollowStart) == 0 {
			atomic.StoreInt32(&isOpenlistFollowStart, 1)
			go OpenlistFollowStart(nsCfg, logger)
		}
		// 扩展程序
		if atomic.LoadInt32(&isExtProgramFollowStart) == 0 {
			atomic.StoreInt32(&isExtProgramFollowStart, 1)
			go Nascore_extended_followStart(nsCfg, logger)
		}
		time.Sleep(time.Second)
		atomic.StoreInt32(&isLoopOneSecondrun, 0)
	}
}
