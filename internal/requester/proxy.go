package requester

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"fingerPrintRequester/internal/config"

	"golang.org/x/net/proxy"
)

func DialWithProxy(addr string, cfg *config.Config) (net.Conn, error) {
	var dialer proxy.Dialer = &net.Dialer{
		Timeout: time.Duration(cfg.Timeout.Connect) * time.Second,
	}

	if cfg.Proxy.Enabled {
		proxyURL, err := url.Parse(cfg.Proxy.URL)
		if err != nil {
			return nil, err
		}
		if cfg.Proxy.Type == "socks5" {
			dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, nil, dialer.(*net.Dialer))
			if err != nil {
				return nil, err
			}
		}
	}

	var conn net.Conn
	var err error

	if cfg.Proxy.Enabled && cfg.Proxy.Type == "http" {
		proxyURL, _ := url.Parse(cfg.Proxy.URL)
		conn, err = dialer.Dial("tcp", proxyURL.Host)
		if err != nil {
			return nil, err
		}
		connectReq := &http.Request{
			Method: "CONNECT",
			URL:    &url.URL{Host: addr},
			Host:   addr,
			Header: make(http.Header),
		}
		connectReq.Write(conn)
		br := bufio.NewReader(conn)
		resp, err := http.ReadResponse(br, connectReq)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if resp.StatusCode != 200 {
			conn.Close()
			return nil, fmt.Errorf("proxy connect failed: %s", resp.Status)
		}
	} else {
		conn, err = dialer.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}
