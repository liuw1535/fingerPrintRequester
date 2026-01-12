# 项目重构说明

## 新的目录结构

```
fingerPrintRequester/
├── cmd/
│   └── tlsRequester/
│       └── main.go              # 程序入口，处理 stdin/stdout
├── internal/
│   ├── config/
│   │   ├── types.go             # 配置和请求类型定义
│   │   └── config.go            # 配置加载函数
│   ├── fingerprint/
│   │   ├── mappings.go          # Cipher/Curve/Signature 映射表
│   │   ├── extensions.go        # TLS 扩展构建逻辑
│   │   └── builder.go           # 指纹构建主函数
│   ├── requester/
│   │   ├── proxy.go             # 代理连接处理
│   │   ├── response.go          # HTTP 响应转发
│   │   └── client.go            # HTTP 请求客户端
│   └── utils/
│       └── utils.go             # 工具函数
├── bin/                         # 编译产物目录
├── config.json                  # 配置文件
├── go.mod
├── build.bat                    # Windows 构建脚本
└── build.sh                     # Linux 构建脚本
```

## 模块职责

### cmd/tlsRequester
- 程序入口点
- 处理 stdin 输入解析
- 处理 stdout 输出和错误信息
- 调用 internal 包的功能

### internal/config
- **types.go**: 定义所有配置和请求相关的数据结构
- **config.go**: 提供配置文件加载功能

### internal/fingerprint
- **mappings.go**: 存储 TLS cipher、curve、signature 的映射表
- **extensions.go**: 构建各种 TLS 扩展（SNI、ALPN、KeyShare 等）
- **builder.go**: 组装完整的 TLS 指纹

### internal/requester
- **proxy.go**: 处理 HTTP/SOCKS5 代理连接
- **response.go**: 转发 HTTP 响应到 stdout（支持流式传输）
- **client.go**: 核心请求逻辑（TLS 握手、HTTP 请求）

### internal/utils
- **utils.go**: 通用工具函数（GREASE 生成、Hex 解析等）

## 编译方法

### Windows
```bash
build.bat
```

### Linux/macOS
```bash
chmod +x build.sh
./build.sh
```

### 手动编译
```bash
go build -o bin/tlsRequester.exe ./cmd/tlsRequester
```

## 重构优势

1. **清晰的分层架构**: cmd 层处理 I/O，internal 层处理业务逻辑
2. **模块化设计**: 每个包职责单一，易于测试和维护
3. **符合 Go 规范**: 使用标准的 cmd 和 internal 目录结构
4. **避免循环依赖**: 依赖关系清晰明确
5. **易于扩展**: 新增功能只需修改对应的包
6. **代码复用**: 各模块可独立使用和测试

## 迁移说明

旧文件已保留在根目录，可以安全删除：
- main.go
- config.go
- fingerprint.go
- extensions.go
- requester.go
- utils.go

新代码已完全替代旧代码的功能，并保持 API 兼容。
