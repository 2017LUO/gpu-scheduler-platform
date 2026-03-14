package bootstrap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	appcfg "gpu-scheduler-platform/internal/config"
)

// BuildClientTLSConfig 根据配置构造客户端 TLS 配置。
// 支持：
// 1. 自定义 CA
// 2. 可选双向认证（cert/key）
// 3. 自定义 server_name
// 4. insecure_skip_verify（仅测试环境建议）
func BuildClientTLSConfig(cfg appcfg.TLSConfig) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		ServerName:         cfg.ServerName,
	}

	if cfg.CAFile != "" {
		caPEM, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("read ca file %q: %w", cfg.CAFile, err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("append ca certs from %q failed", cfg.CAFile)
		}
		tlsCfg.RootCAs = pool
	}

	// 若配置了 cert/key，则启用客户端证书（mTLS）
	if cfg.CertFile != "" || cfg.KeyFile != "" {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			return nil, fmt.Errorf("both cert_file and key_file are required for client tls cert")
		}

		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("load x509 key pair cert=%q key=%q: %w", cfg.CertFile, cfg.KeyFile, err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return tlsCfg, nil
}
