package server_http

import (
	"log"
	"nascore_v3/handlers_http/api"
	"nascore_v3/handlers_http/default_staticfileserver"
	"nascore_v3/handlers_http/index_and_favicon"
	"nascore_v3/handlers_http/subReverseproxy"
	"nascore_v3/handlers_http/webUI_ssi"
	"nascore_v3/handlers_http/webdav_core"
	"nascore_v3/pkgs/ddnsgo"
	rclonefollowstart "nascore_v3/pkgs/rcloneExe/rcloneFollowStart"
	"nascore_v3/system_config"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinStart() {
	nsCfg, configFilePath, err := initFlagsAndSysCfg()
	if err != nil {
		log.Fatal("system init failed:", err)
		os.Exit(1)
	}
	system_config.ConfigFilePath = configFilePath
	logger := zap.NewExample().Sugar()
	var qpsCounter *uint64

	go rclonefollowstart.RcloneFollowStart(nsCfg, logger, qpsCounter)
	ddnsgo.DdnsSGOFollowStart(nsCfg, logger, qpsCounter)
	r := gin.Default()
	// ping  ==========
	r.GET("/@ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// favicon  ==========
	r.GET("/favicon.ico", func(c *gin.Context) {
		index_and_favicon.Favicon_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// 静态文件  ==========
	r.GET(nsCfg.Server.DefaultStaticFileServicePrefix+"*filepath", func(c *gin.Context) {
		default_staticfileserver.Default_staticfileserver_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// webui  ==========

	r.GET(nsCfg.Server.WebUIPrefix+"*filepath", func(c *gin.Context) {
		webUI_ssi.Webui_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// api  ==========
	r.POST(system_config.PrefixApi+"*filepath", func(c *gin.Context) {
		api.Api_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	r.GET(system_config.PrefixApi+"*filepath", func(c *gin.Context) {
		api.Api_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// DDNS-go
	r.Any(system_config.PrefixDdnsGo+"*filepath", func(c *gin.Context) {
		subReverseproxy.SubDDnsGO(system_config.PrefixDdnsGo, &nsCfg.ThirdPartyExt.DdnsGO.ReverseproxyUrl, nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// adguardhome
	r.Any(system_config.PrefixAdguardhome+"*filepath", func(c *gin.Context) {
		subReverseproxy.SubAdguardhome(system_config.PrefixAdguardhome, &nsCfg.ThirdPartyExt.AdGuard.ReverseproxyUrl, nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// 欢迎页  ==========
	r.GET("/", func(c *gin.Context) {
		index_and_favicon.HanderWellcome(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	// 未匹配的 就到webdav  ==========
	r.NoRoute(func(c *gin.Context) {
		webdav_core.Webdav_core_handler(nsCfg, logger, qpsCounter)(c.Writer, c.Request)
	})
	r.Run(":" + strconv.Itoa(nsCfg.Server.HttpPort))
}
