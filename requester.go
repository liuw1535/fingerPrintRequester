package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
)

func MakeRequest(req *Request, cfg *Config) error {
	spec, err := BuildFingerprint(&cfg.Fingerprint, req.URL)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return err
	}

	// Build dialer with proxy support
	var dialer proxy.Dialer = &net.Dialer{
		Timeout: time.Duration(cfg.Timeout.Connect) * time.Second,
	}

	if cfg.Proxy.Enabled {
		proxyURL, err := url.Parse(cfg.Proxy.URL)
		if err != nil {
			return err
		}
		if cfg.Proxy.Type == "socks5" {
			dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, nil, dialer.(*net.Dialer))
			if err != nil {
				return err
			}
		}
	}

	// Dial connection
	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	addr := net.JoinHostPort(parsedURL.Hostname(), port)

	var conn net.Conn
	if cfg.Proxy.Enabled && cfg.Proxy.Type == "http" {
		proxyURL, _ := url.Parse(cfg.Proxy.URL)
		conn, err = dialer.Dial("tcp", proxyURL.Host)
		if err != nil {
			return err
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
			return err
		}
		if resp.StatusCode != 200 {
			conn.Close()
			return fmt.Errorf("proxy connect failed: %s", resp.Status)
		}
	} else {
		conn, err = dialer.Dial("tcp", addr)
		if err != nil {
			return err
		}
	}

	// Apply read timeout for handshake
	conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.Timeout.Read) * time.Second))

	// TLS handshake
	if parsedURL.Scheme == "https" {
		tlsConfig := &utls.Config{
			ServerName:         parsedURL.Hostname(),
			InsecureSkipVerify: true,
		}
		uConn := utls.UClient(conn, tlsConfig, utls.HelloCustom)
		if err := uConn.ApplyPreset(spec); err != nil {
			conn.Close()
			return err
		}
		if err := uConn.Handshake(); err != nil {
			conn.Close()
			return err
		}
		conn = uConn
	}

	// Send HTTP request
	httpReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
	if err != nil {
		conn.Close()
		return err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	if cfg.Fingerprint.HTTP2 {
		// Use HTTP/2
		transport := &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return conn, nil
			},
			AllowHTTP: true,
		}
		client := &http.Client{Transport: transport}
		resp, err := client.Do(httpReq)
		if err != nil {
			conn.Close()
			return err
		}
		return forwardResponse(resp, conn)
	} else {
		// Use HTTP/1.1
		if err := httpReq.Write(conn); err != nil {
			conn.Close()
			return err
		}
		// Read response with timeout
		br := bufio.NewReader(conn)
		resp, err := http.ReadResponse(br, httpReq)
		if err != nil {
			conn.Close()
			return err
		}
		return forwardResponse(resp, conn)
	}
}

func forwardResponse(resp *http.Response, conn net.Conn) error {
	defer conn.Close()
	
	// Cancel read timeout for streaming
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetReadDeadline(time.Time{})
	} else if utlsConn, ok := conn.(*utls.UConn); ok {
		utlsConn.SetReadDeadline(time.Time{})
	}

	// Write status line
	fmt.Fprintf(os.Stdout, "HTTP/%d.%d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)

	// Write headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			fmt.Fprintf(os.Stdout, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprintf(os.Stdout, "\r\n")
	os.Stdout.Sync() // Flush headers immediately

	// Stream body chunk by chunk
	buf := make([]byte, 8192)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
			os.Stdout.Sync() // Flush each chunk immediately
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			resp.Body.Close()
			return err
		}
	}
	resp.Body.Close()
	return nil
}

