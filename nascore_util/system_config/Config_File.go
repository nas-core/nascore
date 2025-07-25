package system_config

import (
	"runtime"
)

type SysCfg struct {
	Server         ServerStru `mapstructure:"Server"`
	JWT            JwtStru    `mapstructure:"JWT"`
	Secret         SecretStru `mapstructure:"Secret"`
	WebUICdnPrefix string     `mapstructure:"WebUICdnPrefix"`
	//	Users          []map[string]string `mapstructure:"users"`
	//	WebSites       []WebsiteEntry      `mapstructure:"WebSites"`
	Limit LimitStru `mapstructure:"Limit"`

	NascoreExt    NascoreExtStru    `mapstructure:"NascoreExt"`
	ThirdPartyExt ThirdPartyExtStru `mapstructure:"ThirdPartyExt"`
}

type NascoreExtStru struct {
	UserID  string     `mapstructure:"UserID"`
	UserKey string     `mapstructure:"UserKey"`
	Vod     VodExtStru `mapstructure:"Vod"`
	Links   LinksStru  `mapstructure:"Links"`
}

type ThirdPartyExtStru struct {
	GitHubDownloadMirror string        `mapstructure:"GitHubDownloadMirror"`
	Rclone               RcloneExtStru `mapstructure:"Rclone"`
	DdnsGO               DdnsgoStru    `mapstructure:"DdnsGO"`
	AdGuard              AdGuardStru   `mapstructure:"AdGuard"`
	AcmeLego             AcmeLegoStru  `mapstructure:"AcmeLego"`
	Caddy2               Caddy2Stru    `mapstructure:"Caddy2"`
	Openlist             OpenlistStru  `mapstructure:"Openlist"`
}
type OpenlistStru struct {
	AutoStartEnable bool   `mapstructure:"AutoStartEnable"`
	Version         string `mapstructure:"Version"`
	BinPath         string `mapstructure:"BinPath"`
	DataPath        string `mapstructure:"DataPath"`
}

func newOpenlistStru() OpenlistStru {
	var path string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/openlist.exe"
	} else {
		path = "./ThirdPartyExt/openlist"
	}
	return OpenlistStru{
		Version:         "4.0.1",
		BinPath:         path,
		DataPath:        "./ThirdPartyExt/openlist_data",
		AutoStartEnable: false,
	}
}

type Caddy2Stru struct {
	AutoStartEnable bool   `mapstructure:"AutoStartEnable"`
	Version         string `mapstructure:"Version"`
	BinPath         string `mapstructure:"BinPath"`
	ConfigPath      string `mapstructure:"ConfigPath"`
}

func newCaddy2Config() Caddy2Stru {
	var path string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/caddy.exe"
	} else {
		path = "./ThirdPartyExt/caddy"
	}
	return Caddy2Stru{ // https://github.com/caddyserver/caddy/releases/download/v2.10.0/caddy_2.10.0_linux_amd64.tar.gz
		Version:         "2.10.0",
		BinPath:         path, // 实际解压到 caddy_2.10.0_linux_amd64/caddy
		ConfigPath:      "./ThirdPartyExt/Caddyfile",
		AutoStartEnable: false,
	}
}

type AcmeLegoStru struct {
	IsLegoAutoRenew         bool   `mapstructure:"IsLegoAutoRenew"`
	Version                 string `mapstructure:"Version"`
	BinPath                 string `mapstructure:"BinPath"`
	AutoUpdateCheckInterval int    `mapstructure:"AutoUpdateCheckInterval"` // 单位是小时
	Command                 string `mapstructure:"Command"`
	LEGO_PATH               string `mapstructure:"LEGO_PATH"`
}

func newAcmeLegoConfig() AcmeLegoStru {
	var path string
	var command string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/lego.exe"
		command = `
set LEGO_DEBUG_CLIENT_VERBOSE_ERROR=true
set LEGO_DEBUG_ACME_HTTP_CLIENT=true
set LEGO_EMAIL=you@example.com
set LEGO_PATH=${LEGO_PATH}
set CF_DNS_API_TOKEN=your-api-token
set LEGO_SERVER=https://acme.zerossl.com/v2/DV90
set LEGO_EAB_HMAC=your-hmac
set LEGO_EAB_KID=your-kid
${BinPath} --accept-tos  --dns cloudflare  -d exp1.com -d *.exp1.com  --eab -k ec256 renew &nascore
set ALICLOUD_ACCESS_KEY=abcdefghijklmnopqrstuvwx
set ALICLOUD_SECRET_KEY=your-secret-key
${BinPath} --accept-tos  --dns alidns  -d exp2.com -d *.exp2.com --eab -k ec256 renew &nascore
set CF_DNS_API_TOKEN=your-api-token2
${BinPath} --accept-tos  --dns cloudflare  -d exp3.com -d '*.exp3.com' --eab -k ec256 renew &nascore
`

	} else {
		path = "./ThirdPartyExt/lego"
		command = `
export LEGO_DEBUG_CLIENT_VERBOSE_ERROR=true
export LEGO_DEBUG_ACME_HTTP_CLIENT=true
export LEGO_EMAIL=you@example.com
export LEGO_PATH=${LEGO_PATH}
export CF_DNS_API_TOKEN=your-api-token
export LEGO_SERVER=https://acme.zerossl.com/v2/DV90
export LEGO_EAB_HMAC=your-hmac
export LEGO_EAB_KID=your-kid
${BinPath} --accept-tos  --dns cloudflare  -d exp1.com -d *.exp1.com  --eab -k ec256 renew &nascore
export ALICLOUD_ACCESS_KEY=abcdefghijklmnopqrstuvwx
export ALICLOUD_SECRET_KEY=your-secret-key
${BinPath} --accept-tos  --dns alidns  -d exp2.com -d *.exp2.com --eab -k ec256 renew &nascore
export CF_DNS_API_TOKEN=your-api-token2
${BinPath} --accept-tos  --dns cloudflare  -d exp3.com -d '*.exp3.com' --eab -k ec256 renew &nascore
`
	}
	return AcmeLegoStru{
		IsLegoAutoRenew:         false,
		Version:                 "4.25.1",
		BinPath:                 path,
		LEGO_PATH:               "./ThirdPartyExt/lego_cert",
		AutoUpdateCheckInterval: 24,
		Command:                 command,
	}
}

type AdGuardStru struct {
	IsAdGuardProxyEnable       bool   `mapstructure:"IsAdGuardProxyEnable"`
	ReverseproxyUrl            string `mapstructure:"ReverseproxyUrl"`
	Upstream_dns_file          string `mapstructure:"Upstream_dns_file"`
	Upstream_dns_fileUpdateUrl string `mapstructure:"Upstream_dns_fileUpdateUrl"`
	YouDohUrlDomain            string `mapstructure:"YouDohUrlDomain"`
	YouDohUrlSuffix            string `mapstructure:"YouDohUrlSuffix"`
	AutoUpdateRulesEnable      bool   `mapstructure:"AutoUpdateRulesEnable"`
	AutoUpdateRulesInterval    int    `mapstructure:"AutoUpdateRulesInterval"`
}

func newAdGuardConfig() AdGuardStru {
	return AdGuardStru{
		IsAdGuardProxyEnable:       false,
		ReverseproxyUrl:            "http://192.168.1.1:3000/",
		Upstream_dns_file:          "./adguard_upstream_dns_file.txt",
		Upstream_dns_fileUpdateUrl: "https://raw.githubusercontent.com/joyanhui/adguardhome-rules/refs/heads/release_file/ADG_chinaDirect_WinUpdate_Gfw.txt",
		YouDohUrlDomain:            "dns.cloudflare.com",
		YouDohUrlSuffix:            "dns-query",
		AutoUpdateRulesEnable:      false,
		AutoUpdateRulesInterval:    48,
	}
}

type RcloneExtStru struct {
	AutoMountEnable    bool   `mapstructure:"AutoMountEnable"`
	AutoMountCommand   string `mapstructure:"AutoMountCommand"`
	AutoUnMountCommand string `mapstructure:"AutoUnMountCommand"`
	Version            string `mapstructure:"Version"`
	BinPath            string `mapstructure:"BinPath"`
	ConfigFilePath     string `mapstructure:"ConfigFilePath"`
}
type DdnsgoStru struct {
	AutoStartEnable     bool   `mapstructure:"AutoStartEnable"`
	IsDDnsGOProxyEnable bool   `mapstructure:"IsDDnsGOProxyEnable"`
	ReverseproxyUrl     string `mapstructure:"ReverseproxyUrl"`
	ConfigFilePath      string `mapstructure:"ConfigFilePath"`
	BinPath             string `mapstructure:"BinPath"`
	Version             string `mapstructure:"Version"`
}

func newDefaultDDSN() DdnsgoStru {
	var path string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/ddns-go.exe"
	} else {
		path = "./ThirdPartyExt/ddns-go"
	}
	return DdnsgoStru{
		AutoStartEnable:     false,
		IsDDnsGOProxyEnable: false,
		Version:             "6.11.0",
		ReverseproxyUrl:     "http://localhost:9876/",
		BinPath:             path,
		ConfigFilePath:      "./ThirdPartyExt/ddnsgo_config.yaml",
	}
}

type ServerStru struct {
	HttpPort    int    `mapstructure:"httpPort"`
	HttpsEnable bool   `mapstructure:"HttpsEnable"`
	HttpsPort   int    `mapstructure:"httpsPort"`
	TlsCert     string `mapstructure:"tlscert"`
	TlsKey      string `mapstructure:"tlskey"`

	IsRunInServerLess bool `mapstructure:"IsRunInServerLess"` // 会让某些异步任务失效

	WebUIPrefix       string `mapstructure:"PrefixWebUI"`
	WebuiAndApiEnable bool   `mapstructure:"WebuiAndApiEnable"`
	ApiEnable         bool   `mapstructure:"ApiEnable"`
	WebDavEnable      bool   `mapstructure:"WebDavEnable"`
	TempFilePath      string `mapstructure:"TempFilePath"`

	DefaultStaticFileServicePrefix string `mapstructure:"DefaultStaticFileService"`
	DefaultStaticFileServiceEnable bool   `mapstructure:"DefaultStaticFileServiceEnable"`
	DefaultStaticFileServiceRoot   string `mapstructure:"DefaultStaticFileServiceRoot"`
}
type LimitStru struct {
	OnlineEditMaxSizeKB        int64 `mapstructure:"OnlineEditMaxSizeKB"`
	MaxFailedLoginsIpMap       int
	MaxFailedLoginSleepTimeSec int
}

func newDefaultServerConfig() ServerStru {
	var tempFilePath string
	if runtime.GOOS == "windows" {
		tempFilePath = "C:/Windows/Temp/nascore_socket/"
	} else {
		tempFilePath = "/tmp/nascore_socket/"
	}
	return ServerStru{
		HttpPort:          9000,
		HttpsEnable:       false,
		HttpsPort:         8181,
		TlsCert:           "domain.crt",
		TlsKey:            "domain.key",
		IsRunInServerLess: false,
		WebDavEnable:      true,
		ApiEnable:         true,
		TempFilePath:      tempFilePath,

		WebUIPrefix:                    "/@webui/",
		WebuiAndApiEnable:              true,
		DefaultStaticFileServicePrefix: "/@static/",
		DefaultStaticFileServiceEnable: true,
		DefaultStaticFileServiceRoot:   "./static/",
	}
}

// JwtStru JWT配置
type JwtStru struct {
	UserAccessTokenExpires  int64  `mapstructure:"user_access_token_expires"`
	UserRefreshTokenExpires int64  `mapstructure:"user_refresh_token_expires"`
	Issuer                  string `mapstructure:"issuer"`
}

// newDefaultJWTConfig 返回默认JWT配置
func newDefaultJWTConfig() JwtStru {
	return JwtStru{
		UserAccessTokenExpires:  2592000,
		UserRefreshTokenExpires: 7776000,
		Issuer:                  "nascore",
	}
}

type VodCacheStru struct {
	DoubanExpire    int `mapstructure:"DoubanExpire"`
	DoubanMax       int `mapstructure:"DoubanMax"`
	OtherExpire     int `mapstructure:"OtherExpire"`
	OtherMax        int `mapstructure:"OtherMax"`
	VoddetailExpire int `mapstructure:"VoddetailExpire"`
	VoddetailMax    int `mapstructure:"VoddetailMax"`
	VodlistExpire   int `mapstructure:"VodlistExpire"`
	VodlistMax      int `mapstructure:"VodlistMax"`
}

type VodSubscriptionStru struct {
	DefaultSelectedAPISite []string `mapstructure:"DefaultSelectedAPISite"`
	IntervalHour           int      `mapstructure:"IntervalHour"`
	Urls                   []string `mapstructure:"Urls"`
}
type LinksStru struct {
	LinksEnable bool `mapstructure:"LinksEnable"`
}
type VodExtStru struct {
	IsNeedLoginUse  bool                `mapstructure:"IsNeedLoginUse"`
	VodEnable       bool                `mapstructure:"VodEnable"`
	VodCache        VodCacheStru        `mapstructure:"VodCache"`
	VodSubscription VodSubscriptionStru `mapstructure:"VodSubscription"`
}

func newNascoreExtStru() NascoreExtStru {
	return NascoreExtStru{
		UserID:  "username",
		UserKey: "sdsds",
		Links: LinksStru{
			LinksEnable: true,
		},
		Vod: VodExtStru{
			IsNeedLoginUse: true,
			VodEnable:      true,
			VodCache: VodCacheStru{
				DoubanExpire:    150,
				DoubanMax:       50,
				OtherExpire:     25,
				OtherMax:        120,
				VoddetailExpire: 120,
				VoddetailMax:    120,
				VodlistExpire:   150,
				VodlistMax:      120,
			},
			VodSubscription: VodSubscriptionStru{
				DefaultSelectedAPISite: []string{"tyyszy", "bfzy", "dyttzy", "ruyi"},
				IntervalHour:           22,
				Urls: []string{
					"https://raw.githubusercontent.com/nas-core/nascore-website/refs/heads/main/docs/.vuepress/public/nascore_tv/subscription_example1.toml",
					"https://raw.githubusercontent.com/nas-core/nascore-website/refs/heads/main/docs/.vuepress/public/nascore_tv/subscription_example2.toml",
				},
			},
		},
	}
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *SysCfg {
	cfg := &SysCfg{
		Server: newDefaultServerConfig(),
		JWT:    newDefaultJWTConfig(),
		Secret: SecretStru{
			JwtSecret:      GenerateStr(1),
			Sha256HashSalt: GenerateStr(2),
			AESkey:         GenerateStr(3),
		},
		Limit: LimitStru{
			MaxFailedLoginsIpMap:       1000,
			MaxFailedLoginSleepTimeSec: 10,
			OnlineEditMaxSizeKB:        10240,
		},
		WebUICdnPrefix: "https://cdn.jsdmirror.com/gh/nas-core/nascore_static@main/",
		ThirdPartyExt: ThirdPartyExtStru{
			GitHubDownloadMirror: "https://github.akams.cn/",
			Openlist:             newOpenlistStru(),
			DdnsGO:               newDefaultDDSN(),
			Rclone:               newDefaultRclone(),
			AdGuard:              newAdGuardConfig(),
			AcmeLego:             newAcmeLegoConfig(),
			Caddy2:               newCaddy2Config(),
		},
		NascoreExt: newNascoreExtStru(),
	}
	// 统一补全目录路径结尾
	cfg.Server.TempFilePath = EnsureDirPathSuffix(cfg.Server.TempFilePath)
	cfg.ThirdPartyExt.Openlist.DataPath = EnsureDirPathSuffix(cfg.ThirdPartyExt.Openlist.DataPath)
	cfg.ThirdPartyExt.AcmeLego.LEGO_PATH = EnsureDirPathSuffix(cfg.ThirdPartyExt.AcmeLego.LEGO_PATH)
	return cfg
}

// 恢复 SecretStru 结构体
// SecretStru 密钥配置
type SecretStru struct {
	JwtSecret      string `mapstructure:"JwtSecret"`
	Sha256HashSalt string `mapstructure:"Sha256HashSalt"`
	AESkey         string `mapstructure:"AESkey"`
}

// 恢复 newDefaultRclone 函数
func newDefaultRclone() RcloneExtStru {
	var path string
	var autoMountCommand string
	var autoUnMountCommand string
	var configFilePath string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/rclone.exe"
		configFilePath = "N:/Scoop/apps/rclone/current/rclone.conf"
		autoMountCommand = `
${BinPath} mount oss_qd: D:/test/rclone/oss_qd --vfs-cache-mode writes --allow-non-empty ${ConfigFilePath} &nascore
${BinPath} mount jianguoyun: D:/test/rclone/jianguoyun --vfs-cache-mode writes --allow-non-empty ${ConfigFilePath} &nascore
`
		autoUnMountCommand = `
net use  D:/test/rclone/jianguoyun  /delete &nascore
net use  D:/test/rclone/oss_qd   /delete &nascore
`
	} else {
		path = "./ThirdPartyExt/rclone"
		configFilePath = "./ThirdPartyExt/rclone.conf"
		autoMountCommand = `
${BinPath} mount oss_qd: /home/yh/tmp/oss_qd --vfs-cache-mode writes --allow-non-empty ${ConfigFilePath} &nascore
${BinPath}  mount jianguoyun: /home/yh/tmp/jianguoyun --vfs-cache-mode writes --allow-non-empty ${ConfigFilePath} &nascore
`
		autoUnMountCommand = `
fusermount3 -u /home/yh/tmp &nascore
fusermount3 -u /home/yh/jianguoyun &nascore
fusermount3 -u /home/yh/oss_qd &nascore
`
	}
	return RcloneExtStru{
		Version:            "1.70.3",
		BinPath:            path,
		ConfigFilePath:     configFilePath,
		AutoMountEnable:    false,
		AutoMountCommand:   autoMountCommand,
		AutoUnMountCommand: autoUnMountCommand,
	}
}
