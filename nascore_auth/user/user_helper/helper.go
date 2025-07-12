package user_helper

import (
	"errors"
	"net/http"

	"github.com/nas-core/nascore/nascore_auth/jwt"
	"github.com/nas-core/nascore/nascore_auth/user/user_get_info"
	"github.com/nas-core/nascore/nascore_util/system_config"
)

// TokenUserInfo 包含token验证后的用户信息
type TokenUserInfo struct {
	Username  string
	HomeDir   string
	UserPerms string
	IsAdmin   bool
}

// ValidateTokenAndGetUserInfo 验证token并获取用户信息
func ValidateTokenAndGetUserInfo(r *http.Request, sys_cfg *system_config.SysCfg) (*TokenUserInfo, error) {
	var accessToken string
	var err error

	// .尝试从 Authorization header 获取 Bearer Token
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		accessToken, err = jwt.ExtractBearerToken(authHeader)
		if err != nil {
			accessToken = "" // 清空，确保后续会检查URL参数
		}
	}

	// 如果 header 中没有有效的 token, 尝试从 URL query parameter "token" 获取
	if accessToken == "" {
		accessTokenCookie, err := r.Cookie("cookieName")
		if err != nil {
			return nil, err
		}
		accessToken = accessTokenCookie.Value
		if accessToken == "" {
			return nil, errors.New("missing authorization token in header or query parameter")
		}
	}

	// 验证token
	claims, err := jwt.ValidateToken(accessToken, jwt.AccessTokenType, sys_cfg.Secret.JwtSecret)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	user, err := user_get_info.GetUserInfo(claims.Username, sys_cfg.Users)
	if err != nil {
		return nil, err
	}

	return &TokenUserInfo{
		Username: claims.Username,
		HomeDir:  user.HomeDir,
		IsAdmin:  user.IsAdmin,
	}, nil
}
