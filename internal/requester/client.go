package requester

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/fingerprint"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

func MakeRequest(req *config.Request, cfg *config.Config) error {
	spec, err := fingerprint.Build(&cfg.Fingerprint, req.URL)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return err
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

	conn, err := DialWithProxy(addr, cfg)
	if err != nil {
		return err
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
		return ForwardResponse(resp, conn)
	} else {
		// Use HTTP/1.1
		if err := httpReq.Write(conn); err != nil {
			conn.Close()
			return err
		}
		br := bufio.NewReader(conn)
		resp, err := http.ReadResponse(br, httpReq)
		if err != nil {
			conn.Close()
			return err
		}
		return ForwardResponse(resp, conn)
	}
}
