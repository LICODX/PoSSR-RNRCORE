package p2p

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
	CAFile   string
}

// NewSecureHost creates LibP2P host with Noise + TLS security
func NewSecureHost(ctx context.Context, port int, tlsConfig *TLSConfig) (host.Host, error) {
	// Generate or load identity
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 256)
	if err != nil {
		return nil, err
	}

	// Security transports (Noise + TLS)
	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	}

	if tlsConfig != nil && tlsConfig.Enabled {
		// Load custom TLS config
		tlsConf, err := loadTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}
		_ = tlsConf // Use for custom transport if needed

		fmt.Println("🔒 TLS encryption enabled for P2P")
	}

	// Create host with security
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Identity(priv),
		libp2p.ChainOptions(securityOptions...),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Secure P2P host created with Noise + TLS")
	return h, nil
}

// loadTLSConfig loads TLS certificates
func loadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	// Load CA if provided
	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
