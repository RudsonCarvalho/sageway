// Package config holds all configuration constants and structs for the AASG system.
// All timing parameters, window sizes, and operational thresholds are defined here.
// No magic numbers anywhere else in the codebase — all constants live in this file.
package config

import "time"

// L2 Control Plane timing parameters (Spec §5.2).
const (
	// DefaultHeartbeatInterval is the frequency at which the AASG publishes its state.
	// Decreasing this value increases network load; increasing it slows failure detection.
	DefaultHeartbeatInterval = 5 * time.Second

	// DefaultMissTimeout is the time the Orchestrator waits without a heartbeat before
	// transitioning the AASG to STALE. Must be >= 2 * DefaultHeartbeatInterval.
	DefaultMissTimeout = 10 * time.Second

	// DefaultProbeTimeout is the maximum time the AASG has to respond to a CHALLENGE.
	// Must be < DefaultMissTimeout / DefaultMaxProbes.
	DefaultProbeTimeout = 2 * time.Second

	// DefaultMaxProbes is the number of consecutive unanswered CHALLENGEs before
	// the Orchestrator transitions the AASG to UNRESPONSIVE.
	DefaultMaxProbes = 3

	// DefaultSessionKeyTTL is the HMAC session key rotation interval.
	// A 30-second drain window (DefaultSessionKeyDrainWindow) is applied after rotation
	// to allow in-flight messages signed with the previous key to be validated.
	DefaultSessionKeyTTL = 3600 * time.Second

	// DefaultSessionKeyDrainWindow is the overlap period after a key rotation during
	// which both old and new session keys are accepted for verification.
	DefaultSessionKeyDrainWindow = 30 * time.Second

	// DefaultCertTTL is the maximum certificate validity period enforced by SPIFFE/SPIRE.
	// Certificates must be renewed before expiry; the system rejects expired certs immediately.
	DefaultCertTTL = 86400 * time.Second

	// DefaultReplayWindowSize is the size of the sliding window of accepted seq_nums.
	// Any seq_num already present in the window is rejected as a replay attack.
	DefaultReplayWindowSize = 1000
)

// HKDF session key derivation parameters (Spec §5, user confirmation).
const (
	// HKDFInfo is the static info string used in HKDF derivation. Changing this value
	// invalidates all existing session keys — requires coordinated deployment.
	HKDFInfo = "aasg-l2-session-v1"

	// HKDFSaltSize is the size in bytes of the random nonce exchanged during the
	// initial mTLS handshake, used as the HKDF salt.
	HKDFSaltSize = 32

	// HMACKeySize is the output key size in bytes for HMAC-SHA256.
	HMACKeySize = 32
)

// ChallengeNonceSize is the size in bytes of the random nonce included in CHALLENGE messages.
// Must be cryptographically random (crypto/rand).
const ChallengeNonceSize = 32

// Dual-sign quarantine parameters (Spec §6.2, user confirmation).
const (
	// QuarantineApprovalsRequired is the number of distinct operator Ed25519 signatures
	// required before the Orchestrator issues a QUARANTINE_CMD.
	QuarantineApprovalsRequired = 2
)

// Performance SLA targets (Spec §7.1).
// These constants are used in benchmark tests to assert SLA compliance.
const (
	// TargetAPIThroughputRPS is the minimum sustained API throughput in requests per second.
	TargetAPIThroughputRPS = 5000

	// TargetL2ThroughputMPS is the minimum sustained L2 heartbeat throughput in messages per second.
	TargetL2ThroughputMPS = 500

	// TargetAddedLatencyP99Ms is the maximum latency added by the AASG at p99, in milliseconds.
	TargetAddedLatencyP99Ms = 10

	// TargetHeartbeatLatencyP99Ms is the maximum heartbeat latency (AASG→Orchestrator) at p99, ms.
	TargetHeartbeatLatencyP99Ms = 50

	// TargetFailureDetection is the maximum time to detect a STALE state (T_miss + network jitter).
	TargetFailureDetection = 12 * time.Second

	// TargetQuarantineTime is the maximum time from N_max failures to QUARANTINED state.
	// Computed as T_miss + N_max * T_probe = 10 + 3*2 = 16s + processing overhead.
	TargetQuarantineTime = 20 * time.Second

	// TargetHMACLatencyP99Ms is the maximum HMAC sign+verify latency at p99, in milliseconds.
	TargetHMACLatencyP99Ms = 1
)

// Circuit breaker thresholds (Spec §8.3).
const (
	// CBErrorRateThreshold is the error rate (0.0–1.0) that triggers the circuit breaker OPEN state.
	CBErrorRateThreshold = 0.50

	// CBWindowSeconds is the sliding window duration for error rate calculation.
	CBWindowSeconds = 10

	// CBConsecutiveFailures is the number of consecutive failures that open the circuit breaker.
	CBConsecutiveFailures = 5

	// CBHalfOpenTimeout is the duration in OPEN state before transitioning to HALF-OPEN.
	CBHalfOpenTimeout = 30 * time.Second

	// CBHalfOpenRequests is the number of probe requests allowed in HALF-OPEN state.
	CBHalfOpenRequests = 1
)

// Connection and reconnect parameters (Spec §8.1).
const (
	// BackoffInitial is the initial reconnect backoff for the L2 channel.
	BackoffInitial = 1 * time.Second

	// BackoffMax is the maximum reconnect backoff for the L2 channel.
	BackoffMax = 30 * time.Second

	// BackoffFactor is the multiplicative factor applied to backoff on each retry.
	BackoffFactor = 2.0
)

// Network addresses (defaults).
const (
	DefaultHTTPAddr = ":8080"
	DefaultGRPCAddr = ":9443"
	DefaultWSAddr   = ":8443"
)

// API defaults (Spec §6.1).
const (
	// DefaultUpstreamTimeout is the maximum time to wait for an upstream service response.
	DefaultUpstreamTimeout = 30 * time.Second

	// DefaultRateLimitPerClient is the default token bucket rate per client in requests per second.
	DefaultRateLimitPerClient = 100

	// DefaultRateLimitGlobal is the global rate limit in requests per second across all clients.
	DefaultRateLimitGlobal = 5000
)

// Storage and retention parameters.
const (
	// AuditRetentionDays is the minimum audit log retention period required by SCC compliance.
	AuditRetentionDays = 365

	// EtcdSeqNumKeyPrefix is the etcd key prefix for persisting L2 sequence numbers per agent.
	EtcdSeqNumKeyPrefix = "/aasg/l2/seqnum/"

	// DefaultEtcdDialTimeout is the maximum time to establish an etcd connection.
	DefaultEtcdDialTimeout = 5 * time.Second

	// AgentTaskQueueTimeout is the maximum time an SCC agent task can remain queued.
	AgentTaskQueueTimeout = 30 * time.Minute
)

// ClockDriftTolerance is the maximum allowed difference between the message timestamp and
// the current time. Messages with timestamps outside this window are rejected as potential
// replay attacks or clock misconfiguration.
const ClockDriftTolerance = 30 * time.Second
