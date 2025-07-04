package sub_reverseproxy

import (
	"context" // Needed for DialContext
	"net"     // Needed for DialContext
	"net/http"
	"net/http/httputil"
	"strings"
	"time" // For IdleConnTimeout

	"github.com/nas-core/nascore/nascore_handler_http/index_and_favicon"
	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func SubNasCoreVodSoket(subPathPrefix string, unixSocketPath *string, cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	// 创建一个自定义的 HTTP 传输层，用于连接到 Unix Socket
	unixSocketTransport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", *unixSocketPath)
		},

		DisableKeepAlives: true,
		IdleConnTimeout:   30 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			originalPath := req.URL.Path                                  // 获取原始请求路径
			targetPath := strings.TrimPrefix(originalPath, subPathPrefix) // 移除前缀，得到后端服务的实际路径
			req.URL.Scheme = "http"
			req.URL.Host = "unixsocket"
			req.URL.Path = "/" + targetPath // 后端服务的实际请求路径e
			req.Host = req.URL.Host
			req.Header.Set("Connection", "close")
		},
		Transport: unixSocketTransport,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Errorf("NasCoreVod proxy backend error: %v", err)
			englishMsg := "NasCoreVod might not be running. Please start it or enable auto-start. If already running, ensure Nascore server/container/VM can access the backend URL."
			chineseMsg := "NasCoreVod 可能未启动。请手动启动或开启随启动。如果已启动，请确保 Nascore 服务器/容器/虚拟机可访问后端地址。"
			// 在错误页面中包含后端 Unix Socket 路径，帮助诊断问题
			// Include the backend Unix Socket path in the error page to help diagnose the issue
			index_and_favicon.RenderPage(w, "NasCoreVod Backend Error 502", englishMsg+"     --------------------ErrInfo--------------------  "+*unixSocketPath+" ------------ "+err.Error(), chineseMsg)
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
