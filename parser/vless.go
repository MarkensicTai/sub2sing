package parser

import (
	"net/url"
	"strconv"
)

func parseVLess(u *url.URL) (*ProxyNode, error) {
	uuid := u.User.Username()

	port, _ := strconv.Atoi(u.Port())
	params := u.Query()

	node := &ProxyNode{
		Type:       "vless",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
		UUID:       uuid,
	}

	node.SNI = params.Get("sni")
	node.Insecure = params.Get("insecure") == "1"

	if flow := params.Get("flow"); flow != "" {
		node.Flow = flow
	}
	if fp := params.Get("fp"); fp != "" {
		node.UTLSFingerprint = fp
	}

	// ALPN
	if alpn := params.Get("alpn"); alpn != "" {
		node.ALPN = []string{alpn}
	}

	// 传输层
	switch params.Get("type") {
	case "ws":
		node.Transport = "ws"
		node.Path = params.Get("path")
		node.Host = params.Get("host")
	case "grpc":
		node.Transport = "grpc"
		node.ServiceName = params.Get("serviceName")
	}

	return node, nil
}
