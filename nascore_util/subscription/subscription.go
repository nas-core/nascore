package subscription

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// SiteConfig 单个站点的配置
type SiteConfig struct {
	API    string `toml:"api" json:"api"`
	Name   string `toml:"name" json:"name"`
	Detail string `toml:"detail,omitempty" json:"detail,omitempty"`
	Adult  bool   `toml:"adult,omitempty" json:"adult,omitempty"`
	Hidden bool   `toml:"hidden,omitempty" json:"hidden,omitempty"`
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

// FetchAndMergeSubscriptions 从给定的 URL 列表获取 TOML 配置并合并。
func FetchAndMergeSubscriptions(githubDownloadMirror string, logger *zap.SugaredLogger, urls []string) (ApiSitesConfig, error) {
	mergedConfig := make(ApiSitesConfig)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			if u == "" {
				logger.Errorf("[subscription] Subscription source URL is empty")
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
			logger.Debug("[subscription] Start fetching subscription source: %s", u)
			resp, err := http.Get(u)
			if err != nil {
				logger.Errorf("[subscription] Failed to fetch subscription source %s: %v", u, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Errorf("[subscription] Subscription source %s returned non-OK status: %d", u, resp.StatusCode)
				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("[subscription] Failed to read response body from %s: %v", u, err)
				return
			}

			v := viper.New()
			v.SetConfigType("toml")

			err = v.ReadConfig(bytes.NewReader(data))
			if err != nil {
				logger.Errorf("[subscription] Failed to parse TOML from %s: %v", u, err)
				return
			}

			var tempConfig ApiSitesConfig
			err = v.Unmarshal(&tempConfig)
			if err != nil {
				logger.Errorf("[subscription] Failed to unmarshal config from %s: %v", u, err)
				return
			}

			mu.Lock()
			for k, v := range tempConfig {
				mergedConfig[k] = v
			}
			mu.Unlock()
			logger.Debug("[subscription] Subscription source %s fetched and merged successfully.", u)
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

	logger.Debug("[subscription] All subscription sources merged, total %d sites.", len(sortedMergedConfig))

	return sortedMergedConfig, nil
}

// ReadSubscriptionConfigFromString 从TOML字符串读取订阅配置。
func ReadSubscriptionConfigFromString(logger *zap.SugaredLogger, tomlStr string) (ApiSitesConfig, error) {
	conf := make(ApiSitesConfig)
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewBufferString(tomlStr)); err != nil {
		logger.Warnf("[subscription] load from string get subscription config failed: %v", err)
		return nil, err
	}
	if err := v.Unmarshal(&conf); err != nil {
		logger.Errorf("[subscription] Unmarshal string subscription config failed: %v", err)
		return nil, err
	}
	logger.Debug("[subscription] Loaded from string successfully, got %d sites.", len(conf))
	return conf, nil
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

	// 优先添加默认选中的站点，过滤掉hidden=true的站点
	for _, key := range config.DefaultSelectedAPISite {
		if site, ok := config.ApiSites[key]; ok && !site.Hidden { // 确保站点存在且不是隐藏的
			orderedKeys = append(orderedKeys, key)
		}
	}

	// 添加剩余的站点，并进行字母排序，过滤掉hidden=true的站点
	for key, site := range config.ApiSites {
		if !selectedSiteMap[key] && !site.Hidden {
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

	// 添加 DefaultSelectedAPISite，过滤掉隐藏的站点
	var filteredSelectedSites []string
	for _, site := range config.DefaultSelectedAPISite {
		if siteConfig, ok := config.ApiSites[site]; ok && !siteConfig.Hidden {
			filteredSelectedSites = append(filteredSelectedSites, site)
		}
	}

	selectedSitesJson, err := json.Marshal(filteredSelectedSites)
	if err != nil {
		return "", fmt.Errorf("无法 Marshal DefaultSelectedAPISite 到 JSON: %v", err)
	}
	buffer.WriteString(fmt.Sprintf("const DefaultSelectedAPISite = %s;\n", string(selectedSitesJson)))

	return buffer.String(), nil
}

// SaveSubscriptionToDB 保存TOML字符串到vod_subscription表
func SaveSubscriptionToDB(db *sql.DB, tomlStr string) error {
	now := time.Now().Unix()
	_, err := db.Exec(`INSERT INTO vod_subscription (id, data, lastlogin_at) VALUES (1, ?, ?) ON CONFLICT(id) DO UPDATE SET data=excluded.data, lastlogin_at=excluded.lastlogin_at`, tomlStr, now)
	return err
}

// LoadSubscriptionFromDB 从vod_subscription表加载并解析
func LoadSubscriptionFromDB(db *sql.DB, logger *zap.SugaredLogger) (ApiSitesConfig, error) {
	row := db.QueryRow(`SELECT data FROM vod_subscription WHERE id=1`)
	var tomlStr string
	if err := row.Scan(&tomlStr); err != nil {
		return nil, err
	}
	return ReadSubscriptionConfigFromString(logger, tomlStr)
}

// MergeRemoteSubscriptions 拉取并合并远程订阅，返回结构和TOML字符串
func MergeRemoteSubscriptions(urls []string, mirror string, logger *zap.SugaredLogger) (ApiSitesConfig, string, error) {
	merged, err := FetchAndMergeSubscriptions(mirror, logger, urls)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(merged); err != nil {
		return nil, "", err
	}
	return merged, buf.String(), nil
}

// RefreshSubscriptionAndSaveToDB 拉取合并并写入DB
func RefreshSubscriptionAndSaveToDB(db *sql.DB, urls []string, mirror string, logger *zap.SugaredLogger) error {
	_, tomlStr, err := MergeRemoteSubscriptions(urls, mirror, logger)
	if err != nil {
		return err
	}
	return SaveSubscriptionToDB(db, tomlStr)
}
