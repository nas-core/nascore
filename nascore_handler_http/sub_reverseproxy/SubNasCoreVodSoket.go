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

func SubNasCoreVodSocket(subPathPrefix string, unixSocketPath *string, cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse("http://unix")
		if err != nil {
			logger.Errorw("failed to parse target URL", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		socketFilePathValue := *unixSocketPath
		if len(socketFilePathValue) > 0 && socketFilePathValue[len(socketFilePathValue)-1] != '/' {
			socketFilePathValue += "/"
		}
		socketFilePathValue += system_config.NasCoreTvSocketFile
		transport := &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				return net.Dial("unix", socketFilePathValue)
			},
			ResponseHeaderTimeout: 60 * time.Second, // 设置等待后端响应头的超时时间
			IdleConnTimeout:       90 * time.Second, // 设置连接保持空闲的最长时间
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		// 如果 r.URL.Path 里面不包含 字符串 proxy/
		// 特指 类似 http://localhost:9000/@nascore_vod/proxy/https%3A%2F%2Fmovie.douban.com%2Fj%2Fsearch_subjects%3Ftype%3Dmovie%26tag%3D%E7%83%AD%E9%97%A8%26sort%3Drecommend%26page_limit%3D16%26page_start%3D0
		/*
		   if !strings.Contains(r.URL.Path, "proxy/") {
		   			r.URL.Path = strings.TrimPrefix(r.URL.Path, subPathPrefix)
		   			if r.URL.Path == "" {
		   				r.URL.Path = "/"
		   			}
		   		}
		*/

		r.URL.Path = strings.TrimPrefix(r.URL.Path, subPathPrefix)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}

		proxy.ServeHTTP(w, r)
	}
}
