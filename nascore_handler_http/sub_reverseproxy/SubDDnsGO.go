package sub_reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/nas-core/nascore/nascore_handler_http/index_and_favicon"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func SubDDnsGO(subPathPrefix string, backEndUrl *string, cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.ThirdPartyExt.DdnsGO.IsDDnsGOProxyEnable {
			index_and_favicon.RenderPage(w,
				"DDNSGO proxy is not enabled",
				"The reverse proxy function of DDNSGO is not enabled. Please enable it in the background or configuration file.",
				"DDNSGO 反向代理功能没有启用。请到后台或者配置文件中启用。",
				"system.shtml#ThirdPartyExtDdnsGO", "Goto",
			)
			return
		}
		originalPath := r.URL.Path                                    // 解析目标 URL
		targetPath := strings.TrimPrefix(originalPath, subPathPrefix) // 移除前缀

		backendfullURL, err := url.Parse(*backEndUrl) // 拼接目标 URL
		if err != nil {
			logger.Errorf("get backenURL Parse err: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		backendfullURL.Path = "/" + targetPath
		logger.Debug("backendfullURL  ", backendfullURL)
		proxy := httputil.ReverseProxy{ // 创建反向代理请求
			Director: func(req *http.Request) {
				req.URL = backendfullURL
				req.Host = backendfullURL.Host
			},
			ModifyResponse: func(resp *http.Response) error {
				if resp.StatusCode >= 300 && resp.StatusCode < 400 {
					if location, err := resp.Location(); err == nil && location.Path == "/login" {
						resp.Header.Set("Location", subPathPrefix+"login")
					}
				}
				return nil
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				logger.Errorf("DDNS-GO proxy backend error: %v", err)
				index_and_favicon.RenderPage(w, "DDNS-GO Backend Error 502", "DDNS-GO might not be running. Please start it or enable auto-start. If already running, ensure Nascore server/container/VM can access the backend URL.     --------------------ErrInfo--------------------  "+backendfullURL.String()+" ------------ "+err.Error(), "DDNS-GO 可能未启动。请手动启动或开启随启动。如果已启动，请确保 Nascore 服务器/容器/虚拟机可访问后端地址。", "system.shtml#ThirdPartyExtDdnsGO", "Goto")
			},
		}
		proxy.ServeHTTP(w, r)
	}
}
