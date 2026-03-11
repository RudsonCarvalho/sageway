package config_test

import (
	"testing"
	"time"

	"github.com/RudsonCarvalho/sageway/internal/config"
)

// TestL2TimingInvariants verifies that the L2 timing parameters satisfy
// the spec invariants: T_miss >= 2*T_hb and T_probe < T_miss/N_max.
func TestL2TimingInvariants(t *testing.T) {
	t.Parallel()
	if config.DefaultMissTimeout < 2*config.DefaultHeartbeatInterval {
		t.Errorf("T_miss (%v) must be >= 2*T_hb (%v)",
			config.DefaultMissTimeout, 2*config.DefaultHeartbeatInterval)
	}
	maxProbeWindow := config.DefaultMissTimeout / time.Duration(config.DefaultMaxProbes)
	if config.DefaultProbeTimeout >= maxProbeWindow {
		t.Errorf("T_probe (%v) must be < T_miss/N_max (%v)",
			config.DefaultProbeTimeout, maxProbeWindow)
	}
}

// TestHKDFInfoNotEmpty verifies that the HKDF info string is set.
func TestHKDFInfoNotEmpty(t *testing.T) {
	t.Parallel()
	if config.HKDFInfo == "" {
		t.Error("HKDFInfo must not be empty")
	}
}

// TestHKDFSaltSizeIsSufficient verifies that the salt size meets minimum cryptographic requirements.
func TestHKDFSaltSizeIsSufficient(t *testing.T) {
	t.Parallel()
	const minSaltBytes = 16
	if config.HKDFSaltSize < minSaltBytes {
		t.Errorf("HKDFSaltSize (%d) must be >= %d bytes", config.HKDFSaltSize, minSaltBytes)
	}
}

// TestChallengeNonceSizeIsSufficient verifies that challenge nonces are cryptographically strong.
func TestChallengeNonceSizeIsSufficient(t *testing.T) {
	t.Parallel()
	const minNonceBytes = 16
	if config.ChallengeNonceSize < minNonceBytes {
		t.Errorf("ChallengeNonceSize (%d) must be >= %d bytes", config.ChallengeNonceSize, minNonceBytes)
	}
}

// TestQuarantineRequiresTwoApprovers verifies the dual-sign requirement.
func TestQuarantineRequiresTwoApprovers(t *testing.T) {
	t.Parallel()
	if config.QuarantineApprovalsRequired < 2 {
		t.Errorf("QuarantineApprovalsRequired must be >= 2, got %d", config.QuarantineApprovalsRequired)
	}
}

// TestReplayWindowSize verifies the sliding window is large enough for peak load.
func TestReplayWindowSize(t *testing.T) {
	t.Parallel()
	const minWindow = 100
	if config.DefaultReplayWindowSize < minWindow {
		t.Errorf("ReplayWindowSize (%d) must be >= %d", config.DefaultReplayWindowSize, minWindow)
	}
}

// TestCertTTLNotExceed24h verifies SPIFFE/SPIRE max cert validity (Spec s7.3).
func TestCertTTLNotExceed24h(t *testing.T) {
	t.Parallel()
	maxAllowed := 24 * time.Hour
	if config.DefaultCertTTL > maxAllowed {
		t.Errorf("DefaultCertTTL (%v) exceeds SPIFFE/SPIRE max of %v", config.DefaultCertTTL, maxAllowed)
	}
}
