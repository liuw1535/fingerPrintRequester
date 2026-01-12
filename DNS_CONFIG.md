# DNS 配置说明

## 功能说明

TLS 请求器现在支持自定义 DNS 服务器配置，可以：
- 使用系统默认 DNS（默认行为）
- 指定自定义 DNS 服务器

## 配置方式

### 1. 在 config.json 中配置

```json
{
  "dns": {
    "servers": []
  }
}
```

- `servers`: DNS 服务器列表（格式：`IP:端口`）
  - 空数组 `[]` 或不配置：使用系统默认 DNS
  - 配置服务器：按顺序尝试，支持 fallback

### 2. 在请求中动态覆盖

```json
{
  "method": "POST",
  "url": "https://api.example.com",
  "headers": {...},
  "body": "...",
  "config_path": "./config.json",
  "dns": {
    "servers": ["8.8.8.8:53"]
  }
}
```

## 常用 DNS 服务器

### 公共 DNS

- **Google DNS**: `8.8.8.8:53`, `8.8.4.4:53`
- **Cloudflare DNS**: `1.1.1.1:53`, `1.0.0.1:53`
- **阿里 DNS**: `223.5.5.5:53`, `223.6.6.6:53`
- **腾讯 DNS**: `119.29.29.29:53`
- **114 DNS**: `114.114.114.114:53`

### 使用系统 DNS

```json
{
  "dns": {
    "servers": []
  }
}
```

或者在请求中不传 `dns` 字段。

## 配置示例

### 示例 1: 使用 Google DNS（带 fallback）

**config.json:**
```json
{
  "timeout": {
    "connect": 30,
    "read": 60
  },
  "dns": {
    "servers": ["8.8.8.8:53", "8.8.4.4:53"]
  },
  "proxy": {
    "enabled": false
  },
  "fingerprint": {...}
}
```

### 示例 2: 使用系统 DNS（默认）

**config.json:**
```json
{
  "timeout": {
    "connect": 30,
    "read": 60
  },
  "dns": {
    "servers": []
  },
  "proxy": {
    "enabled": false
  },
  "fingerprint": {...}
}
```

### 示例 3: 多 DNS 服务器 fallback

**config.json:**
```json
{
  "dns": {
    "servers": [
      "8.8.8.8:53",      // 优先使用 Google DNS
      "1.1.1.1:53",      // 失败则使用 Cloudflare DNS
      "223.5.5.5:53"     // 最后使用阿里 DNS
    ]
  }
}
```

### 示例 4: 请求中动态指定 DNS

**Python 调用:**
```python
import subprocess
import json

request = {
    "method": "POST",
    "url": "https://api.openai.com/v1/chat/completions",
    "headers": {
        "Authorization": "Bearer xxx",
        "Content-Type": "application/json"
    },
    "body": json.dumps({"stream": True, "messages": [...]}),
    "config_path": "./config.json",
    "dns": {
        "servers": ["1.1.1.1:53"]  # 使用 Cloudflare DNS
    }
}

proc = subprocess.Popen(
    ["./bin/fingerprint_windows_amd64.exe"],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE
)

stdout, stderr = proc.communicate(json.dumps(request).encode())
```

## 注意事项

1. DNS 服务器格式必须包含端口号（通常是 `:53`）
2. 支持配置多个 DNS 服务器，按顺序尝试（fallback 机制）
3. 如果第一个 DNS 服务器失败，会自动尝试下一个
4. 自定义 DNS 不影响代理功能，两者可以同时使用
5. 如果所有 DNS 服务器都失败，会返回错误
