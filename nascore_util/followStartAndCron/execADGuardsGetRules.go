package followStartAndCron

import (
	"path/filepath"
	"strings"

	"github.com/nas-core/nascore/nascore_util/downfile"
	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

func execADGuardsGetRules(nsCfg *system_config.SysCfg, logger *zap.SugaredLogger) {
	err := DownloadADGuardRules(&nsCfg.ThirdPartyExt.AdGuard.Upstream_dns_fileUpdateUrl, &nsCfg.ThirdPartyExt.GitHubDownloadMirror, &nsCfg.ThirdPartyExt.AdGuard.Upstream_dns_file)
	if err != nil {
		logger.Errorw("Download ADGuard rules failed", "error", err)
	}
}

func DownloadADGuardRules(Upstream_dns_fileUpdateUrl *string, GitHubDownloadMirror *string, Upstream_dns_file *string) error {
	DownLoadlink := *Upstream_dns_fileUpdateUrl

	if len(*GitHubDownloadMirror) > len("https://") {
		if !strings.HasSuffix(*GitHubDownloadMirror, "/") {
			*GitHubDownloadMirror += "/"
		}
		if strings.Contains(DownLoadlink, "github.com/") || strings.Contains(DownLoadlink, "raw.githubusercontent.com/") {
			DownLoadlink = *GitHubDownloadMirror + DownLoadlink
		}
	}
	saveFilename := filepath.Base(*Upstream_dns_file)
	SaveDir := filepath.Dir(*Upstream_dns_file)

	err := downfile.DownloadFile(DownLoadlink, SaveDir, saveFilename)
	if err != nil {
		return err
	}
	return nil
}
