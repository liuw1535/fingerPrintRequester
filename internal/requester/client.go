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

	var negotiatedProtocol string
	// TLS handshake with timeout
	if parsedURL.Scheme == "https" {
		conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.Timeout.Read) * time.Second))
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
		
		// Get negotiated protocol from ALPN
		negotiatedProtocol = uConn.ConnectionState().NegotiatedProtocol
		conn = uConn
	}

	// Clear all timeouts for streaming
	conn.SetReadDeadline(time.Time{})
	conn.SetWriteDeadline(time.Time{})

	// Send HTTP request
	httpReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
	if err != nil {
		conn.Close()
		return err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Check if HTTP/2 should be used based on config AND ALPN negotiation
	// Only use HTTP/2 if it's enabled in config AND server negotiated "h2" via ALPN
	useHTTP2 := cfg.Fingerprint.HTTP2 && negotiatedProtocol == "h2"
	
	if useHTTP2 {
		// Use HTTP/2
		transport := &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return conn, nil
			},
			// Don't allow fallback to HTTP/1.1 if we expect HTTP/2
			AllowHTTP: false,
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
