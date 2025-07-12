package getip

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			// X-Forwarded-For 可能会包含多个 IP，第一个通常是客户端真实 IP
			return strings.TrimSpace(ips[0])
		}
	}
	// 尝试从 X-Real-Ip 头部获取 IP
	if xRealIP := r.Header.Get("X-Real-Ip"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	// 最后，使用 r.RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return strings.TrimSpace(ip)
	}

	return strings.TrimSpace(r.RemoteAddr)
}
