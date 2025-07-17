package subscription

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// SiteConfig 单个站点的配置
type SiteConfig struct {
	API    string `toml:"api" json:"api"`
	Name   string `toml:"name" json:"name"`
	Detail string `toml:"detail,omitempty" json:"detail,omitempty"`
	Adult  bool   `toml:"adult,omitempty" json:"adult,omitempty"`
}

// ApiSitesConfig 包含所有站点的配置
type ApiSitesConfig map[string]SiteConfig

// ApiSiteAndDefault 包含订阅站点配置和默认选择的站点列表
type ApiSiteAndDefault struct {
	ApiSites               ApiSitesConfig
	DefaultSelectedAPISite []string
}

// CurrentSubscriptionConfig 存储当前加载的统一订阅配置，全局可访问
var CurrentSubscriptionConfig atomic.Value

func init() {
	// 初始化 CurrentSubscriptionConfig 为一个空的 UnifiedSubscriptionConfig
	CurrentSubscriptionConfig.Store(ApiSiteAndDefault{
		ApiSites: make(ApiSitesConfig),
	})
}

// ReadSubscriptionConfigFromFile 从本地 TOML 文件读取订阅配置。
func ReadSubscriptionConfigFromFile(logger *zap.SugaredLogger, filePath string) (ApiSitesConfig, error) {
	conf := make(ApiSitesConfig)
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		logger.Warnf("load from %s get subscription config failed: %v", filePath, err)
		return nil, err
	}

	err := v.Unmarshal(&conf)
	if err != nil {
		logger.Errorf("Unmarshal file %s subscription config failed: %v", filePath, err)
		return nil, err
	}

	logger.Debugf("Loaded %s successfully, get %d sites.", filePath, len(conf))
	return conf, nil
}

// SaveSubscriptionConfigToFile 将 ApiSitesConfig 写入到指定的 TOML 文件。

func SaveSubscriptionConfigToFile(logger *zap.SugaredLogger, cfg ApiSitesConfig, filePath string) error {
	vp := viper.New() // 使用 vp 变量名避免与循环变量 v 混淆
	vp.SetConfigType("toml")

	// 遍历配置，将每个站点配置设置到 viper 实例中
	for k, siteConfig := range cfg {
		vp.Set(k, siteConfig)
	}

	err := vp.WriteConfigAs(filePath)
	if err != nil {
		logger.Errorf("将订阅配置写入文件 %s 失败: %v", filePath, err)
		return err
	}
	logger.Infof("订阅配置已成功写入文件: %s", filePath)
	return nil
}

// FetchAndMergeSubscriptions 从给定的 URL 列表获取 TOML 配置并合并。
func FetchAndMergeSubscriptions(githubDownloadMirror string, logger *zap.SugaredLogger, urls []string, subscriptionFilePath string) (ApiSitesConfig, error) {
	mergedConfig := make(ApiSitesConfig)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			if u == "" {
				logger.Errorf("订阅源 URL 为空")
				return
			}
			if len(githubDownloadMirror) > len("https://") {
				if !strings.HasSuffix(githubDownloadMirror, "/") {
					githubDownloadMirror += "/"
				}
				if strings.Contains(u, "github.com/") || strings.Contains(u, "raw.githubusercontent.com/") {
					u = githubDownloadMirror + u
				}
			}
			logger.Debugf("开始获取订阅源: %s", u)
			resp, err := http.Get(u)
			if err != nil {
				logger.Errorf("获取订阅源 %s 失败: %v", u, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Errorf("获取订阅源 %s 返回非正常状态码: %d", u, resp.StatusCode)
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("读取订阅源 %s 响应体失败: %v", u, err)
				return
			}

			// 使用 viper 解析 TOML
			v := viper.New()
			v.SetConfigType("toml") // 明确设置配置类型为 TOML

			// 将字节数据读入 viper
			err = v.ReadConfig(bytes.NewReader(data))
			if err != nil {
				logger.Errorf("解析订阅源 %s TOML 失败: %v", u, err)
				return
			}

			var tempConfig ApiSitesConfig
			// 将 viper 配置反序列化到结构体
			err = v.Unmarshal(&tempConfig)
			if err != nil {
				logger.Errorf("Unmarshal 订阅源 %s 配置失败: %v", u, err)
				return
			}

			mu.Lock()
			for k, v := range tempConfig {
				mergedConfig[k] = v
			}
			mu.Unlock()
			logger.Info("订阅源 %s 获取并合并成功。", u)
		}(url)
	}
	wg.Wait()

	// 按键名排序
	sortedKeys := make([]string, 0, len(mergedConfig))
	for k := range mergedConfig {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	sortedMergedConfig := make(ApiSitesConfig)
	for _, k := range sortedKeys {
		sortedMergedConfig[k] = mergedConfig[k]
	}

	logger.Info("所有订阅源合并成功，共 %d 个站点。", len(sortedMergedConfig))

	// 将合并后的配置写入本地文件
	if subscriptionFilePath != "" {
		if err := SaveSubscriptionConfigToFile(logger, sortedMergedConfig, subscriptionFilePath); err != nil {
			logger.Errorf("保存合并后的订阅配置到文件 %s 失败: %v", subscriptionFilePath, err)
		}
	} else {
		logger.Errorf("未指定订阅配置文件路径")
	}

	return sortedMergedConfig, nil
}

// ToJavaScriptCode 将 UnifiedSubscriptionConfig 转换为 JavaScript 格式的字符串。
func ToJavaScriptCode(config ApiSiteAndDefault) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("const API_SITES = {")

	// 创建一个 map 用于快速查找默认选中的站点
	selectedSiteMap := make(map[string]bool)
	for _, site := range config.DefaultSelectedAPISite {
		selectedSiteMap[site] = true
	}

	// 存放最终排序的键
	var orderedKeys []string
	// 存放非默认选中的站点键
	var remainingKeys []string

	// 优先添加默认选中的站点
	for _, key := range config.DefaultSelectedAPISite {
		if _, ok := config.ApiSites[key]; ok { // 确保站点存在
			orderedKeys = append(orderedKeys, key)
		}
	}

	// 添加剩余的站点，并进行字母排序
	for key := range config.ApiSites {
		if !selectedSiteMap[key] {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Strings(remainingKeys) // 对剩余站点进行字母排序

	// 合并最终的键列表
	allKeys := append(orderedKeys, remainingKeys...)

	for i, key := range allKeys {
		site := config.ApiSites[key]
		jsonBytes, err := json.MarshalIndent(site, "  ", "    ") // 使用 MarshalIndent 格式化 JSON
		if err != nil {
			return "", fmt.Errorf("无法 Marshal SiteConfig 到 JSON: %v", err)
		}
		buffer.WriteString(fmt.Sprintf("  %s: %s", key, string(jsonBytes)))
		if i < len(allKeys)-1 {
			buffer.WriteString(",\n")
		} else {
			buffer.WriteString("\n")
		}
	}
	buffer.WriteString("};\n")

	// 添加 DefaultSelectedAPISite
	selectedSitesJson, err := json.Marshal(config.DefaultSelectedAPISite)
	if err != nil {
		return "", fmt.Errorf("无法 Marshal DefaultSelectedAPISite 到 JSON: %v", err)
	}
	buffer.WriteString(fmt.Sprintf("const DefaultSelectedAPISite = %s;\n", string(selectedSitesJson)))

	return buffer.String(), nil
}
