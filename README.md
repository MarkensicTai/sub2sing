# sub2sing

v2ray 订阅链接转 sing-box 配置文件。

## 用法

```bash
# 从 URL
./sub2sing -url "https://..." -o config.json

# 从本地文件
./sub2sing -file sub.txt -o config.json

# 输出到标准输出
./sub2sing -url "https://..."
```

## 支持的协议

| 协议 | 解析 | 生成 |
|------|------|------|
| trojan | ✓ | ✓ |
| hysteria2 | ✓ | ✓ |
| anytls | ✓ | ✓ |
| shadowsocks (SIP002) | ✓ | ✓ |
| tuic (Tuic) | ✓ | ✓ |
| vless (VLESS) | ✓ | ✓ |

## 参考文档

- **sing-box 配置文档**：1.14.0（2025）
- **URL**：https://sing-box.sagernet.org/zh/configuration/

### 文档更新记录

| 日期 | sing-box 版本 | 同步状态 |
|------|-------------|---------|
| 2026-05-05 | 1.14.0 | 已同步 |

> 后续官方文档更新后，参考此表对比变更，同步更新转换逻辑。

## 交叉编译

```bash
# macOS / Linux
GOOS=android GOARCH=arm64 go build -o sub2sing-android-arm64 .

# 如需 amd64 模拟器 / 旧设备
GOOS=android GOARCH=amd64 go build -o sub2sing-android-amd64 .
```

## 已知问题

- 规则集下载使用 `download_detour: "direct"` 而非 `http_client`。`http_client` 是 1.14.0 引入的替代字段，但官方 stable 版本不识别该字段（报 `unknown field "http_client"`），因此沿用 `download_detour`。
