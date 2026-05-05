package parser

import (
	"net/url"
	"strconv"
)

func parseAnyTLS(u *url.URL) (*ProxyNode, error) {
	password := u.User.Username()
	if password == "" {
		password, _ = u.User.Password()
	}

	port, _ := strconv.Atoi(u.Port())
	params := u.Query()

	node := &ProxyNode{
		Type:       "anytls",
		Tag:        sanitizeTag(u.Fragment),
		Server:     u.Hostname(),
		ServerPort: port,
		Password:   password,
	}

	node.SNI = params.Get("sni")
	node.Insecure = params.Get("insecure") == "1"

	// uTLS fingerprint
	if fp := params.Get("fp"); fp != "" {
		node.UTLSFingerprint = fp
	}

	// ALPN
	if alpn := params.Get("alpn"); alpn != "" {
		node.ALPN = []string{alpn}
	}

	return node, nil
}
