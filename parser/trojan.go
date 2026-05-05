package parser

import (
	"net/url"
	"strconv"
)

func parseTrojan(u *url.URL) (*ProxyNode, error) {
	password := u.User.Username()
	if password == "" {
		password, _ = u.User.Password()
	}

	port, _ := strconv.Atoi(u.Port())
	params := u.Query()

	node := &ProxyNode{
		Type:       "trojan",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
		Password:   password,
	}

	// TLS
	node.Peer = params.Get("peer")
	node.SNI = params.Get("sni")
	if node.SNI == "" {
		node.SNI = node.Peer
	}
	node.Insecure = params.Get("allowInsecure") == "1"

	// ALPN
	if alpn := params.Get("alpn"); alpn != "" {
		node.ALPN = []string{alpn}
	}

	// 传输层
	node.Network = "tcp"
	switch params.Get("type") {
	case "ws":
		node.Transport = "ws"
		node.Path = params.Get("path")
		node.Host = params.Get("host")
	case "grpc":
		node.Transport = "grpc"
		node.ServiceName = params.Get("serviceName")
	case "httpupgrade":
		node.Transport = "httpupgrade"
		node.Path = params.Get("path")
		node.Host = params.Get("host")
	}

	return node, nil
}
