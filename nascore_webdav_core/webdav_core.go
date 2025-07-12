package nascore_webdav_core

import (
	"net/http"

	"github.com/nas-core/nascore/nascore_auth/user/checkpsw"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"github.com/emersion/go-webdav"
	"go.uber.org/zap"
)

func Webdav_core_handler(sys_cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// 检查WebDAV服务是否启用
		if !sys_cfg.Server.WebDavEnable {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("WebDAV服务未启用"))
			return
		}
		// 获取用户名/密码
		username, password, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		passwdOk, homepath, err := checkpsw.AuthUserAndGetUserInfo(req, logger, sys_cfg, username, password)
		if !passwdOk || err != nil {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		fs_webdav := &webdav.Handler{
			FileSystem: webdav.LocalFileSystem(homepath),
		}
		fs_webdav.ServeHTTP(w, req)
	}
}
