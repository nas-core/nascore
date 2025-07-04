package sub_reverseproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/nas-core/nascore/nascore_util/system_config"
	"go.uber.org/zap"
)

func SubNasCoreVodSoket(subPathPrefix string, unixSocketPath *string, cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse("http://unix")
		if err != nil {
			logger.Errorw("failed to parse target URL", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		transport := &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				return net.Dial("unix", *unixSocketPath)
			},
			ResponseHeaderTimeout: 60 * time.Second, // 设置等待后端响应头的超时时间
			IdleConnTimeout:       90 * time.Second, // 设置连接保持空闲的最长时间
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		// 如果 r.URL.Path 里面不包含 字符串 proxy/
		if !strings.Contains(r.URL.Path, "proxy/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, subPathPrefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}
		proxy.ServeHTTP(w, r)
	}
}
