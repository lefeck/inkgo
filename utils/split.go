package utils

import "strings"

// ParseEmail 分割邮箱地址为用户名和域名
func ParseEmail(email string) (username string, domain string) {
	if email == "" {
		return "", ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
