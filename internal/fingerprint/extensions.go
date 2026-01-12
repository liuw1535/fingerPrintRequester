package fingerprint

import (
	"crypto/rand"
	"fmt"

	"fingerPrintRequester/internal/config"
	"fingerPrintRequester/internal/utils"

	utls "github.com/refraction-networking/utls"
)

func BuildExtension(cfg config.ExtensionConfig, serverName string) (utls.TLSExtension, error) {
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
				if curveID, exists := CurveMap[name.(string)]; exists {
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
				if algo, exists := SignatureAlgorithmMap[name.(string)]; exists {
					algorithms = append(algorithms, algo)
				}
			}
		}
		return &utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: algorithms}, nil
	case "signature_algorithms_cert":
		algorithms := []utls.SignatureScheme{}
		if algoNames, ok := cfg.Data["algorithms"].([]interface{}); ok {
			for _, name := range algoNames {
				if algo, exists := SignatureAlgorithmMap[name.(string)]; exists {
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
				if curveID, exists := CurveMap[name.(string)]; exists {
					groups = append(groups, curveID)
				}
			}
		}
		keyShares := []utls.KeyShare{}
		for _, group := range groups {
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
				versions[i] = utils.ParseHex(v.(string))
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
					ObfuscatedTicketAge: utils.GenerateRandomObfuscatedTicketAge(),
				},
			},
			Binders: [][]byte{binder},
		}, nil
	case "encrypted_client_hello":
		ext := &utls.GREASEEncryptedClientHelloExtension{}
		
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
