package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Timeout     TimeoutConfig     `json:"timeout"`
	Proxy       ProxyConfig       `json:"proxy"`
	Fingerprint FingerprintConfig `json:"fingerprint"`
}

type TimeoutConfig struct {
	Connect int `json:"connect"`
	Read    int `json:"read"`
}

type ProxyConfig struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	URL     string `json:"url"`
}

type FingerprintConfig struct {
	TLSVersionMin      string              `json:"tls_version_min"`
	TLSVersionMax      string              `json:"tls_version_max"`
	HTTP2              bool                `json:"http2"`
	GREASE             bool                `json:"grease"`
	Ciphers            []string            `json:"ciphers"`
	CompressionMethods []byte              `json:"compression_methods"`
	Extensions         []ExtensionConfig   `json:"extensions"`
}

type ExtensionConfig struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type Request struct {
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	ConfigPath string            `json:"config_path"`
	Timeout    *TimeoutConfig    `json:"timeout,omitempty"`
	Proxy      *ProxyConfig      `json:"proxy,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
