package main

import (
	"net"
	"net/http"
	"strings"
)

func getRealIP(r *http.Request) string {
	// 优先使用 X-Real-IP（常见于 nginx/apache）
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// Fallback to X-Forwarded-For（可能包含多个 IP）
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0] // 拿第一个
	}
	// 最后用默认的 remote address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
