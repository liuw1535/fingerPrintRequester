package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/requester"
)

func main() {
	// Read request from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		outputError(fmt.Sprintf("failed to read stdin: %v", err))
		os.Exit(1)
	}

	var req config.Request
	if err := json.Unmarshal(input, &req); err != nil {
		outputError(fmt.Sprintf("failed to parse request: %v", err))
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig(req.ConfigPath)
	if err != nil {
		outputError(fmt.Sprintf("failed to load config: %v", err))
		os.Exit(1)
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

	// Make request
	if err := requester.MakeRequest(&req, cfg); err != nil {
		outputError(fmt.Sprintf("request failed: %v", err))
		os.Exit(1)
	}
}

func outputError(msg string) {
	errResp := map[string]interface{}{
		"success": false,
		"error":   msg,
	}
	data, _ := json.Marshal(errResp)
	fmt.Println(string(data))
}
