package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/requester"
)

func main() {
	// Check if running in curl mode (has command line args)
	if len(os.Args) > 1 {
		runCurlMode()
		return
	}

	// Original stdin JSON mode
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		outputError("INPUT_ERROR", fmt.Sprintf("failed to read stdin: %v", err), 1)
	}

	var req config.Request
	if err := json.Unmarshal(input, &req); err != nil {
		outputError("INPUT_ERROR", fmt.Sprintf("failed to parse request: %v", err), 1)
	}

	// Load config
	cfg, err := config.LoadConfig(req.ConfigPath)
	if err != nil {
		outputError("CONFIG_ERROR", fmt.Sprintf("failed to load config: %v", err), 4)
	}

	// Override config with request parameters
	if req.Timeout != nil {
		if req.Timeout.Connect > 0 {
			cfg.Timeout.Connect = req.Timeout.Connect
		}
		if req.Timeout.Read > 0 {
			cfg.Timeout.Read = req.Timeout.Read
		}
	}
	if req.Proxy != nil {
		cfg.Proxy = *req.Proxy
	}
	if req.DNS != nil {
		cfg.DNS = *req.DNS
	}

	// Make request
	if err := requester.MakeRequest(&req, cfg); err != nil {
		errMsg := err.Error()
		code := 2 // Default: network error
		errType := "NETWORK_ERROR"
		
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline") {
			code = 3
			errType = "TIMEOUT_ERROR"
		}
		
		outputError(errType, errMsg, code)
	}
}

const version = "1.0.0"

func runCurlMode() {
	var (
		method     = flag.String("X", "GET", "HTTP method")
		headerArgs = flag.String("H", "", "Headers (JSON object or key:value)")
		data       = flag.String("d", "", "Request body")
		configPath = flag.String("c", "config.json", "Config file path")
		proxy      = flag.String("x", "", "Proxy URL")
		showVersion = flag.Bool("v", false, "Show version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("TLS Requester v%s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "TLS Requester v%s\n", version)
		fmt.Fprintln(os.Stderr, "Usage: tlsRequester [options] <url>")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)
	req := config.Request{
		Method:     *method,
		URL:        url,
		Headers:    make(map[string]string),
		Body:       *data,
		ConfigPath: *configPath,
	}

	// Parse headers
	if *headerArgs != "" {
		if strings.HasPrefix(*headerArgs, "{") {
			json.Unmarshal([]byte(*headerArgs), &req.Headers)
		} else {
			parts := strings.SplitN(*headerArgs, ":", 2)
			if len(parts) == 2 {
				req.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	// Load config
	cfg, err := config.LoadConfig(req.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(4)
	}

	// Set proxy if specified
	if *proxy != "" {
		proxyType := "http"
		if strings.HasPrefix(*proxy, "socks") {
			proxyType = "socks5"
		}
		cfg.Proxy = config.ProxyConfig{
			Enabled: true,
			Type:    proxyType,
			URL:     *proxy,
		}
	}

	// Make request
	if err := requester.MakeRequest(&req, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}

func outputError(errType, msg string, exitCode int) {
	errResp := map[string]interface{}{
		"success":    false,
		"error":      msg,
		"error_type": errType,
	}
	data, _ := json.Marshal(errResp)
	fmt.Fprintln(os.Stderr, string(data))
	os.Exit(exitCode)
}
