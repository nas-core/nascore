package default_staticfileserver

import (
	"net/http"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func Default_staticfileserver_handler(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !nsCfg.Server.DefaultStaticFileServiceEnable {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Default static file server is not enabled"))
			return
		}

		/*
			 * // 之前的 代码 		fs := http.StripPrefix(sys_cfg.Server.DefaultStaticFileServiceRoot, http.FileServer(http.Dir(sys_cfg.Server.DefaultStaticFileServiceRoot)))
				logger.Info("Default_staticfileserver_handler  sys_cfg.Server.DefaultStaticFileServicePrefix:", sys_cfg.Server.DefaultStaticFileServicePrefix)
				logger.Info("Default_staticfileserver_handler  sys_cfg.Server.DefaultStaticFileServiceRoot:", sys_cfg.Server.DefaultStaticFileServiceRoot)
		*/
		fs := http.StripPrefix(nsCfg.Server.DefaultStaticFileServicePrefix, http.FileServer(http.Dir(nsCfg.Server.DefaultStaticFileServiceRoot)))

		fs.ServeHTTP(w, r)
	}
}
