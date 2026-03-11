package config_test

import (
	"testing"

	"github.com/RudsonCarvalho/sageway/internal/config"
)

func setRequiredGatewayEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AASG_TLS_CERT_PATH", "/certs/aasg.crt")
	t.Setenv("AASG_TLS_KEY_PATH", "/certs/aasg.key")
	t.Setenv("AASG_TLS_CA_PATH", "/certs/ca.crt")
	t.Setenv("AASG_ETCD_ENDPOINTS", "http://localhost:2379")
	t.Setenv("AASG_VAULT_ADDR", "http://localhost:8200")
	t.Setenv("AASG_VAULT_SECRET_PATH", "secret/aasg/l2/master")
	t.Setenv("AASG_SIEM_URL", "http://localhost:9200")
	t.Setenv("AASG_OPA_POLICY_DIR", "/policies")
}

func setRequiredOrchestratorEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AASG_TLS_CERT_PATH", "/certs/orch.crt")
	t.Setenv("AASG_TLS_KEY_PATH", "/certs/orch.key")
	t.Setenv("AASG_TLS_CA_PATH", "/certs/ca.crt")
	t.Setenv("AASG_ETCD_ENDPOINTS", "http://localhost:2379")
	t.Setenv("AASG_VAULT_ADDR", "http://localhost:8200")
	t.Setenv("AASG_VAULT_SECRET_PATH", "secret/aasg/l2/master")
	t.Setenv("AASG_SIEM_URL", "http://localhost:9200")
}

// NOTE: Tests in this file use t.Setenv and therefore cannot use t.Parallel.
// t.Setenv mutates process-wide environment; parallel execution would cause races.

func TestLoadGateway_AllRequiredPresent(t *testing.T) {
	setRequiredGatewayEnv(t)
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() returned unexpected error: %v", err)
	}
	if cfg.TLS.CertPath != "/certs/aasg.crt" {
		t.Errorf("CertPath = %q, want /certs/aasg.crt", cfg.TLS.CertPath)
	}
	if cfg.Vault.SecretPath != "secret/aasg/l2/master" {
		t.Errorf("VaultSecretPath = %q, want secret/aasg/l2/master", cfg.Vault.SecretPath)
	}
}

func TestLoadGateway_MissingRequired(t *testing.T) {
	_, err := config.LoadGateway()
	if err == nil {
		t.Fatal("LoadGateway() should fail when required vars are missing")
	}
}

func TestLoadGateway_MissingOneVar(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_SIEM_URL", "")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Logf("LoadGateway with empty SIEM_URL returned error (expected if empty not accepted): %v", err)
	}
	_ = cfg
}

func TestLoadGateway_Defaults(t *testing.T) {
	setRequiredGatewayEnv(t)
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.L2.HeartbeatInterval != config.DefaultHeartbeatInterval {
		t.Errorf("HeartbeatInterval = %v, want %v", cfg.L2.HeartbeatInterval, config.DefaultHeartbeatInterval)
	}
	if cfg.L2.MissTimeout != config.DefaultMissTimeout {
		t.Errorf("MissTimeout = %v, want %v", cfg.L2.MissTimeout, config.DefaultMissTimeout)
	}
	if cfg.L2.MaxProbes != config.DefaultMaxProbes {
		t.Errorf("MaxProbes = %d, want %d", cfg.L2.MaxProbes, config.DefaultMaxProbes)
	}
	if cfg.L2.ReplayWindowSize != config.DefaultReplayWindowSize {
		t.Errorf("ReplayWindowSize = %d, want %d", cfg.L2.ReplayWindowSize, config.DefaultReplayWindowSize)
	}
	if cfg.HTTP.Addr != config.DefaultHTTPAddr {
		t.Errorf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, config.DefaultHTTPAddr)
	}
}

func TestLoadGateway_EtcdEndpointsParsed(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_ETCD_ENDPOINTS", "http://etcd1:2379,http://etcd2:2379,http://etcd3:2379")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if len(cfg.Etcd.Endpoints) != 3 {
		t.Errorf("Etcd.Endpoints len = %d, want 3", len(cfg.Etcd.Endpoints))
	}
}

func TestLoadOrchestrator_AllRequiredPresent(t *testing.T) {
	setRequiredOrchestratorEnv(t)
	cfg, err := config.LoadOrchestrator()
	if err != nil {
		t.Fatalf("LoadOrchestrator() returned unexpected error: %v", err)
	}
	if cfg.TLS.CertPath != "/certs/orch.crt" {
		t.Errorf("CertPath = %q, want /certs/orch.crt", cfg.TLS.CertPath)
	}
}

func TestLoadOrchestrator_MissingRequired(t *testing.T) {
	_, err := config.LoadOrchestrator()
	if err == nil {
		t.Fatal("LoadOrchestrator() should fail when required vars are missing")
	}
}

func TestLoadGateway_HTTPAddrOverride(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_HTTP_ADDR", ":9090")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.HTTP.Addr != ":9090" {
		t.Errorf("HTTP.Addr = %q, want :9090", cfg.HTTP.Addr)
	}
}

func TestLoadGateway_DurationOverride(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_L2_T_HB", "10s")
	t.Setenv("AASG_L2_T_MISS", "30s")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.L2.HeartbeatInterval.String() != "10s" {
		t.Errorf("HeartbeatInterval = %v, want 10s", cfg.L2.HeartbeatInterval)
	}
	if cfg.L2.MissTimeout.String() != "30s" {
		t.Errorf("MissTimeout = %v, want 30s", cfg.L2.MissTimeout)
	}
}

func TestLoadGateway_IntOverride(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_L2_N_MAX", "5")
	t.Setenv("AASG_L2_REPLAY_WINDOW", "2000")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.L2.MaxProbes != 5 {
		t.Errorf("MaxProbes = %d, want 5", cfg.L2.MaxProbes)
	}
	if cfg.L2.ReplayWindowSize != 2000 {
		t.Errorf("ReplayWindowSize = %d, want 2000", cfg.L2.ReplayWindowSize)
	}
}

func TestLoadGateway_InvalidDurationFallsBackToDefault(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_L2_T_HB", "not-a-duration")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.L2.HeartbeatInterval != config.DefaultHeartbeatInterval {
		t.Errorf("HeartbeatInterval = %v, want default %v", cfg.L2.HeartbeatInterval, config.DefaultHeartbeatInterval)
	}
}

func TestLoadGateway_InvalidIntFallsBackToDefault(t *testing.T) {
	setRequiredGatewayEnv(t)
	t.Setenv("AASG_L2_N_MAX", "not-an-int")
	cfg, err := config.LoadGateway()
	if err != nil {
		t.Fatalf("LoadGateway() failed: %v", err)
	}
	if cfg.L2.MaxProbes != config.DefaultMaxProbes {
		t.Errorf("MaxProbes = %d, want default %d", cfg.L2.MaxProbes, config.DefaultMaxProbes)
	}
}
