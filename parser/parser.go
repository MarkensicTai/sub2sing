package parser

import (
	"encoding/base64"
	"net/url"
	"regexp"
	"strings"
)

// ProxyNode 统一的代理节点中间结构
type ProxyNode struct {
	Type       string // trojan, hysteria2, anytls, shadowsocks
	Tag        string // 节点名称
	Server     string
	ServerPort int

	// 通用认证
	Password string

	// TLS 相关
	SNI              string
	Insecure         bool
	UTLSFingerprint  string // anytls 的 fp 参数
	ALPN             []string
	Peer             string // trojan 的 peer 参数

	// Shadowsocks 专用
	Method string

	// Hysteria2 专用
	ServerPorts  []string // 端口范围列表，如 ["20000:50000"]
	ObfsType     string
	ObfsPassword string
	HopInterval  string

	// VLESS 专用
	UUID           string
	Flow           string
	PacketEncoding string

	// TUIC 专用
	CongestionControl string
	UDPRelayMode      string
	ZeroRTTHandshake  bool

	// 传输层
	Network     string // tcp, udp
	Transport   string // ws, grpc, httpupgrade
	Host        string // transport host
	Path        string // ws/httpupgrade path
	ServiceName string // grpc service name
}

// ParseShareLink 解析单行 share link
func ParseShareLink(raw string) (*ProxyNode, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "trojan":
		return parseTrojan(u)
	case "hysteria2", "hy2":
		return parseHysteria2(u)
	case "anytls":
		return parseAnyTLS(u)
	case "ss":
		return parseShadowsocks(u)
	case "tuic":
		return parseTuic(u)
	case "vless":
		return parseVLess(u)
	default:
		return nil, nil // 忽略不支持的协议
	}
}

// ParseSubscription 解析 base64 编码的订阅内容
func ParseSubscription(data []byte) ([]*ProxyNode, error) {
	raw := strings.TrimSpace(string(data))
	// 如果以协议头开头，说明是明文格式
	if strings.HasPrefix(raw, "trojan://") || strings.HasPrefix(raw, "ss://") ||
		strings.HasPrefix(raw, "hysteria2://") || strings.HasPrefix(raw, "hy2://") ||
		strings.HasPrefix(raw, "anytls://") || strings.HasPrefix(raw, "tuic://") ||
		strings.HasPrefix(raw, "vless://") {
		return parseLines(strings.Split(raw, "\n")), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		// 尝试 URL-safe base64
		decoded, err = base64.URLEncoding.DecodeString(strings.TrimSpace(string(data)))
		if err != nil {
			// 尝试以 raw 格式解码（添加 padding）
			s := strings.TrimSpace(string(data))
			if missing := len(s) % 4; missing != 0 {
				s += strings.Repeat("=", 4-missing)
			}
			decoded, err = base64.URLEncoding.DecodeString(s)
			if err != nil {
				decoded, err = base64.StdEncoding.DecodeString(s)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return parseLines(strings.Split(string(decoded), "\n")), nil
}

func parseLines(lines []string) (nodes []*ProxyNode) {
	for _, line := range lines {
		node, err := ParseShareLink(line)
		if err != nil {
			continue
		}
		if node != nil {
			nodes = append(nodes, node)
		}
	}
	return
}

// sanitizeTag 生成合法的 tag 名称
func sanitizeTag(name string) string {
	re := regexp.MustCompile(`[\p{So}\p{Sk}]`)
	name = re.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "-")
	if len([]rune(name)) > 30 {
		name = string([]rune(name)[:30])
	}
	if name == "" {
		name = "proxy"
	}
	return name
}

// sanitizeTag