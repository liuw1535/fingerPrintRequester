package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/requester"
)

func main() {
	// Read request from stdin
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
