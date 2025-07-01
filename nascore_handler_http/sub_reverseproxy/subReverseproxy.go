package sub_reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func SubReverseproxy(subPathPrefix string, backEndUrl *string, cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		}
		proxy.ServeHTTP(w, r)
	}
}
