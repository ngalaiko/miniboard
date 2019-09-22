package api

// TLSConfig contains ssl certificates.
type TLSConfig struct {
	CertPath string
	KeyPath  string
}

func (cfg *TLSConfig) valid() bool {
	return cfg.CertPath != "" || cfg.KeyPath != ""
}
