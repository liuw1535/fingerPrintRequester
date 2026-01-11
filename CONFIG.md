# 配置文件说明

## 基本结构

```json
{
  "timeout": { ... },
  "proxy": { ... },
  "fingerprint": { ... }
}
```

## 超时配置 (timeout)

```json
"timeout": {
  "connect": 30,  // 连接超时（秒）
  "read": 60      // 读取超时（秒）
}
```

## 代理配置 (proxy)

```json
"proxy": {
  "enabled": false,              // 是否启用代理
  "type": "http",                // 代理类型: "http" 或 "socks5"
  "url": "http://127.0.0.1:7890" // 代理地址
}
```

## TLS 指纹配置 (fingerprint)

### 基本参数

```json
"fingerprint": {
  "tls_version_min": "0x0303",  // 最小 TLS 版本 (0x0303=TLS1.2, 0x0304=TLS1.3)
  "tls_version_max": "0x0304",  // 最大 TLS 版本
  "http2": true,                // 是否使用 HTTP/2
  "grease": true,               // 是否启用 GREASE
  "compression_methods": [0],   // 压缩方法 (通常为 [0])
  "ciphers": [...],             // 密码套件列表
  "extensions": [...]           // 扩展列表
}
```

### 支持的密码套件 (ciphers)

**TLS 1.3 密码套件:**
- `TLS_AES_128_GCM_SHA256`
- `TLS_AES_256_GCM_SHA384`
- `TLS_CHACHA20_POLY1305_SHA256`

**TLS 1.2 ECDHE 密码套件:**
- `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`
- `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`
- `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`
- `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`
- `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256`
- `TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256`
- `TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256`
- `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA`
- `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`
- `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA`
- `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`

**TLS 1.2 RSA 密码套件:**
- `TLS_RSA_WITH_AES_128_GCM_SHA256`
- `TLS_RSA_WITH_AES_256_GCM_SHA384`
- `TLS_RSA_WITH_AES_128_CBC_SHA`
- `TLS_RSA_WITH_AES_256_CBC_SHA`

### 支持的扩展 (extensions)

#### 1. server_name (SNI)
自动使用目标 URL 的主机名
```json
{"name": "server_name"}
```

#### 2. extended_master_secret
扩展主密钥
```json
{"name": "extended_master_secret"}
```

#### 3. renegotiation_info
重协商信息
```json
{"name": "renegotiation_info"}
```

#### 4. supported_groups (椭圆曲线)
```json
{
  "name": "supported_groups",
  "data": {
    "curves": ["X25519", "CurveP256", "CurveP384", "CurveP521"]
  }
}
```

**支持的曲线:**
- `X25519` / `CurveP256` / `CurveP384` / `CurveP521`
- `secp256r1` / `secp384r1` / `secp521r1`
- `X25519MLKEM768` (混合后量子)
- `SecP256r1MLKEM768` (混合后量子)
- `SecP384r1MLKEM1024` (混合后量子)

#### 5. ec_point_formats
椭圆曲线点格式
```json
{
  "name": "ec_point_formats",
  "data": {
    "formats": [0]  // 0=uncompressed
  }
}
```

#### 6. session_ticket
会话票据
```json
{"name": "session_ticket"}
```

#### 7. application_layer_protocol_negotiation (ALPN)
```json
{
  "name": "application_layer_protocol_negotiation",
  "data": {
    "protocols": ["h2", "http/1.1"]
  }
}
```

#### 8. status_request (OCSP Stapling)
```json
{"name": "status_request"}
```

#### 9. signature_algorithms
签名算法
```json
{
  "name": "signature_algorithms",
  "data": {
    "algorithms": [
      "ecdsa_secp256r1_sha256",
      "rsa_pss_rsae_sha256",
      "rsa_pkcs1_sha256",
      "ecdsa_secp384r1_sha384",
      "rsa_pss_rsae_sha384",
      "rsa_pkcs1_sha384",
      "rsa_pss_rsae_sha512",
      "rsa_pkcs1_sha512",
      "rsa_pkcs1_sha1",
      "ecdsa_sha1",
      "ed25519"
    ]
  }
}
```

**支持的算法:**
- ECDSA: `ECDSAWithP256AndSHA256`, `ECDSAWithP384AndSHA384`, `ECDSAWithP521AndSHA512`, `ECDSAWithSHA1`
- RSA PSS: `PSSWithSHA256`, `PSSWithSHA384`, `PSSWithSHA512`
- RSA PKCS1: `PKCS1WithSHA256`, `PKCS1WithSHA384`, `PKCS1WithSHA512`, `PKCS1WithSHA1`
- EdDSA: `Ed25519`
- 别名格式: `ecdsa_secp256r1_sha256`, `rsa_pss_rsae_sha256`, `rsa_pkcs1_sha256` 等

#### 10. signature_algorithms_cert
证书签名算法（格式同 signature_algorithms）
```json
{
  "name": "signature_algorithms_cert",
  "data": {
    "algorithms": ["ecdsa_secp256r1_sha256", "rsa_pss_rsae_sha256"]
  }
}
```

#### 11. signed_certificate_timestamp (SCT)
```json
{"name": "signed_certificate_timestamp"}
```

#### 12. key_share
密钥共享
```json
{
  "name": "key_share",
  "data": {
    "groups": ["X25519", "CurveP256"]
  }
}
```

#### 13. psk_key_exchange_modes
PSK 密钥交换模式
```json
{
  "name": "psk_key_exchange_modes",
  "data": {
    "modes": [1]  // 1=psk_dhe_ke
  }
}
```

#### 14. supported_versions
支持的 TLS 版本
```json
{
  "name": "supported_versions",
  "data": {
    "versions": ["0x0304", "0x0303"]  // 0x0304=TLS1.3, 0x0303=TLS1.2
  }
}
```

#### 15. padding
填充扩展
```json
{
  "name": "padding",
  "data": {
    "length": 75  // 填充长度
  }
}
```

#### 16. compress_certificate
证书压缩
```json
{
  "name": "compress_certificate",
  "data": {
    "algorithms": [2]  // 2=brotli
  }
}
```

#### 17. application_settings (ALPS)
```json
{
  "name": "application_settings",
  "data": {
    "protocols": ["h2"]
  }
}
```

#### 18. encrypted_client_hello (ECH)
加密客户端 Hello（使用 GREASE ECH）
```json
{
  "name": "encrypted_client_hello",
  "data": {
    "payload_lengths": [128, 160, 192, 224]  // 可选：payload 长度候选列表，随机选择一个
  }
}
```

**可选参数：**
- `payload_lengths`: payload 长度数组，utls 会随机选择（默认 [128]）
- `payload_length`: 单个 payload 长度（如果不使用数组）
- `cipher_suites`: HPKE 密码套件数组（高级用法）
  ```json
  "cipher_suites": [
    {"kdf_id": 1, "aead_id": 1},  // HKDF-SHA256 + AES-128-GCM
    {"kdf_id": 1, "aead_id": 2}   // HKDF-SHA256 + AES-256-GCM
  ]
  ```

**说明：**
- 使用 utls 内置的 `GREASEEncryptedClientHelloExtension`
- 自动生成符合 RFC 标准的 ECH 结构
- Config ID、Encapsulated Key、Payload 均自动生成
- 模拟真实浏览器的 ECH GREASE 行为

#### 19. pre_shared_key (PSK)
预共享密钥（必须放在最后）
```json
{
  "name": "pre_shared_key",
  "data": {
    "identity_length": 138,  // 身份标识长度
    "binder_length": 32      // 绑定器长度
  }
}
```

## 扩展顺序说明

1. **扩展顺序很重要**，会影响 TLS 指纹
2. **GREASE** 启用时会自动插入到合适位置
3. **pre_shared_key** 必须是最后一个扩展
4. 参考真实浏览器的扩展顺序以获得最佳兼容性

## 完整示例

参考 `config-chrome141.json` 查看 Chrome 141 的完整配置示例。

## 提示

- 使用浏览器开发者工具或 Wireshark 抓包分析真实浏览器的 TLS 指纹
- 访问 https://tls.peet.ws/api/all 查看当前浏览器的 TLS 指纹
- 密码套件和扩展的顺序会影响指纹识别
- `encrypted_client_hello` 使用 utls 的 GREASE ECH 实现，自动生成符合标准的数据
- 后量子密码学曲线（MLKEM）需要服务器支持才能正常工作
