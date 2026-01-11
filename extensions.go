package main

import (
	"crypto/rand"
	"fmt"
	utls "github.com/refraction-networking/utls"
)

var cipherMap = map[string]uint16{
	"TLS_AES_128_GCM_SHA256":                      utls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":                      utls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256":                utls.TLS_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":       utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":     utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":       utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":     utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":       utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":        utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":          utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":        utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":          utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":             utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":             utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_RSA_WITH_AES_128_CBC_SHA":                utls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":                utls.TLS_RSA_WITH_AES_256_CBC_SHA,
}

var curveMap = map[string]utls.CurveID{
	"X25519":              utls.X25519,
	"CurveP256":           utls.CurveP256,
	"CurveP384":           utls.CurveP384,
	"CurveP521":           utls.CurveP521,
	"secp256r1":           utls.CurveP256,
	"secp384r1":           utls.CurveP384,
	"secp521r1":           utls.CurveP521,
	"X25519MLKEM768":      utls.CurveID(0x11ec),
	"SecP256r1MLKEM768":   utls.CurveID(0x11eb),
	"SecP384r1MLKEM1024":  utls.CurveID(0x11ed),
}

var signatureAlgorithmMap = map[string]utls.SignatureScheme{
	"ECDSAWithP256AndSHA256":   utls.ECDSAWithP256AndSHA256,
	"ECDSAWithP384AndSHA384":   utls.ECDSAWithP384AndSHA384,
	"ECDSAWithP521AndSHA512":   utls.ECDSAWithP521AndSHA512,
	"PSSWithSHA256":            utls.PSSWithSHA256,
	"PSSWithSHA384":            utls.PSSWithSHA384,
	"PSSWithSHA512":            utls.PSSWithSHA512,
	"PKCS1WithSHA256":          utls.PKCS1WithSHA256,
	"PKCS1WithSHA384":          utls.PKCS1WithSHA384,
	"PKCS1WithSHA512":          utls.PKCS1WithSHA512,
	"PKCS1WithSHA1":            utls.PKCS1WithSHA1,
	"ECDSAWithSHA1":            utls.ECDSAWithSHA1,
	"Ed25519":                  utls.Ed25519,
	"ecdsa_secp256r1_sha256":   utls.ECDSAWithP256AndSHA256,
	"ecdsa_secp384r1_sha384":   utls.ECDSAWithP384AndSHA384,
	"ecdsa_secp521r1_sha512":   utls.ECDSAWithP521AndSHA512,
	"rsa_pss_rsae_sha256":      utls.PSSWithSHA256,
	"rsa_pss_rsae_sha384":      utls.PSSWithSHA384,
	"rsa_pss_rsae_sha512":      utls.PSSWithSHA512,
	"rsa_pkcs1_sha256":         utls.PKCS1WithSHA256,
	"rsa_pkcs1_sha384":         utls.PKCS1WithSHA384,
	"rsa_pkcs1_sha512":         utls.PKCS1WithSHA512,
	"rsa_pkcs1_sha1":           utls.PKCS1WithSHA1,
	"ecdsa_sha1":               utls.ECDSAWithSHA1,
	"ed25519":                  utls.Ed25519,
}

func buildExtension(cfg ExtensionConfig, serverName string) (utls.TLSExtension, error) {
	switch cfg.Name {
	case "server_name":
		return &utls.SNIExtension{ServerName: serverName}, nil
	case "extended_master_secret":
		return &utls.UtlsExtendedMasterSecretExtension{}, nil
	case "renegotiation_info":
		return &utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, nil
	case "supported_groups":
		curves := []utls.CurveID{}
		if curveNames, ok := cfg.Data["curves"].([]interface{}); ok {
			for _, name := range curveNames {
				if curveID, exists := curveMap[name.(string)]; exists {
					curves = append(curves, curveID)
				}
			}
		}
		return &utls.SupportedCurvesExtension{Curves: curves}, nil
	case "ec_point_formats":
		formats := []byte{0}
		if fmts, ok := cfg.Data["formats"].([]interface{}); ok {
			formats = make([]byte, len(fmts))
			for i, f := range fmts {
				formats[i] = byte(f.(float64))
			}
		}
		return &utls.SupportedPointsExtension{SupportedPoints: formats}, nil
	case "session_ticket":
		return &utls.SessionTicketExtension{}, nil
	case "application_layer_protocol_negotiation":
		protocols := []string{"h2", "http/1.1"}
		if protos, ok := cfg.Data["protocols"].([]interface{}); ok {
			protocols = make([]string, len(protos))
			for i, p := range protos {
				protocols[i] = p.(string)
			}
		}
		return &utls.ALPNExtension{AlpnProtocols: protocols}, nil
	case "status_request":
		return &utls.StatusRequestExtension{}, nil
	case "signature_algorithms":
		algorithms := []utls.SignatureScheme{}
		if algoNames, ok := cfg.Data["algorithms"].([]interface{}); ok {
			for _, name := range algoNames {
				if algo, exists := signatureAlgorithmMap[name.(string)]; exists {
					algorithms = append(algorithms, algo)
				}
			}
		}
		return &utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: algorithms}, nil
	case "signature_algorithms_cert":
		algorithms := []utls.SignatureScheme{}
		if algoNames, ok := cfg.Data["algorithms"].([]interface{}); ok {
			for _, name := range algoNames {
				if algo, exists := signatureAlgorithmMap[name.(string)]; exists {
					algorithms = append(algorithms, algo)
				}
			}
		}
		return &utls.SignatureAlgorithmsCertExtension{SupportedSignatureAlgorithms: algorithms}, nil
	case "signed_certificate_timestamp":
		return &utls.SCTExtension{}, nil
	case "key_share":
		groups := []utls.CurveID{utls.X25519}
		if groupNames, ok := cfg.Data["groups"].([]interface{}); ok {
			groups = make([]utls.CurveID, 0, len(groupNames))
			for _, name := range groupNames {
				if curveID, exists := curveMap[name.(string)]; exists {
					groups = append(groups, curveID)
				}
			}
		}
		keyShares := []utls.KeyShare{}
		for _, group := range groups {
			// 让 uTLS 自动生成所有曲线的密钥数据
			keyShares = append(keyShares, utls.KeyShare{Group: group})
		}
		return &utls.KeyShareExtension{KeyShares: keyShares}, nil
	case "psk_key_exchange_modes":
		modes := []uint8{utls.PskModeDHE}
		if modeList, ok := cfg.Data["modes"].([]interface{}); ok {
			modes = make([]uint8, len(modeList))
			for i, m := range modeList {
				modes[i] = uint8(m.(float64))
			}
		}
		return &utls.PSKKeyExchangeModesExtension{Modes: modes}, nil
	case "supported_versions":
		versions := []uint16{utls.VersionTLS13, utls.VersionTLS12}
		if verList, ok := cfg.Data["versions"].([]interface{}); ok {
			versions = make([]uint16, len(verList))
			for i, v := range verList {
				versions[i] = parseHex(v.(string))
			}
		}
		return &utls.SupportedVersionsExtension{Versions: versions}, nil
	case "padding":
		length := 0
		if l, ok := cfg.Data["length"].(float64); ok {
			length = int(l)
		}
		return &utls.UtlsPaddingExtension{
			GetPaddingLen: func(clientHelloUnpaddedLen int) (int, bool) {
				return length, true
			},
		}, nil
	case "compress_certificate":
		algorithms := []utls.CertCompressionAlgo{utls.CertCompressionBrotli}
		if algos, ok := cfg.Data["algorithms"].([]interface{}); ok {
			algorithms = make([]utls.CertCompressionAlgo, len(algos))
			for i, a := range algos {
				algorithms[i] = utls.CertCompressionAlgo(a.(float64))
			}
		}
		return &utls.UtlsCompressCertExtension{Algorithms: algorithms}, nil
	case "application_settings":
		protocols := []string{"h2"}
		if protos, ok := cfg.Data["protocols"].([]interface{}); ok {
			protocols = make([]string, len(protos))
			for i, p := range protos {
				protocols[i] = p.(string)
			}
		}
		return &utls.ApplicationSettingsExtension{SupportedProtocols: protocols}, nil
	case "pre_shared_key":
		identityLen := 138
		binderLen := 32
		if l, ok := cfg.Data["identity_length"].(float64); ok {
			identityLen = int(l)
		}
		if l, ok := cfg.Data["binder_length"].(float64); ok {
			binderLen = int(l)
		}
		identity := make([]byte, identityLen)
		rand.Read(identity)
		binder := make([]byte, binderLen)
		rand.Read(binder)
		return &utls.FakePreSharedKeyExtension{
			Identities: []utls.PskIdentity{
				{
					Label:               identity,
					ObfuscatedTicketAge: generateRandomObfuscatedTicketAge(),
				},
			},
			Binders: [][]byte{binder},
		}, nil
	case "encrypted_client_hello":
		// 使用 utls 内置的 GREASE ECH 扩展
		ext := &utls.GREASEEncryptedClientHelloExtension{}
		
		// 可选：配置 cipher suites
		if cipherSuites, ok := cfg.Data["cipher_suites"].([]interface{}); ok {
			ext.CandidateCipherSuites = make([]utls.HPKESymmetricCipherSuite, len(cipherSuites))
			for i, cs := range cipherSuites {
				if csMap, ok := cs.(map[string]interface{}); ok {
					ext.CandidateCipherSuites[i] = utls.HPKESymmetricCipherSuite{
						KdfId:  uint16(csMap["kdf_id"].(float64)),
						AeadId: uint16(csMap["aead_id"].(float64)),
					}
				}
			}
		}
		
		// 可选：配置 payload 长度
		if payloadLens, ok := cfg.Data["payload_lengths"].([]interface{}); ok {
			ext.CandidatePayloadLens = make([]uint16, len(payloadLens))
			for i, l := range payloadLens {
				ext.CandidatePayloadLens[i] = uint16(l.(float64))
			}
		} else if payloadLen, ok := cfg.Data["payload_length"].(float64); ok {
			ext.CandidatePayloadLens = []uint16{uint16(payloadLen)}
		}
		
		return ext, nil
	case "GREASE":
		return &utls.UtlsGREASEExtension{}, nil
	default:
		return nil, fmt.Errorf("unknown extension: %s", cfg.Name)
	}
}
