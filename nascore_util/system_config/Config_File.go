package system_config

import (
	"log"

	"github.com/spf13/viper"
)

type SysCfg struct {
	Server         ServerStru `mapstructure:"Server"`
	JWT            JwtStru    `mapstructure:"JWT"`
	Secret         SecretStru `mapstructure:"Secret"`
	WebUIPubLicCdn WebUIStru  `mapstructure:"WebUIPubLicCdn"`
	//	Users          []map[string]string `mapstructure:"users"`
	//	WebSites       []WebsiteEntry      `mapstructure:"WebSites"`
	Limit LimitStru `mapstructure:"Limit"`

	NascoreExt    NascoreExtStru    `mapstructure:"NascoreExt"`
	ThirdPartyExt ThirdPartyExtStru `mapstructure:"ThirdPartyExt"`
}
type NascoreExtStru struct {
	UserID         string `mapstructure:"UserID"`
	UserKey        string `mapstructure:"UserKey"`
	UnixSocketPath string `mapstructure:"UnixSocketPath"`
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
	return OpenlistStru{
		DownLoadlink:    "https://github.com/OpenListTeam/OpenList/releases/download/v{ver}/openlist-{os}-{arch}.tar.gz",
		Version:         "4.0.1",
		BinPath:         "./openlist",
		DataPath:        "./openlist_data",
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
	return Caddy2Stru{ // https://github.com/caddyserver/caddy/releases/download/v2.10.0/caddy_2.10.0_linux_amd64.tar.gz
		DownLoadlink:    "https://github.com/caddyserver/caddy/releases/download/v{ver}/caddy_{ver}_{os}_{arch}.tar.gz",
		Version:         "2.10.0",
		BinPath:         "./caddy", // 实际解压到 caddy_2.10.0_linux_amd64/caddy
		ConfigPath:      "./Caddyfile",
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
	return AcmeLegoStru{
		IsLegoAutoRenew:         false,
		DownLoadlink:            "https://github.com/go-acme/lego/releases/download/v{ver}/lego_v{ver}_{os}_{arch}.tar.gz",
		Version:                 "4.23.1",
		BinPath:                 "./lego",
		LEGO_PATH:               "./lego_cert",
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
	return DdnsgoStru{
		AutoStartEnable:     false,
		IsDDnsGOProxyEnable: false,
		DownLoadlink:        "https://github.com/jeessy2/ddns-go/releases/download/v{ver}/ddns-go_{ver}_{os}_{arch}.tar.gz",
		Version:             "6.11.0",
		ReverseproxyUrl:     "http://localhost:9876/",
		BinPath:             "./ddns-go",
		ConfigFilePath:      "./ddnsgo_config.yaml",
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

// WebUIStru CDN配置
type WebUIStru struct {
	Header      string `mapstructure:"header"`
	Footer      string `mapstructure:"footer"`
	Dropzone    string `mapstructure:"dropzone"`
	Artplayer   string `mapstructure:"artplayer"`
	Tailwindcss string `mapstructure:"tailwindcss"`
}

func newDefaultRclone() RcloneExtStru {
	return RcloneExtStru{
		DownLoadlink:    "https://github.com/rclone/rclone/releases/download/v{ver}/rclone-v{ver}-{os}-{arch}.zip",
		Version:         "1.70.1",
		BinPath:         "./rclone",
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

// newDefaultWebUIPubLicCdn 返回默认CDN配置
func newDefaultWebUIPubLicCdn() WebUIStru {
	return WebUIStru{
		Header: `
<link href="https://lf3-cdn-tos.bytecdntp.com/cdn/expire-1-M/bootstrap/5.1.2/css/bootstrap.min.css" type="text/css"    rel="stylesheet" />
<link href="https://cdn.jsdmirror.com/npm/bootstrap-icons@1.13.1/font/bootstrap-icons.css"    type="text/css" rel="stylesheet" />
<script src="https://lf26-cdn-tos.bytecdntp.com/cdn/expire-1-M/axios/0.26.0/axios.min.js" type="application/javascript"></script>
`,
		Footer: `
<script src="https://lf26-cdn-tos.bytecdntp.com/cdn/expire-1-M/bootstrap/5.1.2/js/bootstrap.bundle.min.js"  type="application/javascript"></script>
`,
		Dropzone: `<script src="https://cdn.jsdmirror.com/npm/dropzone@5.9.3/dist/min/dropzone.min.js"></script><!--cdn.jsdelivr.net-->`,
		Artplayer: `
<script src="https://cdn.jsdmirror.com/npm/hls.js@1.5.18/dist/hls.min.js"></script>
<script src="https://cdn.jsdmirror.com/npm/artplayer/dist/artplayer.js"></script><!--cdn.jsdelivr.net-->
`,

		Tailwindcss: `<script src="https://cdn.jsdmirror.com/gh/nas-core/nascore_static@main/libs/tailwindcss.min.js"></script>`,
	}
}

type SecretStru struct {
	JwtSecret      string `mapstructure:"JwtSecret"`
	Sha256HashSalt string `mapstructure:"Sha256HashSalt"`
	AESkey         string `mapstructure:"AESkey"`
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *SysCfg {
	return &SysCfg{
		Server: newDefaultServerConfig(),
		//	WebSites: newWebSiteConfig(),
		JWT: newDefaultJWTConfig(),
		Secret: SecretStru{
			JwtSecret:      GenerateStr(1), // 创建
			Sha256HashSalt: GenerateStr(2),
			AESkey:         GenerateStr(3),
		},
		Limit: LimitStru{
			MaxFailedLoginsIpMap:       1000,
			MaxFailedLoginSleepTimeSec: 10,
			OnlineEditMaxSizeKB:        10240,
		},
		WebUIPubLicCdn: newDefaultWebUIPubLicCdn(),
		ThirdPartyExt: ThirdPartyExtStru{
			GitHubDownloadMirror: "https://github.akams.cn/",
			Openlist:             newOpenlistStru(),
			DdnsGO:               newDefaultDDSN(),
			Rclone:               newDefaultRclone(),
			AdGuard:              newAdGuardConfig(),
			AcmeLego:             newAcmeLegoConfig(),
			Caddy2:               newCaddy2Config(),
		},
		/*		Users: []map[string]string{{
					"username": "admin",
					"passwd":   "admin",
					"home":     "/tmp", // 末尾不能是/开头
					"isadmin":  "yes",
				}, {
					"username": "nascore",
					"passwd":   "nascore",
					"home":     "/tmp",
					"isadmin":  "yes",
				}, {
					"username": "yh",
					"passwd":   "yh",
					"home":     "/home/yh/tmp",
					"isadmin":  "no",
				}}, */
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
