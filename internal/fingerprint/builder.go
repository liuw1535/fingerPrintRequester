package fingerprint

import (
	"net/url"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/utils"

	utls "github.com/refraction-networking/utls"
)

func Build(cfg *config.FingerprintConfig, targetURL string) (*utls.ClientHelloSpec, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	serverName := parsedURL.Hostname()

	spec := &utls.ClientHelloSpec{
		TLSVersMin:         utils.ParseHex(cfg.TLSVersionMin),
		TLSVersMax:         utils.ParseHex(cfg.TLSVersionMax),
		CompressionMethods: cfg.CompressionMethods,
	}

	// Build ciphers
	ciphers := []uint16{}
	if cfg.GREASE {
		ciphers = append(ciphers, utils.GenerateGREASEValue())
	}
	for _, cipherName := range cfg.Ciphers {
		if cipher, ok := CipherMap[cipherName]; ok {
			ciphers = append(ciphers, cipher)
		}
	}
	spec.CipherSuites = ciphers

	// Build extensions
	extensions := []utls.TLSExtension{}
	if cfg.GREASE {
		extensions = append(extensions, &utls.UtlsGREASEExtension{})
	}
	
	for _, extCfg := range cfg.Extensions {
		ext, err := BuildExtension(extCfg, serverName)
		if err != nil {
			continue
		}
		extensions = append(extensions, ext)
	}
	
	// Insert GREASE before last extension (if GREASE enabled and pre_shared_key is last)
	if cfg.GREASE && len(extensions) > 0 {
		if len(cfg.Extensions) > 0 && cfg.Extensions[len(cfg.Extensions)-1].Name == "pre_shared_key" {
			extensions = append(extensions[:len(extensions)-1], &utls.UtlsGREASEExtension{}, extensions[len(extensions)-1])
		}
	}
	
	spec.Extensions = extensions
	return spec, nil
}
