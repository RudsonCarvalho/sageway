// Package config holds all configuration structs for the AASG system.
// All fields are loaded from environment variables at startup via os.LookupEnv.
// Missing required variables cause an immediate fatal error -- no defaults for secrets.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GatewayConfig holds all configuration for the AASG gateway.
// Loaded via LoadGateway at startup; never modified after construction.
type GatewayConfig struct {
	Env   string
	HTTP  HTTPConfig
	L2    L2Config
	TLS   TLSConfig
	Etcd  EtcdConfig
	Vault VaultConfig
	SIEM  SIEMConfig
	OPA   OPAConfig
}

// OrchestratorConfig holds all configuration for the AASG orchestrator.
// Loaded via LoadOrchestrator at startup; never modified after construction.
type OrchestratorConfig struct {
	Env      string
	GRPCAddr string
	L2       L2Config
	TLS      TLSConfig
	Etcd     EtcdConfig
	Vault    VaultConfig
	SIEM     SIEMConfig
}

// HTTPConfig holds network address configuration for the gateway.
type HTTPConfig struct {
	Addr            string
	GRPCAddr        string
	WSAddr          string
	UpstreamTimeout time.Duration
}

// L2Config holds Control Plane timing and security parameters (Spec s5.2).
type L2Config struct {
	HeartbeatInterval time.Duration
	MissTimeout       time.Duration
	ProbeTimeout      time.Duration
	MaxProbes         int
	SessionKeyTTL     time.Duration
	SessionKeyDrain   time.Duration
	CertTTL           time.Duration
	ReplayWindowSize  int
}

// TLSConfig holds certificate paths and SPIFFE configuration.
// All paths must reference existing files; validated at startup.
type TLSConfig struct {
	CertPath   string
	KeyPath    string
	CAPath     string
	SPIFFEAddr string
}

// EtcdConfig holds etcd cluster connection parameters.
// Endpoints is a comma-separated list parsed from AASG_ETCD_ENDPOINTS.
type EtcdConfig struct {
	Endpoints   []string
	DialTimeout time.Duration
}

// VaultConfig holds Vault connection parameters.
// SecretPath references the Vault path for the L2 master secret (never stored in env vars).
type VaultConfig struct {
	Address    string
	SecretPath string
	MountPath  string
}

// SIEMConfig holds Elasticsearch SIEM exporter configuration.
type SIEMConfig struct {
	ElasticsearchURL string
	IndexPrefix      string
	RetryMax         int
	RetryBackoff     time.Duration
}

// OPAConfig holds Open Policy Agent configuration.
type OPAConfig struct {
	PolicyDir string
	BundleURL string
}

// requiredGatewayVars lists env var names that must be present for the gateway to start.
var requiredGatewayVars = []string{
	"AASG_TLS_CERT_PATH",
	"AASG_TLS_KEY_PATH",
	"AASG_TLS_CA_PATH",
	"AASG_ETCD_ENDPOINTS",
	"AASG_VAULT_ADDR",
	"AASG_VAULT_SECRET_PATH",
	"AASG_SIEM_URL",
	"AASG_OPA_POLICY_DIR",
}

// requiredOrchestratorVars lists env var names that must be present for the orchestrator.
var requiredOrchestratorVars = []string{
	"AASG_TLS_CERT_PATH",
	"AASG_TLS_KEY_PATH",
	"AASG_TLS_CA_PATH",
	"AASG_ETCD_ENDPOINTS",
	"AASG_VAULT_ADDR",
	"AASG_VAULT_SECRET_PATH",
	"AASG_SIEM_URL",
}

// loadRequired validates that all required env vars are present.
// Returns a map of var name to value, or an error listing all missing vars.
func loadRequired(required []string) (map[string]string, error) {
	values := make(map[string]string, len(required))
	var missing []string
	for _, key := range required {
		v, ok := os.LookupEnv(key)
		if !ok {
			missing = append(missing, key)
			continue
		}
		values[key] = v
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return values, nil
}

// buildL2Config constructs an L2Config from environment variables, applying defaults.
func buildL2Config() L2Config {
	return L2Config{
		HeartbeatInterval: lookupDurationOrDefault("AASG_L2_T_HB", DefaultHeartbeatInterval),
		MissTimeout:       lookupDurationOrDefault("AASG_L2_T_MISS", DefaultMissTimeout),
		ProbeTimeout:      lookupDurationOrDefault("AASG_L2_T_PROBE", DefaultProbeTimeout),
		MaxProbes:         lookupIntOrDefault("AASG_L2_N_MAX", DefaultMaxProbes),
		SessionKeyTTL:     lookupDurationOrDefault("AASG_L2_SESSION_KEY_TTL", DefaultSessionKeyTTL),
		SessionKeyDrain:   lookupDurationOrDefault("AASG_L2_SESSION_KEY_DRAIN", DefaultSessionKeyDrainWindow),
		CertTTL:           lookupDurationOrDefault("AASG_L2_CERT_TTL", DefaultCertTTL),
		ReplayWindowSize:  lookupIntOrDefault("AASG_L2_REPLAY_WINDOW", DefaultReplayWindowSize),
	}
}

// buildTLSConfig constructs a TLSConfig from the required vars map.
func buildTLSConfig(required map[string]string) TLSConfig {
	return TLSConfig{
		CertPath:   required["AASG_TLS_CERT_PATH"],
		KeyPath:    required["AASG_TLS_KEY_PATH"],
		CAPath:     required["AASG_TLS_CA_PATH"],
		SPIFFEAddr: lookupEnvOrDefault("AASG_SPIFFE_ADDR", ""),
	}
}

// buildEtcdConfig constructs an EtcdConfig from the required vars map.
func buildEtcdConfig(required map[string]string) EtcdConfig {
	return EtcdConfig{
		Endpoints:   strings.Split(required["AASG_ETCD_ENDPOINTS"], ","),
		DialTimeout: lookupDurationOrDefault("AASG_ETCD_DIAL_TIMEOUT", DefaultEtcdDialTimeout),
	}
}

// buildVaultConfig constructs a VaultConfig from the required vars map.
func buildVaultConfig(required map[string]string) VaultConfig {
	return VaultConfig{
		Address:    required["AASG_VAULT_ADDR"],
		SecretPath: required["AASG_VAULT_SECRET_PATH"],
		MountPath:  lookupEnvOrDefault("AASG_VAULT_MOUNT", "secret"),
	}
}

// buildSIEMConfig constructs a SIEMConfig from the required vars map.
func buildSIEMConfig(required map[string]string) SIEMConfig {
	return SIEMConfig{
		ElasticsearchURL: required["AASG_SIEM_URL"],
		IndexPrefix:      lookupEnvOrDefault("AASG_SIEM_INDEX_PREFIX", "aasg-audit"),
		RetryMax:         lookupIntOrDefault("AASG_SIEM_RETRY_MAX", 3),
		RetryBackoff:     lookupDurationOrDefault("AASG_SIEM_RETRY_BACKOFF", BackoffInitial),
	}
}

// LoadGateway loads GatewayConfig from environment variables.
// Fails fast on any missing required variable, listing all missing vars in the error.
func LoadGateway() (*GatewayConfig, error) {
	required, err := loadRequired(requiredGatewayVars)
	if err != nil {
		return nil, err
	}
	return &GatewayConfig{
		Env: lookupEnvOrDefault("AASG_ENV", "development"),
		HTTP: HTTPConfig{
			Addr:            lookupEnvOrDefault("AASG_HTTP_ADDR", DefaultHTTPAddr),
			GRPCAddr:        lookupEnvOrDefault("AASG_GRPC_ADDR", DefaultGRPCAddr),
			WSAddr:          lookupEnvOrDefault("AASG_WS_ADDR", DefaultWSAddr),
			UpstreamTimeout: lookupDurationOrDefault("AASG_UPSTREAM_TIMEOUT", DefaultUpstreamTimeout),
		},
		L2:    buildL2Config(),
		TLS:   buildTLSConfig(required),
		Etcd:  buildEtcdConfig(required),
		Vault: buildVaultConfig(required),
		SIEM:  buildSIEMConfig(required),
		OPA: OPAConfig{
			PolicyDir: required["AASG_OPA_POLICY_DIR"],
			BundleURL: lookupEnvOrDefault("AASG_OPA_BUNDLE_URL", ""),
		},
	}, nil
}

// LoadOrchestrator loads OrchestratorConfig from environment variables.
// Fails fast on any missing required variable.
func LoadOrchestrator() (*OrchestratorConfig, error) {
	required, err := loadRequired(requiredOrchestratorVars)
	if err != nil {
		return nil, err
	}
	return &OrchestratorConfig{
		Env:      lookupEnvOrDefault("AASG_ENV", "development"),
		GRPCAddr: lookupEnvOrDefault("AASG_GRPC_ADDR", DefaultGRPCAddr),
		L2:       buildL2Config(),
		TLS:      buildTLSConfig(required),
		Etcd:     buildEtcdConfig(required),
		Vault:    buildVaultConfig(required),
		SIEM:     buildSIEMConfig(required),
	}, nil
}

// lookupEnvOrDefault returns the env var value or the default if not set.
func lookupEnvOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// lookupDurationOrDefault returns the env var parsed as a duration, or the default.
func lookupDurationOrDefault(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

// lookupIntOrDefault returns the env var parsed as an int, or the default.
func lookupIntOrDefault(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
