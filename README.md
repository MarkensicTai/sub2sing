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
| tuic (Tuic) | x | x |
| vless (VLESS) | x | x |

## 参考文档

- **sing-box 配置文档**：1.14.0（2025）
- **URL**：https://sing-box.sagernet.org/zh/configuration/

### 文档更新记录

| 日期 | sing-box 版本 | 同步状态 |
|------|-------------|---------|
| 2026-05-05 | 1.14.0 | 已同步 |

> 后续官方文档更新后，参考此表对比变更，同步更新转换逻辑。
