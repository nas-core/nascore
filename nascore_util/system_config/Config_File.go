package system_config

import (
	"log"
	"runtime"

	"github.com/spf13/viper"
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
	DownLoadlink    string `mapstructure:"DownLoadlink"`
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
		DownLoadlink:    "https://github.com/OpenListTeam/OpenList/releases/download/v{ver}/openlist-{os}-{arch}.tar.gz",
		Version:         "4.0.1",
		BinPath:         path,
		DataPath:        "./ThirdPartyExt/openlist_data",
		AutoStartEnable: false,
	}
}

type Caddy2Stru struct {
	AutoStartEnable bool   `mapstructure:"AutoStartEnable"`
	DownLoadlink    string `mapstructure:"DownLoadlink"`
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
		DownLoadlink:    "https://github.com/caddyserver/caddy/releases/download/v{ver}/caddy_{ver}_{os}_{arch}.tar.gz",
		Version:         "2.10.0",
		BinPath:         path, // 实际解压到 caddy_2.10.0_linux_amd64/caddy
		ConfigPath:      "./ThirdPartyExt/Caddyfile",
		AutoStartEnable: false,
	}
}

type AcmeLegoStru struct {
	IsLegoAutoRenew         bool   `mapstructure:"IsLegoAutoRenew"`
	DownLoadlink            string `mapstructure:"DownLoadlink"`
	Version                 string `mapstructure:"Version"`
	BinPath                 string `mapstructure:"BinPath"`
	AutoUpdateCheckInterval int    `mapstructure:"AutoUpdateCheckInterval"` // 单位是小时
	Command                 string `mapstructure:"Command"`
	LEGO_PATH               string `mapstructure:"LEGO_PATH"`
}

func newAcmeLegoConfig() AcmeLegoStru {
	var path string
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/lego.exe"
	} else {
		path = "./ThirdPartyExt/lego"
	}
	return AcmeLegoStru{
		IsLegoAutoRenew:         false,
		DownLoadlink:            "https://github.com/go-acme/lego/releases/download/v{ver}/lego_v{ver}_{os}_{arch}.tar.gz",
		Version:                 "4.23.1",
		BinPath:                 path,
		LEGO_PATH:               "./ThirdPartyExt/lego_cert",
		AutoUpdateCheckInterval: 24,
		Command: `
LEGO_DEBUG_CLIENT_VERBOSE_ERROR=true
LEGO_DEBUG_ACME_HTTP_CLIENT=true
export LEGO_EMAIL="you@example.com"
export LEGO_PATH=${LEGO_PATH}

export CF_DNS_API_TOKEN=b9841238feb177a84330febba8a83208921177bffe733
${BinPath}  --dns cloudflare  -d example.com -d *.example.com --key-type ec256 run  &nascore
export ALICLOUD_ACCESS_KEY=abcdefghijklmnopqrstuvwx
export ALICLOUD_SECRET_KEY=your-secret-key
${BinPath}  --dns alidns  -d example2.com -d *.example2.com --key-type ec256 run  &nascore

`,
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
	DownLoadlink       string `mapstructure:"DownLoadlink"`
	AutoMountEnable    bool   `mapstructure:"AutoMountEnable"`
	AutoMountCommand   string `mapstructure:"AutoMountCommand"`
	AutoUnMountCommand string `mapstructure:"AutoUnMountCommand"`
	Version            string `mapstructure:"Version"`
	BinPath            string `mapstructure:"BinPath"`
}
type DdnsgoStru struct {
	AutoStartEnable     bool   `mapstructure:"AutoStartEnable"`
	IsDDnsGOProxyEnable bool   `mapstructure:"IsDDnsGOProxyEnable"`
	ReverseproxyUrl     string `mapstructure:"ReverseproxyUrl"`
	ConfigFilePath      string `mapstructure:"ConfigFilePath"`
	BinPath             string `mapstructure:"BinPath"`
	DownLoadlink        string `mapstructure:"DownLoadlink"`
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
		DownLoadlink:        "https://github.com/jeessy2/ddns-go/releases/download/v{ver}/ddns-go_{ver}_{os}_{arch}.tar.gz",
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

	DefaultStaticFileServicePrefix      string `mapstructure:"DefaultStaticFileService"`
	DefaultStaticFileServiceEnable      bool   `mapstructure:"DefaultStaticFileServiceEnable"`
	DefaultStaticFileServiceRoot        string `mapstructure:"DefaultStaticFileServiceRoot"`
	DefaultStaticFileServiceDownloadUrl string `mapstructure:"DefaultStaticFileServiceDownloadUrl"`
	UnixSocketFilePath                  string `mapstructure:"UnixSocketFilePath"`
}
type LimitStru struct {
	OnlineEditMaxSizeKB        int64 `mapstructure:"OnlineEditMaxSizeKB"`
	MaxFailedLoginsIpMap       int
	MaxFailedLoginSleepTimeSec int
}

func newDefaultServerConfig() ServerStru {
	return ServerStru{
		HttpPort:          9000,
		HttpsEnable:       false,
		HttpsPort:         8181,
		TlsCert:           "cert.pem",
		TlsKey:            "key.pem",
		IsRunInServerLess: false,
		WebDavEnable:      true,
		ApiEnable:         true,

		WebUIPrefix:                         "/@webui/",
		WebuiAndApiEnable:                   true,
		DefaultStaticFileServicePrefix:      "/@static/",
		DefaultStaticFileServiceEnable:      true,
		DefaultStaticFileServiceRoot:        "./static/",
		DefaultStaticFileServiceDownloadUrl: "https://github.com/nas-core/nascore_static/archive/refs/heads/main.zip",
		UnixSocketFilePath:                  "/tmp/nascore_socket/",
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

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*SysCfg, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("viper.ReadInConfig file failed: ", err)
	}
	config := NewDefaultConfig() // 初始化 config 为指针类型
	err := viper.Unmarshal(config)
	if err != nil {
		log.Println("viper.Unmarshal failed: ", err)
	}
	// log.Println("the system configuration loaded from file", configPath)

	if config.Secret.JwtSecret == "" {
		config.Secret.JwtSecret = GenerateStr(1)
		log.Println("config.Secret.JwtSecret is empty set :", config.Secret.JwtSecret)
	}
	if config.Secret.Sha256HashSalt == "" {
		config.Secret.Sha256HashSalt = GenerateStr(2)
		log.Println("config.Secret.Sha256HashSalt is empty set :", config.Secret.Sha256HashSalt)
	}
	if config.Secret.AESkey == "" {
		config.Secret.AESkey = GenerateStr(3)
		log.Println("config.Secret.AESkey is empty set :", config.Secret.AESkey)
	}
	/*	if len(config.Users) > MaxUserLength {
		config.Users = config.Users[:5] // 裁剪到5个
	} */
	return config, err
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *SysCfg {
	return &SysCfg{
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
	if runtime.GOOS == "windows" {
		path = "./ThirdPartyExt/rclone.exe"
	} else {
		path = "./ThirdPartyExt/rclone"
	}
	return RcloneExtStru{
		DownLoadlink:    "https://github.com/rclone/rclone/releases/download/v{ver}/rclone-v{ver}-{os}-{arch}.zip",
		Version:         "1.70.1",
		BinPath:         path,
		AutoMountEnable: false,
		AutoMountCommand: `
${BinPath} mount oss_qd: /home/yh/tmp/oss_qd --vfs-cache-mode writes --allow-non-empty  --config=/home/yh/.config/rclone/rclone.conf &nascore
${BinPath}  mount jianguoyun: /home/yh/tmp/jianguoyun --vfs-cache-mode writes --allow-non-empty  --config=/home/yh/.config/rclone/rclone.conf &nascore
`,
		AutoUnMountCommand: `
fusermount3 -u /home/yh/tmp &nascore
fusermount3 -u /home/yh/jianguoyun &nascore
fusermount3 -u /home/yh/oss_qd &nascore
`,
	}
}
