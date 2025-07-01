package index_and_favicon

import (
	"embed" // 引入 embed 包
	"fmt"
	"html/template" // 引入 html/template 包
	"net/http"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"go.uber.org/zap"
)

//go:embed welcome_page.html
var welcomePageContent embed.FS // 声明一个 embed.FS 变量来嵌入 HTML 文件

var welcomePageTemplate *template.Template

// init 函数在包加载时执行，用于解析 HTML 模板
func init() {
	var err error
	welcomePageTemplate, err = template.ParseFS(welcomePageContent, "welcome_page.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse welcome_page.html template: %v", err))
	}
}

// WelcomePageData 定义了传递给欢迎页面的数据结构
type WelcomePageData struct {
	WebUIPrefix string
}

func HanderWellcome(cfg *system_config.SysCfg, logger *zap.SugaredLogger, qpsCounter *uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 准备数据
		data := WelcomePageData{
			WebUIPrefix: cfg.Server.WebUIPrefix,
		}

		// 执行模板并写入响应
		err := welcomePageTemplate.Execute(w, data)
		if err != nil {
			logger.Errorf("Failed to render welcome_page.html: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
