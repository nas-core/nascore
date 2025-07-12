package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/nas-core/nascore/nascore_auth/user/user_get_info"

	"github.com/nas-core/nascore/nascore_util/system_config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 定义 token 类型
type TokenType string

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

// Claims 定义 JWT claims
type Claims struct {
	Username string    `json:"username"`
	Type     TokenType `json:"type"`
	IsAdmin  bool      `json:"IsAdmin"`
	jwt.RegisteredClaims
}

// TokenResponse 令牌响应结构
type TokenResponse struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

// CreateTowTokens 创建访问令牌和刷新令牌
func CreateTowTokens(username string, nsCore_cfg *system_config.SysCfg) (*TokenResponse, error) {
	user, err := user_get_info.GetUserInfo(username, nsCore_cfg.Users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	// 创建访问令牌
	accessTokenExpires := time.Now().Add(time.Duration(nsCore_cfg.JWT.UserAccessTokenExpires) * time.Second)
	accessClaims := Claims{
		Username: username,
		Type:     AccessTokenType,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    nsCore_cfg.JWT.Issuer,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(nsCore_cfg.Secret.JwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// 创建刷新令牌
	refreshTokenExpires := time.Now().Add(time.Duration(nsCore_cfg.JWT.UserRefreshTokenExpires) * time.Second)
	refreshClaims := Claims{
		Username: username,
		Type:     RefreshTokenType,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    nsCore_cfg.JWT.Issuer,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(nsCore_cfg.Secret.JwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:         accessTokenString,
		RefreshToken:        refreshTokenString,
		AccessTokenExpires:  accessTokenExpires.Unix(),
		RefreshTokenExpires: refreshTokenExpires.Unix(),
	}, nil
}

// ValidateToken 验证令牌并返回 claims
func ValidateToken(tokenString string, expectedType TokenType, Secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Type != expectedType {
			return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractBearerToken 从 Authorization header 提取 Bearer token
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
