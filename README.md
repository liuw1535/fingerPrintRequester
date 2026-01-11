# TLS Requester

基于 Go utls 库的 TLS 指纹请求器，支持自定义 TLS 指纹、GREASE、代理和流式传输。

## 编译

```bash
go mod tidy
go build -o tlsRequester.exe
```

## 配置文件

编辑 `config.json` 配置 TLS 指纹参数：

- `timeout`: 连接和读取超时设置
- `proxy`: 代理配置（支持 HTTP 和 SOCKS5）
- `fingerprint`: TLS 指纹配置
  - `grease`: 是否启用 GREASE
  - `http2`: 是否使用 HTTP/2
  - `ciphers`: 密码套件列表（保证顺序）
  - `extensions`: 扩展列表（保证顺序）

### 扩展参数说明

- `supported_groups`: `{"curves": ["X25519", "CurveP256", "CurveP384"]}`
- `signature_algorithms`: `{"algorithms": ["ECDSAWithP256AndSHA256", ...]}`
- `key_share`: `{"groups": ["X25519"]}`
- `psk_key_exchange_modes`: `{"modes": [1]}`
- `supported_versions`: `{"versions": ["0x0304", "0x0303"]}`
- `padding`: `{"length": 75}`
- `pre_shared_key`: `{"identity_length": 138, "binder_length": 32}`

## 使用方法

### 输入格式（stdin JSON）

```json
{
  "method": "POST",
  "url": "https://api.example.com/v1/chat",
  "headers": {
    "Authorization": "Bearer xxx",
    "Content-Type": "application/json"
  },
  "body": "{\"stream\":true}",
  "config_path": "./config.json"
}
```

### 输出格式（stdout）

成功时直接输出完整 HTTP 响应（包括响应头和响应体）：

```
HTTP/1.1 200 OK
Content-Type: text/event-stream
Transfer-Encoding: chunked

data: {"chunk": 1}
data: {"chunk": 2}
...
```

失败时输出 JSON 错误信息：

```json
{"success": false, "error": "connection timeout"}
```

## 调用示例

### Python

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
    "config_path": "./config.json"
}

proc = subprocess.Popen(
    ["./tlsRequester.exe"],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE
)

stdout, stderr = proc.communicate(json.dumps(request).encode())

# 解析响应
response = stdout.decode()
if response.startswith("HTTP/"):
    # 成功，解析 HTTP 响应
    print(response)
else:
    # 失败，解析错误信息
    error = json.loads(response)
    print(f"Error: {error['error']}")
```

### Node.js

```javascript
const { spawn } = require('child_process');

const request = {
  method: 'POST',
  url: 'https://api.openai.com/v1/chat/completions',
  headers: {
    'Authorization': 'Bearer xxx',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ stream: true, messages: [...] }),
  config_path: './config.json'
};

const proc = spawn('./tlsRequester.exe');

proc.stdin.write(JSON.stringify(request));
proc.stdin.end();

proc.stdout.on('data', (data) => {
  console.log(data.toString());
});

proc.stderr.on('data', (data) => {
  console.error(data.toString());
});
```

## 特性

- ✅ 完全自定义 TLS 指纹
- ✅ 保证 cipher 和 extension 顺序
- ✅ 支持 GREASE（自动插入到合适位置）
- ✅ 支持 HTTP/2 和 HTTP/1.1
- ✅ 支持 HTTP 和 SOCKS5 代理
- ✅ 流式传输（适合 AI 对话）
- ✅ 直接转发响应（不做任何处理）
- ✅ 子进程调用（语言无关）
- ✅ 超时控制（连接超时和读取超时）
