package fingerprint

import (
	utls "github.com/refraction-networking/utls"
)

var CipherMap = map[string]uint16{
	// TLS 1.3 cipher suites
	"TLS_AES_128_GCM_SHA256":       utls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":       utls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256": utls.TLS_CHACHA20_POLY1305_SHA256,

	// ECDHE with ECDSA
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384":       0xc024, // Not in utls, use raw value

	// ECDHE with RSA
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":       utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":       utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256": utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":          utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":          utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":       utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384":       0xc028, // Not in utls, use raw value

	// DHE with RSA
	"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256":       0x009e,
	"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384":       0x009f,
	"TLS_DHE_RSA_WITH_CHACHA20_POLY1305_SHA256": 0xccaa,
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA":          0x0033,
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA":          0x0039,
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256":       0x0067,
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256":       0x006b,
	"TLS_DHE_RSA_WITH_AES_128_CCM":              0xc09e,
	"TLS_DHE_RSA_WITH_AES_256_CCM":              0xc09f,
	"TLS_DHE_RSA_WITH_ARIA_128_GCM_SHA256":      0xc052,
	"TLS_DHE_RSA_WITH_ARIA_256_GCM_SHA384":      0xc053,

	// DHE with DSS
	"TLS_DHE_DSS_WITH_AES_128_GCM_SHA256":  0x00a2,
	"TLS_DHE_DSS_WITH_AES_256_GCM_SHA384":  0x00a3,
	"TLS_DHE_DSS_WITH_AES_128_CBC_SHA":     0x0032,
	"TLS_DHE_DSS_WITH_AES_256_CBC_SHA":     0x0038,
	"TLS_DHE_DSS_WITH_AES_128_CBC_SHA256":  0x0040,
	"TLS_DHE_DSS_WITH_AES_256_CBC_SHA256":  0x006a,
	"TLS_DHE_DSS_WITH_ARIA_128_GCM_SHA256": 0xc056,
	"TLS_DHE_DSS_WITH_ARIA_256_GCM_SHA384": 0xc057,

	// RSA
	"TLS_RSA_WITH_AES_128_GCM_SHA256":  utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":  utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_RSA_WITH_AES_128_CBC_SHA":     utls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":     utls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":  0x003c,
	"TLS_RSA_WITH_AES_256_CBC_SHA256":  0x003d,
	"TLS_RSA_WITH_AES_128_CCM":         0xc09c,
	"TLS_RSA_WITH_AES_256_CCM":         0xc09d,
	"TLS_RSA_WITH_ARIA_128_GCM_SHA256": 0xc050,
	"TLS_RSA_WITH_ARIA_256_GCM_SHA384": 0xc051,

	// ECDHE ECDSA with CCM and ARIA
	"TLS_ECDHE_ECDSA_WITH_AES_128_CCM":         0xc0ac,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CCM":         0xc0ad,
	"TLS_ECDHE_ECDSA_WITH_ARIA_128_GCM_SHA256": 0xc05c,
	"TLS_ECDHE_ECDSA_WITH_ARIA_256_GCM_SHA384": 0xc05d,

	// ECDHE RSA with ARIA
	"TLS_ECDHE_RSA_WITH_ARIA_128_GCM_SHA256": 0xc060,
	"TLS_ECDHE_RSA_WITH_ARIA_256_GCM_SHA384": 0xc061,
}

var CurveMap = map[string]utls.CurveID{
	// Standard curves
	"X25519":    utls.X25519,
	"x25519":    utls.X25519,
	"CurveP256": utls.CurveP256,
	"CurveP384": utls.CurveP384,
	"CurveP521": utls.CurveP521,
	"secp256r1": utls.CurveP256,
	"secp384r1": utls.CurveP384,
	"secp521r1": utls.CurveP521,

	// X448 curve
	"X448": utls.CurveID(0x001e),
	"x448": utls.CurveID(0x001e),

	// Post-quantum hybrid curves (ML-KEM)
	"X25519MLKEM768":     utls.CurveID(0x11ec),
	"SecP256r1MLKEM768":  utls.CurveID(0x11eb),
	"SecP384r1MLKEM1024": utls.CurveID(0x11ed),

	// FFDHE groups (Finite Field Diffie-Hellman Ephemeral)
	"ffdhe2048": utls.CurveID(0x0100),
	"ffdhe3072": utls.CurveID(0x0101),
	"ffdhe4096": utls.CurveID(0x0102),
	"ffdhe6144": utls.CurveID(0x0103),
	"ffdhe8192": utls.CurveID(0x0104),
}

var SignatureAlgorithmMap = map[string]utls.SignatureScheme{
	// Standard ECDSA algorithms
	"ECDSAWithP256AndSHA256": utls.ECDSAWithP256AndSHA256,
	"ECDSAWithP384AndSHA384": utls.ECDSAWithP384AndSHA384,
	"ECDSAWithP521AndSHA512": utls.ECDSAWithP521AndSHA512,
	"ECDSAWithSHA1":          utls.ECDSAWithSHA1,

	// RSA-PSS algorithms (rsae - using RSA keys in X.509 certificates)
	"PSSWithSHA256": utls.PSSWithSHA256,
	"PSSWithSHA384": utls.PSSWithSHA384,
	"PSSWithSHA512": utls.PSSWithSHA512,

	// RSA PKCS#1 algorithms
	"PKCS1WithSHA256": utls.PKCS1WithSHA256,
	"PKCS1WithSHA384": utls.PKCS1WithSHA384,
	"PKCS1WithSHA512": utls.PKCS1WithSHA512,
	"PKCS1WithSHA1":   utls.PKCS1WithSHA1,

	// EdDSA algorithms
	"Ed25519": utls.Ed25519,
	"ed25519": utls.Ed25519,
	"Ed448":   utls.SignatureScheme(0x0808),
	"ed448":   utls.SignatureScheme(0x0808),

	// TLS 1.3 style names - ECDSA
	"ecdsa_secp256r1_sha256": utls.ECDSAWithP256AndSHA256,
	"ecdsa_secp384r1_sha384": utls.ECDSAWithP384AndSHA384,
	"ecdsa_secp521r1_sha512": utls.ECDSAWithP521AndSHA512,
	"ecdsa_sha1":             utls.ECDSAWithSHA1,

	// TLS 1.3 style names - RSA-PSS (rsae - RSA keys in X.509)
	"rsa_pss_rsae_sha256": utls.PSSWithSHA256,
	"rsa_pss_rsae_sha384": utls.PSSWithSHA384,
	"rsa_pss_rsae_sha512": utls.PSSWithSHA512,

	// TLS 1.3 style names - RSA-PSS (pss - PSS keys in X.509)
	"rsa_pss_pss_sha256": utls.SignatureScheme(0x0809),
	"rsa_pss_pss_sha384": utls.SignatureScheme(0x080a),
	"rsa_pss_pss_sha512": utls.SignatureScheme(0x080b),

	// TLS 1.3 style names - RSA PKCS#1
	"rsa_pkcs1_sha256": utls.PKCS1WithSHA256,
	"rsa_pkcs1_sha384": utls.PKCS1WithSHA384,
	"rsa_pkcs1_sha512": utls.PKCS1WithSHA512,
	"rsa_pkcs1_sha1":   utls.PKCS1WithSHA1,

	// SHA224 algorithms (legacy)
	"sha224_ecdsa": utls.SignatureScheme(0x0303),
	"sha224_rsa":   utls.SignatureScheme(0x0301),
	"sha224_dsa":   utls.SignatureScheme(0x0302),

	// DSA algorithms
	"sha256_dsa": utls.SignatureScheme(0x0402),
	"sha384_dsa": utls.SignatureScheme(0x0502),
	"sha512_dsa": utls.SignatureScheme(0x0602),

	// Brainpool curves (TLS 1.3)
	"ecdsa_brainpoolP256r1tls13_sha256": utls.SignatureScheme(0x081a),
	"ecdsa_brainpoolP384r1tls13_sha384": utls.SignatureScheme(0x081b),
	"ecdsa_brainpoolP512r1tls13_sha512": utls.SignatureScheme(0x081c),

	// ML-DSA (post-quantum) algorithms - TLS 1.3
	"mldsa44":          utls.SignatureScheme(0x0905),
	"mldsa65":          utls.SignatureScheme(0x0906),
	"mldsa87":          utls.SignatureScheme(0x0907),
	"mldsa44_rsa2048":  utls.SignatureScheme(0x0904),
	"mldsa65_ecdsa256": utls.SignatureScheme(0x0905),
	"mldsa87_ecdsa384": utls.SignatureScheme(0x0906),
}
