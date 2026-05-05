package parser

import (
	"net/url"
	"strconv"
	"strings"
)

func parseTuic(u *url.URL) (*ProxyNode, error) {
	uuid := u.User.Username()
	password, _ := u.User.Password()

	port, _ := strconv.Atoi(u.Port())
	params := u.Query()

	node := &ProxyNode{
		Type:       "tuic",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
		UUID:       uuid,
		Password:   password,
	}

	node.SNI = params.Get("sni")
	node.Insecure = params.Get("insecure") == "1"

	if cc := params.Get("congestion_control"); cc != "" {
		node.CongestionControl = cc
	}
	if relay := params.Get("udp_relay_mode"); relay != "" {
		node.UDPRelayMode = relay
	}
	if alpn := params.Get("alpn"); alpn != "" {
		node.ALPN = strings.Split(alpn, ",")
	}

	return node, nil
}
