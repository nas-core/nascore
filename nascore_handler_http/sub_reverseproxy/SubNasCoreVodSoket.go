package sub_reverseproxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/nas-core/nascore/nascore_auth/user/user_helper"
	"github.com/nas-core/nascore/nascore_handler_http/index_and_favicon"
	"github.com/nas-core/nascore/nascore_util/system_config"
	"go.uber.org/zap"
)

func SubNasCoreVodSocket(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	// logger.Info("SubNasCoreVodSocket started")

	return func(w http.ResponseWriter, r *http.Request) {
		// logger.Info("SubNasCoreVodSocket started url path ", r.URL.Path)
		if r.URL.Path == "admin_setting.html" {
			userInfo, err := user_helper.ValidateTokenAndGetUserInfo(r, nsCfg)
			if err != nil {
				logger.Errorw("token get err", "error", err) // [auth]
				index_and_favicon.RenderPage(w, "Token validation failed ", "Token validation failed ", "无法验证您的身份，请重新登录。")
				return
			}
			if !userInfo.IsAdmin {
				logger.Warnw("用户不是管理员", "user", userInfo.Username) // [auth]
				index_and_favicon.RenderPage(w, "Insufficient permissions", "Insufficient permissions", "您不是管理员，无法访问此页面。")
				return
			}
		}

		socketFilePathValue := nsCfg.Server.UnixSocketFilePath
		if len(socketFilePathValue) > 0 && socketFilePathValue[len(socketFilePathValue)-1] != '/' {
			socketFilePathValue += "/"
		}
		socketFilePathValue += system_config.NasCoreTvSocketFile

		// 检查 Socket 文件是否存在
		if _, err := os.Stat(socketFilePathValue); os.IsNotExist(err) {
			logger.Errorw("socket 文件不存在", "path", socketFilePathValue, "error", err) // [socket]
			index_and_favicon.RenderPage(w, "服务不可用", "Service Unavailable", fmt.Sprintf("后端服务 socket 文件 %s 不存在。", socketFilePathValue))
			return
		}

		target, err := url.Parse("http://unix")
		if err != nil {
			logger.Errorw("failed to parse target URL", "error", err) // [proxy]
			index_and_favicon.RenderPage(w, "502 错误", "Bad Gateway", "无法连接到后端服务。")
			return
		}

		transport := &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				return net.Dial("unix", socketFilePathValue)
			},
			ResponseHeaderTimeout: 60 * time.Second, // 设置等待后端响应头的超时时间
			IdleConnTimeout:       90 * time.Second, // 设置连接保持空闲的最长时间
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Errorw("反向代理出错", "error", err) // [proxy]
			index_and_favicon.RenderPage(w, "502 错误", "Bad Gateway", "后端服务无响应，可能没有安装nascore_vod扩展 没启动。")
		}

		//	r.URL.Path = strings.TrimPrefix(r.URL.Path, system_config.PrefixNasCoreTv)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}

		proxy.ServeHTTP(w, r)
	}
}
