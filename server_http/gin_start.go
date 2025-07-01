package server_http

import (
	"nascore_v3/handlers_http/api"
	"nascore_v3/handlers_http/default_staticfileserver"
	"nascore_v3/handlers_http/index_and_favicon"
	"nascore_v3/handlers_http/subReverseproxy"
	"nascore_v3/handlers_http/webdav_core"
	"net/http"
	"strconv"

	"github.com/nas-core/nascore/nascore_util/followStartAndCron"

	webUI_ssi "github.com/nas-core/webui"

	"github.com/nas-core/nascore/pkgs/isDevMode"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinNasCoreStart(nsCfg *system_config.SysCfg, configFilePath string, logger *zap.SugaredLogger, qpsCounter *uint64) {
	go followStartAndCron.FollowStartAndCronMainLoop_forMachine(nsCfg, logger) // 防止独立部署的情况 没有请求的时候 自动任务不执行
	if !isDevMode.IsDevMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// ping  ==========
	r.GET("/@ping", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// favicon  ==========
	r.GET("/favicon.ico", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		index_and_favicon.Favicon_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// 静态文件  ==========
	r.GET(nsCfg.Server.DefaultStaticFileServicePrefix+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		default_staticfileserver.Default_staticfileserver_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// webui  ==========

	r.GET(nsCfg.Server.WebUIPrefix+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		webUI_ssi.Webui_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// api  ==========
	r.POST(system_config.PrefixApi+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		api.Api_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	r.GET(system_config.PrefixApi+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		api.Api_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	r.OPTIONS(system_config.PrefixApi+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		api.Api_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// DDNS-go
	r.Any(system_config.PrefixDdnsGo+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		subReverseproxy.SubDDnsGO(system_config.PrefixDdnsGo, &nsCfg.ThirdPartyExt.DdnsGO.ReverseproxyUrl, nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// adguardhome
	r.Any(system_config.PrefixAdguardhome+"*filepath", func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		subReverseproxy.SubAdguardhome(system_config.PrefixAdguardhome, &nsCfg.ThirdPartyExt.AdGuard.ReverseproxyUrl, nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// 欢迎页  ==========
	/*
	   r.GET("/", func(c *gin.Context) {
	   		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
	   		index_and_favicon.HanderWellcome(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	   	})
	*/
	// 未匹配的 就到webdav  ==========
	r.NoRoute(func(c *gin.Context) {
		followStartAndCron.FollowStartAndCronMain_forStateless(nsCfg, logger)
		webdav_core.Webdav_core_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	r.Run(":" + strconv.Itoa(nsCfg.Server.HttpPort))
}
