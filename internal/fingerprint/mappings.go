package fingerprint

import (
	utls "github.com/refraction-networking/utls"
)

var CipherMap = map[string]uint16{
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

var CurveMap = map[string]utls.CurveID{
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

var SignatureAlgorithmMap = map[string]utls.SignatureScheme{
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
