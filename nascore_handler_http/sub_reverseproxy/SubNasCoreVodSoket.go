package sub_reverseproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/joyanhui/golang-pkgs/pkgs/response_yh"
	"github.com/nas-core/nascore/nascore_util/system_config"
	"go.uber.org/zap"
)

func SubNasCoreVodSocket(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	// logger.Info("SubNasCoreVodSocket started")

	return func(w http.ResponseWriter, r *http.Request) {
		// logger.Info("SubNasCoreVodSocket started url path ", r.URL.Path)
		if r.URL.Path == "admin_setting.html" {
			var err error
			userInfo, err := user_helper.ValidateTokenAndGetUserInfo(r, nsCfg)
			if err != nil {
				response_yh.SendError(w, "Authorization failed: "+err.Error(), http.StatusUnauthorized)
				return
			}
			if !userInfo.IsAdmin { // 仅限管理员
				response_yh.SendError(w, "you are not admin ", http.StatusUnauthorized)
				return
			}
		}
		target, err := url.Parse("http://unix")
		if err != nil {
			logger.Errorw("failed to parse target URL", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		socketFilePathValue := nsCfg.Server.UnixSocketFilePath
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

		//	r.URL.Path = strings.TrimPrefix(r.URL.Path, system_config.PrefixNasCoreTv)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}

		proxy.ServeHTTP(w, r)
	}
}
