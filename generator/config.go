package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"sub2sing/parser"
)

// ---------- sing-box 配置结构体 ----------

type Config struct {
	Log          *Log          `json:"log,omitempty"`
	DNS          *DNS          `json:"dns,omitempty"`
	NTP          *NTP          `json:"ntp,omitempty"`
	Inbounds     []Inbound     `json:"inbounds"`
	Outbounds    []Outbound    `json:"outbounds"`
	Route        *Route        `json:"route,omitempty"`
	Experimental *Experimental `json:"experimental,omitempty"`
}

type Log struct {
	Level string `json:"level"`
}

type NTP struct {
	Enabled    bool   `json:"enabled"`
	Server     string `json:"server,omitempty"`
	ServerPort int    `json:"server_port,omitempty"`
	Interval   string `json:"interval,omitempty"`
	Detour     string `json:"detour,omitempty"`
}

// ---------- DNS ----------

type DNS struct {
	Servers          []DNSServer `json:"servers"`
	Rules            []DNSRule   `json:"rules"`
	Final            string      `json:"final"`
	Strategy       string `json:"strategy,omitempty"`
	ReverseMapping bool   `json:"reverse_mapping,omitempty"`
}

type DNSServer struct {
	Type            string `json:"type,omitempty"`
	Tag             string `json:"tag"`
	Address         string `json:"address,omitempty"`
	Server          string `json:"server,omitempty"`
	ServerPort      int    `json:"server_port,omitempty"`
	AddressResolver string `json:"address_resolver,omitempty"`
	Detour          string `json:"detour,omitempty"`
	Inet4Range      string `json:"inet4_range,omitempty"`
	Inet6Range      string `json:"inet6_range,omitempty"`
}

type DNSRule struct {
	RuleSet   []string `json:"rule_set,omitempty"`
	Geosite   []string `json:"geosite,omitempty"`
	QueryType []string `json:"query_type,omitempty"`
	ClashMode string   `json:"clash_mode,omitempty"`
	Action    string   `json:"action,omitempty"`
	Server    string   `json:"server,omitempty"`
}

// ---------- Inbound ----------

type Inbound struct {
	Type          string      `json:"type"`
	Tag           string      `json:"tag"`
	Listen        string      `json:"listen,omitempty"`
	ListenPort    int         `json:"listen_port,omitempty"`
	Address       []string    `json:"address,omitempty"`
	InterfaceName string      `json:"interface_name,omitempty"`
	MTU           int         `json:"mtu,omitempty"`
	Stack         string      `json:"stack,omitempty"`
	AutoRoute     bool        `json:"auto_route,omitempty"`
	StrictRoute   bool        `json:"strict_route,omitempty"`
	Platform      *Platform   `json:"platform,omitempty"`
}

type Platform struct {
	HTTPProxy *HTTPProxy `json:"http_proxy,omitempty"`
}

type HTTPProxy struct {
	Enabled    bool   `json:"enabled"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
}

// ---------- Outbound ----------

type Outbound struct {
	IdleTimeout string `json:"idle_timeout,omitempty"`
	Type        string      `json:"type"`
	Tag         string      `json:"tag"`
	Server      string      `json:"server,omitempty"`
	ServerPort  int         `json:"server_port,omitempty"`
	Password    string      `json:"password,omitempty"`
	Method      string      `json:"method,omitempty"`
	Network     string      `json:"network,omitempty"`
	TLS         *TLS        `json:"tls,omitempty"`
	Transport   *Transport  `json:"transport,omitempty"`
	Multiplex   *Multiplex  `json:"multiplex,omitempty"`
	Outbounds   []string    `json:"outbounds,omitempty"`
	URL         string      `json:"url,omitempty"`
	Interval    string      `json:"interval,omitempty"`
	Tolerance   int         `json:"tolerance,omitempty"`
	Default     string      `json:"default,omitempty"`
	InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`

	// Hysteria2 专用
	ServerPorts    []string `json:"server_ports,omitempty"`
	HopInterval    string   `json:"hop_interval,omitempty"`
	HopIntervalMax string   `json:"hop_interval_max,omitempty"`
	Obfs           *Obfs    `json:"obfs,omitempty"`

	// AnyTLS 专用
	IdleSessionCheckInterval string `json:"idle_session_check_interval,omitempty"`
	IdleSessionTimeout      string `json:"idle_session_timeout,omitempty"`
	MinIdleSession          int    `json:"min_idle_session,omitempty"`

	// VLESS 专用
	UUID           string `json:"uuid,omitempty"`
	Flow           string `json:"flow,omitempty"`
	PacketEncoding string `json:"packet_encoding,omitempty"`

	// TUIC 专用
	CongestionControl string `json:"congestion_control,omitempty"`
	UDPRelayMode      string `json:"udp_relay_mode,omitempty"`
	ZeroRTTHandshake  bool   `json:"zero_rtt_handshake,omitempty"`
	Heartbeat         string `json:"heartbeat,omitempty"`

	// 域名解析
	DomainResolver any `json:"domain_resolver,omitempty"`
}

type TLS struct {
	Enabled    bool     `json:"enabled"`
	DisableSNI bool     `json:"disable_sni,omitempty"`
	ServerName string   `json:"server_name,omitempty"`
	Insecure   bool     `json:"insecure,omitempty"`
	ALPN       []string `json:"alpn,omitempty"`
	UTLS       *UTLS    `json:"utls,omitempty"`
}

type UTLS struct {
	Enabled     bool   `json:"enabled"`
	Fingerprint string `json:"fingerprint"`
}

type Transport struct {
	Type        string            `json:"type"`
	Path        string            `json:"path,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	ServiceName string            `json:"service_name,omitempty"`
}

type Multiplex struct {
	Enabled  bool   `json:"enabled"`
	Protocol string `json:"protocol,omitempty"`
}

type Obfs struct {
	Type     string `json:"type"`
	Password string `json:"password"`
}

// ---------- Route ----------

type Route struct {
	RuleSet               []RuleSet   `json:"rule_set,omitempty"`
	Rules                 []RouteRule `json:"rules"`
	Final                 string      `json:"final"`
	AutoDetectInterface   bool        `json:"auto_detect_interface"`
	DefaultDomainResolver any         `json:"default_domain_resolver,omitempty"`
}

type RuleSet struct {
	Type           string `json:"type"`
	Tag            string `json:"tag"`
	Format         string `json:"format"`
	URL            string `json:"url,omitempty"`
	UpdateInterval string `json:"update_interval,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty"`
}

type RouteRule struct {
	Method  string   `json:"method,omitempty"`
	NoDrop  bool     `json:"no_drop,omitempty"`
	Protocol     string   `json:"protocol,omitempty"`
	RuleSet      []string `json:"rule_set,omitempty"`
	Inbound      []string `json:"inbound,omitempty"`
	Geosite      []string `json:"geosite,omitempty"`
	Geoip        []string `json:"geoip,omitempty"`
	DomainSuffix []string `json:"domain_suffix,omitempty"`
	IPIsPrivate  bool     `json:"ip_is_private,omitempty"`
	ClashMode    string   `json:"clash_mode,omitempty"`
	Action       string   `json:"action,omitempty"`
	Outbound     string   `json:"outbound,omitempty"`
	Timeout      string   `json:"timeout,omitempty"`
}

// ---------- Experimental ----------

type Experimental struct {
	CacheFile *CacheFile `json:"cache_file,omitempty"`
	ClashAPI  *ClashAPI  `json:"clash_api,omitempty"`
}

type CacheFile struct {
	Enabled     bool   `json:"enabled"`
	Path        string `json:"path,omitempty"`
	StoreFakeIP bool   `json:"store_fakeip,omitempty"`
	StoreDNS    bool   `json:"store_dns,omitempty"`
}

type ClashAPI struct {
	DefaultMode        string `json:"default_mode,omitempty"`
	ExternalController string `json:"external_controller,omitempty"`
	Secret             string `json:"secret,omitempty"`
}

// ---------- 生成逻辑 ----------

func firstProxyTag(tags []string) string {
	if len(tags) > 0 {
		return tags[0]
	}
	return ""
}

// Generate 生成完整的 sing-box config
func Generate(nodes []*parser.ProxyNode) *Config {
	// 收集所有代理节点的标签
	var proxyTags []string
	var outbounds []Outbound

	for _, n := range nodes {
		out := convertNode(n)
		outbounds = append(outbounds, out)
		proxyTags = append(proxyTags, out.Tag)
	}

	// 基础出站
	baseOutbounds := []Outbound{
		{Type: "direct", Tag: "direct", DomainResolver: map[string]interface{}{"server": "local"}},
}

	// urltest 自动选择
	if len(proxyTags) > 0 {
		autoOutbound := Outbound{
			Type:      "urltest",
			Tag:       "auto",
			Outbounds: proxyTags,
			URL:       "https://www.gstatic.com/generate_204",
			Interval:  "10m",
			Tolerance: 50,
			IdleTimeout: "30m",
			InterruptExistConnections: false,
		}
		baseOutbounds = append(baseOutbounds, autoOutbound)

		// selector 手动选择，包含 auto + 所有节点
		selectorTags := append([]string{"auto"}, proxyTags...)
		baseOutbounds = append(baseOutbounds, Outbound{
			Type:      "selector",
			Tag:       "select",
			Outbounds: selectorTags,
			Default:   "auto",
			InterruptExistConnections: true,
		})
	}

	outbounds = append(baseOutbounds, outbounds...)

	config := &Config{
		Log: &Log{Level: "info"},
		NTP: &NTP{
			Enabled:    true,
			Server:     "time.apple.com",
			ServerPort: 123,
			Interval:   "30m",
			Detour:     "direct",
		},
		DNS:       buildDNS(),
		Inbounds:  buildInbounds(),
		Outbounds: outbounds,
		Route:     buildRoute(),
		Experimental: &Experimental{
			CacheFile: &CacheFile{Enabled: true},
			ClashAPI: &ClashAPI{
				DefaultMode:        "Rule",
				ExternalController: "127.0.0.1:9090",
				Secret:             "",
			},
		},
	}

	return config
}

func buildDNS() *DNS {
	return &DNS{
		Servers: []DNSServer{
			{
				Type: "local",
				Tag:  "local",
			},
			{
				Type:   "udp",
				Tag:    "remote",
				Server: "1.1.1.1",
			},
			{
				Type:   "udp",
				Tag:    "cn",
				Server: "223.5.5.5",
			},
		},
		Rules: []DNSRule{
			{
				ClashMode: "Direct",
				Server:    "local",
				Action:    "route",
			},
			{
				ClashMode: "Global",
				Server:    "remote",
				Action:    "route",
			},
			{
				RuleSet: []string{"geosite-cn"},
				Server:  "cn",
				Action:  "route",
			},
		},
		Final:            "remote",
	}
}

func buildInbounds() []Inbound {
	return []Inbound{
		{
			Type:       "mixed",
			Tag:        "mixed-in",
			Listen:     "127.0.0.1",
			ListenPort: 2080,
		},
		{
			Type:          "tun",
			Tag:           "tun-in",
			InterfaceName: "tun0",
			Address:       []string{"172.19.0.1/30", "fdfe:dcba:9876::1/126"},
			MTU:           9000,
			Stack:         "system",
			AutoRoute:     true,
			StrictRoute:   true,
			Platform: &Platform{
				HTTPProxy: &HTTPProxy{
					Enabled:    false,
					Server:     "127.0.0.1",
					ServerPort: 2080,
				},
			},
		},
	}
}

func buildRoute() *Route {
	return &Route{
		DefaultDomainResolver: map[string]interface{}{"server": "remote"},
		Rules: []RouteRule{
			{
				Action: "sniff",
			},
			{
				Protocol: "dns",
				Action:   "hijack-dns",
			},
			{
				ClashMode: "Direct",
				Outbound:  "direct",
				Action:    "route",
			},
			{
				ClashMode: "Global",
				Outbound:  "select",
				Action:    "route",
			},
			{
				RuleSet:  []string{"category-ads-all"},
				Action:   "reject",
			},
			{
				IPIsPrivate: true,
				Outbound:    "direct",
				Action:      "route",
			},
			{
				RuleSet:  []string{"geosite-cn", "geoip-cn"},
				Outbound: "direct",
				Action:   "route",
			},
		},
		Final:              "select",
		AutoDetectInterface: true,
		RuleSet: []RuleSet{
			{
				Tag:            "geosite-geolocation-!cn",
				Type:           "remote",
				Format:         "binary",
				URL:            "https://cdn.jsdelivr.net/gh/SagerNet/sing-geosite@rule-set/geosite-geolocation-!cn.srs",
				DownloadDetour: "direct",
			},
			{
				Tag:            "geoip-cn",
				Type:           "remote",
				Format:         "binary",
				URL:            "https://cdn.jsdelivr.net/gh/SagerNet/sing-geoip@rule-set/geoip-cn.srs",
				DownloadDetour: "direct",
			},
			{
				Tag:            "geosite-cn",
				Type:           "remote",
				Format:         "binary",
				URL:            "https://cdn.jsdelivr.net/gh/SagerNet/sing-geosite@rule-set/geosite-cn.srs",
				DownloadDetour: "direct",
			},
			{
				Tag:            "category-ads-all",
				Type:           "remote",
				Format:         "binary",
				URL:            "https://cdn.jsdelivr.net/gh/SagerNet/sing-geosite@rule-set/geosite-category-ads-all.srs",
				DownloadDetour: "direct",
			},
		},
	}
}

// convertNode 将解析后的节点转为 sing-box outbound
func convertNode(n *parser.ProxyNode) Outbound {
	out := Outbound{
		Type:           n.Type,
		Tag:            ensureUniqueTag(n.Tag),
		Server:         n.Server,
		ServerPort:     n.ServerPort,
		Password:       n.Password,
		DomainResolver: "local",
	}

	// multiplex: 仅 trojan 支持（anytls/hysteria2 不支持此字段, ss 与 udp_over_tcp 冲突）
	if n.Type == "trojan" {
		out.Multiplex = &Multiplex{Enabled: true, Protocol: "smux"}
	}

	// TLS 配置：trojan/hysteria2/anytls 都需要
	if n.Type == "trojan" || n.Type == "hysteria2" || n.Type == "anytls" || n.Type == "tuic" || n.Type == "vless" {
		out.TLS = &TLS{
			Enabled:    true,
			ServerName: n.SNI,
			Insecure:   n.Insecure,
			ALPN:       n.ALPN,
		}
		// uTLS 指纹
		if n.UTLSFingerprint != "" && (n.Type == "anytls" || n.Type == "trojan" || n.Type == "vless") {
			out.TLS.UTLS = &UTLS{
				Enabled:     true,
				Fingerprint: n.UTLSFingerprint,
			}
		}
	}

	// 传输层
	if n.Transport != "" {
		out.Transport = &Transport{
			Type:        n.Transport,
			Path:        n.Path,
			ServiceName: n.ServiceName,
		}
		if n.Host != "" {
			out.Transport.Headers = map[string]string{"Host": n.Host}
		}
	}

	// Shadowsocks 专用
	if n.Type == "shadowsocks" {
		out.Method = n.Method
		out.TLS = nil
		out.Multiplex = nil // SS 2022 与 udp_over_tcp 冲突
	}

	// Hysteria2 专用
	if n.Type == "hysteria2" {
		if len(n.ServerPorts) > 0 {
			out.ServerPorts = n.ServerPorts
			out.ServerPort = 0
		}
		if n.HopInterval != "" {
			out.HopInterval = n.HopInterval
		}
		if n.ObfsType != "" {
			out.Obfs = &Obfs{
				Type:     n.ObfsType,
				Password: n.ObfsPassword,
			}
		}
		out.Multiplex = nil // Hysteria2 自带多路复用
	}

	// AnyTLS 闲置会话管理
	if n.Type == "anytls" {
		out.IdleSessionCheckInterval = "30s"
		out.IdleSessionTimeout = "30s"
		out.MinIdleSession = 0
	}

	// VLESS 专用
	if n.Type == "vless" {
		out.UUID = n.UUID
		out.Flow = n.Flow
		out.PacketEncoding = "xudp"
		out.Multiplex = nil // vless 配套 reality 不适合 smux
	}

	// TUIC 专用
	if n.Type == "tuic" {
		out.UUID = n.UUID
		if n.CongestionControl != "" {
			out.CongestionControl = n.CongestionControl
		}
		if n.UDPRelayMode != "" {
			out.UDPRelayMode = n.UDPRelayMode
		}
		out.Multiplex = nil // tuic 自带多路复用
	}

	return out
}

var tagCounter = make(map[string]int)

func ensureUniqueTag(tag string) string {
	tag = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == '"' || r == '\'' {
			return '-'
		}
		return r
	}, tag)

	count := tagCounter[tag]
	tagCounter[tag]++
	if count > 0 {
		return fmt.Sprintf("%s-%d", tag, count)
	}
	return tag
}

// ToJSON 序列化为 JSON
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
