package parser

import (
	"net/url"
	"strconv"
	"strings"
)

func parseHysteria2(u *url.URL) (*ProxyNode, error) {
	password := u.User.Username()
	if password == "" {
		password, _ = u.User.Password()
	}

	port, _ := strconv.Atoi(u.Port())
	params := u.Query()

	node := &ProxyNode{
		Type:       "hysteria2",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
		Password:   password,
	}

	node.SNI = params.Get("sni")
	node.Insecure = params.Get("insecure") == "1"

	// 端口跳跃范围
	if mport := params.Get("mport"); mport != "" {
		var ports []string
		for _, p := range strings.Split(mport, ",") {
			// sing-box 要求冒号分隔，v2ray 订阅用短横线
			if strings.Contains(p, "-") && !strings.Contains(p, ":") {
				parts := strings.SplitN(p, "-", 2)
				if len(parts) == 2 {
					p = parts[0] + ":" + parts[1]
				}
			}
			ports = append(ports, p)
		}
		node.ServerPorts = ports
	}

	// Obfs
	if obfs := params.Get("obfs"); obfs != "" {
		node.ObfsType = obfs
		node.ObfsPassword = params.Get("obfs-password")
	}

	// 端口跳跃间隔
	if hop := params.Get("hopinterval"); hop != "" {
		node.HopInterval = hop
	}

	return node, nil
}
