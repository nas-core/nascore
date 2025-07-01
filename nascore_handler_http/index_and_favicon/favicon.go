package index_and_favicon

import (
	"embed"
	"net/http"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

//go:embed favicon.ico
var faviconFS embed.FS

func Favicon_handler(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := faviconFS.ReadFile("favicon.ico")
		if err != nil {
			logger.Errorln("Error reading embedded file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(content)
	}
}
