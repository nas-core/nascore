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

	// 这些变量用于跟踪在当前进程生命周期中是否已启动随从启动操作。在无服务器环境中，
	// 每个请求可能会启动一个新进程，因此理想情况下，如果外部程序需要在每个新实例上启动， 则这些变量应在每个请求时重置或重新评估。

	isRcloneMountFollowStart int32
	isDdnsSGOFollowStart     int32
	isCaddy2FollowStart      int32
	isOpenlistFollowStart    int32
	isExtProgramFollowStart  int32
)

// 防止独立部署的情况下，没有请求时自动任务不执行。
func FollowStartAndCronMainLoop_forMachine(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	for {
		if !nsCfg.Server.IsRunInServerLess { // 仅在 服务器部署的模式下运行此循环。		// 在无服务器模式下，每个请求会触发无状态检查。
			FollowStartAndCronMain_forStateless_andForMachine(nsCfg, logger)
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

// 插入到每个路由前面，方便兼容无状态服务器。
func FollowStartAndCronMain_forStateless_andForMachine(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	// 如果在无服务器模式下，重置随从启动标志以确保外部程序在每个新实例上被检查/启动。
	if nsCfg.Server.IsRunInServerLess {
		atomic.StoreInt32(&isRcloneMountFollowStart, 0)
		atomic.StoreInt32(&isDdnsSGOFollowStart, 0)
		atomic.StoreInt32(&isCaddy2FollowStart, 0)
		atomic.StoreInt32(&isOpenlistFollowStart, 0)
		atomic.StoreInt32(&isExtProgramFollowStart, 0)
		atomic.StoreInt32(&isLoopOneSecondrun, 0) // 确保 loopCheckFollowStart 每请求运行一次
		atomic.StoreInt32(&isCheckingCron, 0)     // 确保 cronFunc 每请求运行一次
	}
	if nsCfg.Server.IsRunInServerLess {
		CheckAllExtensionStatusOnce()
	}
	if atomic.LoadInt32(&isLoopOneSecondrun) == 0 { // 避免循环启动
		if nsCfg.Server.IsRunInServerLess {
			loopCheckFollowStart(nsCfg, logger) // 同步执行，无睡眠
		} else {
			go loopCheckFollowStart(nsCfg, logger) // 异步执行，含睡眠循环
		}
	}
	if atomic.LoadInt32(&isCheckingCron) == 0 { // 如果已经在检查 那么不检查
		if nsCfg.Server.IsRunInServerLess {
			cronFunc(nsCfg, logger) // 同步执行，无睡眠
		} else {
			go cronFunc(nsCfg, logger) // 异步执行，含睡眠循环
		}
	}
}

/**
 * 计划任务主函数
 */
func cronFunc(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	// 防止 cronFunc 在同一进程内并发执行。 在无服务器模式下，确保它每个请求运行一次。 在非无服务器模式下，防止多个 cron goroutine 运行。
	if !atomic.CompareAndSwapInt32(&isCheckingCron, 0, 1) {
		return // 另一个 cronFunc 正在运行或被检查
	}
	defer atomic.StoreInt32(&isCheckingCron, 0)

	nowTimeInt64 := time.Now().Unix() // 秒

	if nsCfg.ThirdPartyExt.AdGuard.AutoUpdateRulesEnable {
		if nowTimeInt64-lastExecADGuardsGetRulesTime > int64(nsCfg.ThirdPartyExt.AdGuard.AutoUpdateRulesInterval*3600) {
			logger.Info("Start ADGuards update Rules as scheduled.")
			execADGuardsGetRules(nsCfg, logger)
			atomic.StoreInt64(&lastExecADGuardsGetRulesTime, nowTimeInt64)
		}
	}

	if nsCfg.ThirdPartyExt.AcmeLego.IsLegoAutoRenew {
		if nowTimeInt64-lastExecLegoRenewOrGetTime > int64(nsCfg.ThirdPartyExt.AcmeLego.AutoUpdateCheckInterval*3600) {
			logger.Info("Start lego as scheduled.")
			execLegoRenewOrGet(nsCfg, logger)
			atomic.StoreInt64(&lastExecLegoRenewOrGetTime, nowTimeInt64)
		}
	}

	// 热重载配置。仅在未重新加载时尝试。
	if atomic.CompareAndSwapInt32(&isReloadingNascoreToml, 0, 1) {
		reloadNascoreToml(nsCfg)
		atomic.StoreInt32(&isReloadingNascoreToml, 0)
	}
}

func loopCheckFollowStart(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	// 防止 loopCheckFollowStart 在同一进程内并发执行。
	if !atomic.CompareAndSwapInt32(&isLoopOneSecondrun, 0, 1) {
		return // 另一个 loopCheckFollowStart 正在运行或被检查
	}
	defer atomic.StoreInt32(&isLoopOneSecondrun, 0) // 在无服务器模式下，运行一次后重置；在非无服务器模式下，由于无限循环，此 defer 是多余的。

	// 在非无服务器模式下，此循环持续运行。
	// 在无服务器模式下，它每个请求周期运行一次。
	for {
		if nsCfg.ThirdPartyExt.Rclone.AutoMountEnable && atomic.CompareAndSwapInt32(&isRcloneMountFollowStart, 0, 1) {
			logger.Info("Starting Rclone AutoMount FollowStart.")
			if nsCfg.Server.IsRunInServerLess {
				RcloneFollowStart(nsCfg, logger)
			} else {
				go RcloneFollowStart(nsCfg, logger)
			}
		}
		if nsCfg.ThirdPartyExt.DdnsGO.AutoStartEnable && atomic.CompareAndSwapInt32(&isDdnsSGOFollowStart, 0, 1) {
			logger.Info("Starting DdnsSGOFollowStart.")
			if nsCfg.Server.IsRunInServerLess {
				DdnsSGOFollowStart(nsCfg, logger)
			} else {
				go DdnsSGOFollowStart(nsCfg, logger)
			}
		}
		// caddy2
		if nsCfg.ThirdPartyExt.Caddy2.AutoStartEnable && atomic.CompareAndSwapInt32(&isCaddy2FollowStart, 0, 1) {
			logger.Info("Starting Caddy2FollowStart.")
			if nsCfg.Server.IsRunInServerLess {
				Caddy2FollowStart(nsCfg, logger)
			} else {
				go Caddy2FollowStart(nsCfg, logger)
			}
		}
		// openlost
		if nsCfg.ThirdPartyExt.Openlist.AutoStartEnable && atomic.CompareAndSwapInt32(&isOpenlistFollowStart, 0, 1) {
			logger.Info("Starting OpenlistFollowStart.")
			if nsCfg.Server.IsRunInServerLess {
				OpenlistFollowStart(nsCfg, logger)
			} else {
				go OpenlistFollowStart(nsCfg, logger)
			}
		}
		// 扩展程序
		if atomic.CompareAndSwapInt32(&isExtProgramFollowStart, 0, 1) {
			logger.Info("Starting nascore_vod")
			if nsCfg.Server.IsRunInServerLess {
				Nascore_extended_followStart(nsCfg, logger)
			} else {
				go Nascore_extended_followStart(nsCfg, logger)
			}
		}
		if !nsCfg.Server.IsRunInServerLess {
			CheckAllExtensionStatusOnce()
		}
		if nsCfg.Server.IsRunInServerLess {
			break // 在无服务器模式下，我们不希望睡眠并阻塞请求。			// 它应该运行一次检查并返回。
		} else {
			// 在非无服务器模式下，每秒循环一次。
			time.Sleep(time.Second)
		}
	}
}
