package system_config

import (
	"log"

	"github.com/spf13/viper"
)

type SysCfg struct {
	Server         ServerConfig        `mapstructure:"Server"`
	JWT            JWTConfig           `mapstructure:"JWT"`
	Secret         SecretConfig        `mapstructure:"Secret"`
	WebUIPubLicCdn WebUIPubLicCdn      `mapstructure:"WebUIPubLicCdn"`
	Users          []map[string]string `mapstructure:"users"`
	Limit          LimitConfig         `mapstructure:"Limit"`
	ThirdPartyExt  ThirdPartyExtConfig `mapstructure:"ThirdPartyExt"`
}
type ThirdPartyExtConfig struct {
	Rclone  RcloneExtConfig `mapstructure:"Rclone"`
	DdnsGO  DdnsGOConfig    `mapstructure:"DdnsGO"`
	AdGuard AdGuardConfig   `mapstructure:"AdGuard"`
}
type AdGuardConfig struct {
	IsAdGuardProxyEnable       bool   `mapstructure:"IsAdGuardProxyEnable"`
	ReverseproxyUrl            string `mapstructure:"ReverseproxyUrl"`
	GitHubDownloadMirror       string `mapstructure:"GitHubDownloadMirror"`
	Upstream_dns_file          string `mapstructure:"Upstream_dns_file"`
	Upstream_dns_fileUpdateUrl string `mapstructure:"Upstream_dns_fileUpdateUrl"`
}

func newAdGuardConfig() AdGuardConfig {
	return AdGuardConfig{
		IsAdGuardProxyEnable:       false,
		ReverseproxyUrl:            "http://192.168.1.1:3000/",
		Upstream_dns_file:          "/overlay/data/adguard_upstream_dns_file.txt",
		GitHubDownloadMirror:       "https://github.akams.cn/",
		Upstream_dns_fileUpdateUrl: "https://raw.githubusercontent.com/joyanhui/adguardhome-rules/refs/heads/release_file/ADG_chinaDirect_WinUpdate_Gfw.txt",
	}
}

type RcloneExtConfig struct {
	DownLoadlink         string `mapstructure:"DownLoadlink"`
	GitHubDownloadMirror string `mapstructure:"GitHubDownloadMirror"`
	AutoMountEnable      bool   `mapstructure:"AutoMountEnable"`
	AutoMountCommand     string `mapstructure:"AutoMount"`
	AutoUnMountCommand   string `mapstructure:"AutoMount"`
}
type DdnsGOConfig struct {
	AutoStartEnable      bool   `mapstructure:"AutoStartEnable"`
	IsDDnsGOProxyEnable  bool   `mapstructure:"IsDDnsGOProxyEnable"`
	ReverseproxyUrl      string `mapstructure:"ReverseproxyUrl"`
	ConfigFilePath       string `mapstructure:"ConfigFilePath"`
	DdnsGOBinPath        string `mapstructure:"DdnsGOBinPath"`
	DownLoadlink         string `mapstructure:"DownLoadlink"`
	GitHubDownloadMirror string `mapstructure:"GitHubDownloadMirror"`
}

func newDefaultDDSN() DdnsGOConfig {
	return DdnsGOConfig{
		AutoStartEnable:      false,
		IsDDnsGOProxyEnable:  false,
		DownLoadlink:         "https://github.com/jeessy2/ddns-go/releases/download/v6.11.0/ddns-go_6.11.0_linux_x86_64.tar.gz",
		GitHubDownloadMirror: "https://github.akams.cn/",
		ReverseproxyUrl:      "http://localhost:9876/",
		DdnsGOBinPath:        "./ddns-go",
		ConfigFilePath:       "/home/yh/myworkspace/nas-core/code-private/ddns-go/config-ddnsgo.yaml",
	}
}

type ServerConfig struct {
	HttpPort    int    `mapstructure:"httpPort"`
	HttpsEnable bool   `mapstructure:"HttpsEnable"`
	HttpsPort   int    `mapstructure:"httpsPort"`
	TlsCert     string `mapstructure:"tlscert"`
	TlsKey      string `mapstructure:"tlskey"`

	WebUIPrefix string `mapstructure:"PrefixWebUI"`
	WebuiEnable bool   `mapstructure:"WebuiEnable"`

	WebDavEnable bool `mapstructure:"WebDavEnable"`

	DefaultStaticFileServicePrefix string `mapstructure:"DefaultStaticFileService"`
	DefaultStaticFileServiceEnable bool   `mapstructure:"DefaultStaticFileServiceEnable"`
	DefaultStaticFileServiceRoot   string `mapstructure:"DefaultStaticFileServiceRoot"`
}
type LimitConfig struct {
	OnlineEditMaxSizeKB        int64 `mapstructure:"OnlineEditMaxSizeKB"`
	MaxFailedLoginsIpMap       int
	MaxFailedLoginSleepTimeSec int
}

func newDefaultServerConfig() ServerConfig {
	return ServerConfig{
		HttpPort:                       9000,
		HttpsEnable:                    false,
		HttpsPort:                      8181,
		TlsCert:                        "cert.pem",
		TlsKey:                         "key.pem",
		WebDavEnable:                   true,
		WebUIPrefix:                    "/@webui/",
		WebuiEnable:                    true,
		DefaultStaticFileServicePrefix: "/@static/",
		DefaultStaticFileServiceEnable: true,
		DefaultStaticFileServiceRoot:   "./static/",
	}
}

// JWTConfig JWT配置
type JWTConfig struct {
	UserAccessTokenExpires  int64  `mapstructure:"user_access_token_expires"`
	UserRefreshTokenExpires int64  `mapstructure:"user_refresh_token_expires"`
	Issuer                  string `mapstructure:"issuer"`
}

// newDefaultJWTConfig 返回默认JWT配置
func newDefaultJWTConfig() JWTConfig {
	return JWTConfig{
		UserAccessTokenExpires:  2592000,
		UserRefreshTokenExpires: 7776000,
		Issuer:                  "nascore",
	}
}

// WebUIPubLicCdn CDN配置
type WebUIPubLicCdn struct {
	Header    string `mapstructure:"header"`
	Footer    string `mapstructure:"footer"`
	Dropzone  string `mapstructure:"dropzone"`
	Artplayer string `mapstructure:"artplayer"`
}

func newDefaultRclone() RcloneExtConfig {
	return RcloneExtConfig{
		DownLoadlink:         "https://github.com/rclone/rclone/releases/download/v1.70.1/rclone-v1.70.1-linux-amd64.zip",
		GitHubDownloadMirror: "https://github.akams.cn/",
		AutoMountEnable:      false,
		AutoMountCommand: `
rclone mount oss_qd: /home/yh/tmp/oss_qd --vfs-cache-mode writes --allow-non-empty  --config=/home/yh/.config/rclone/rclone.conf
rclone mount jianguoyun: /home/yh/tmp/jianguoyun --vfs-cache-mode writes --allow-non-empty  --config=/home/yh/.config/rclone/rclone.conf
`,
		AutoUnMountCommand: `
fusermount3 -u /home/yh/tmp
fusermount3 -u /home/yh/jianguoyun
`,
	}
}

// newDefaultWebUIPubLicCdn 返回默认CDN配置
func newDefaultWebUIPubLicCdn() WebUIPubLicCdn {
	return WebUIPubLicCdn{
		Header: `
<link href="https://lf3-cdn-tos.bytecdntp.com/cdn/expire-1-M/bootstrap/5.1.2/css/bootstrap.min.css" type="text/css"    rel="stylesheet" />
<link href="https://lf6-cdn-tos.bytecdntp.com/cdn/expire-1-M/bootstrap-icons/1.8.1/font/bootstrap-icons.css"    type="text/css" rel="stylesheet" />
<script src="https://lf26-cdn-tos.bytecdntp.com/cdn/expire-1-M/axios/0.26.0/axios.min.js" type="application/javascript"></script>
`,
		Footer: `
<script src="https://lf26-cdn-tos.bytecdntp.com/cdn/expire-1-M/bootstrap/5.1.2/js/bootstrap.bundle.min.js"  type="application/javascript"></script>
`,
		Dropzone: `<script src="https://unpkg.com/dropzone@5.9.3/dist/min/dropzone.min.js"></script>`,
		Artplayer: `
<script src="https://cdn.bootcdn.net/ajax/libs/hls.js/1.5.18/hls.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/artplayer/dist/artplayer.js"></script>
`,
	}
}

type SecretConfig struct {
	JwtSecret      string `mapstructure:"JwtSecret"`
	Sha256HashSalt string `mapstructure:"Sha256HashSalt"`
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *SysCfg {
	return &SysCfg{
		Server: newDefaultServerConfig(),
		JWT:    newDefaultJWTConfig(),
		Secret: SecretConfig{
			JwtSecret:      GenerateStr(1), // 创建
			Sha256HashSalt: GenerateStr(2),
		},
		Limit: LimitConfig{
			MaxFailedLoginsIpMap:       1000,
			MaxFailedLoginSleepTimeSec: 10,
			OnlineEditMaxSizeKB:        10240,
		},
		WebUIPubLicCdn: newDefaultWebUIPubLicCdn(),
		ThirdPartyExt: ThirdPartyExtConfig{
			DdnsGO:  newDefaultDDSN(),
			Rclone:  newDefaultRclone(),
			AdGuard: newAdGuardConfig(),
		},
		Users: []map[string]string{{
			"username": "admin",
			"passwd":   "admin",
			"home":     "/tmp",
			"isadmin":  "yes",
		}, {
			"username": "yh",
			"passwd":   "yh",
			"home":     "/home/yh/tmp", // 末尾不能是/开头
			"isadmin":  "no",
		}},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*SysCfg, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper.ReadInConfig file failed: ", err)
	}
	config := NewDefaultConfig() // 初始化 config 为指针类型
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatal("viper.Unmarshal failed: ", err)
	}
	log.Println("the system configuration loaded from file", configPath)

	if config.Secret.JwtSecret == "" {
		log.Fatal("config.Secret.JwtSecret is empty")
	}
	if config.Secret.Sha256HashSalt == "" {
		log.Fatal("config.Secret.Sha256HashSalt is empty")
	}
	if len(config.Users) > MaxUserLength {
		log.Fatal("config.Users length is greater than 5")
	}
	return config, err
}
