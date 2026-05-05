package parser

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"
)

func parseShadowsocks(u *url.URL) (*ProxyNode, error) {
	port, _ := strconv.Atoi(u.Port())

	node := &ProxyNode{
		Type:       "shadowsocks",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
	}

	// SIP002 格式: ss://base64(method:password)@host:port
	// URL 解析后 username 就是 base64 编码的 method:password
	userinfo := u.User.Username()
	if userinfo == "" {
		// 旧格式: ss://base64(method:password@host:port)
		// 需要重新解析
		return parseLegacySS(u)
	}

	// 尝试 base64 解码
	decoded, err := base64.RawURLEncoding.DecodeString(userinfo)
	if err != nil {
		decoded, err = base64.StdEncoding.DecodeString(userinfo)
		if err != nil {
			// 添加 padding 重试
			s := userinfo
			if missing := len(s) % 4; missing != 0 {
				s += strings.Repeat("=", 4-missing)
			}
			decoded, err = base64.StdEncoding.DecodeString(s)
			if err != nil {
				decoded, _ = base64.RawURLEncoding.DecodeString(s)
			}
		}
	}

	if decoded != nil {
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) == 2 {
			node.Method = parts[0]
			node.Password = parts[1]
		}
	}

	return node, nil
}

func parseLegacySS(u *url.URL) (*ProxyNode, error) {
	// 兼容旧格式，直接跳过
	return nil, nil
}
