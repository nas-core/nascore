package index_and_favicon

import (
	"bytes"
	"embed" // 引入 embed 包
	"fmt"
	"html/template" // 引入 html/template 包
	"net/http"

	"github.com/nas-core/WebUi/pkgs/replacetemplateplaceholders"
)

//go:embed RenderPage.html
var content embed.FS // 声明一个 embed.FS 变量来嵌入 HTML 文件

var errorPageTemplate *template.Template

// init 函数在包加载时执行，用于解析 HTML 模板
func init() {
	var err error
	errorPageTemplate, err = template.ParseFS(content, "RenderPage.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse RenderPage.html template: %v", err))
	}
}

// errorPageData 定义了传递给错误页面的数据结构
type errorPageData struct {
	Title              string
	EnglishDescription string
	ChineseDescription string
	GotoLink           string
	GotoText           string
	WebUICdnPrefix     string
}

// RenderPage 生成并写入一个美观的 HTML 页面
// title：页面标题。
// englishDescription：英文错误描述。
// chineseDescription：中文错误描述。
// webUICdnPrefix：前端CDN前缀。
func RenderPage(w http.ResponseWriter, title, englishDescription, chineseDescription, GotoLink, GotoText, webUICdnPrefix string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := errorPageData{
		Title:              title,
		EnglishDescription: englishDescription,
		ChineseDescription: chineseDescription,
		GotoLink:           GotoLink,
		GotoText:           GotoText,
		WebUICdnPrefix:     webUICdnPrefix,
	}

	var tplBuf bytes.Buffer
	err := errorPageTemplate.Execute(&tplBuf, data)
	if err != nil {
		http.Error(w, "Internal Server Error by RenderPage", http.StatusInternalServerError)
		return
	}
	parsedContent := replacetemplateplaceholders.ReplaceTemplatePlaceholders(tplBuf.String(), webUICdnPrefix, "")
	w.Write([]byte(parsedContent))
}
