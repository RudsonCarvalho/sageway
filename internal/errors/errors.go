// Package errors defines all sentinel errors and the API error response type for the AASG system.
//
// Security contract: internal error details (cryptographic failures, replay detection)
// are NEVER returned to API callers. ToGatewayError maps internal errors to safe external codes.
// Full details are written to the Audit Ledger only.
package errors

import "errors"

// Sentinel errors -- L2 Control Plane.
var (
	// ErrReplayDetected is returned when a seq_num has been seen before within the replay window.
	// This is a security event that must be logged to the Audit Ledger.
	ErrReplayDetected = errors.New("replay detected: seq_num already seen")

	// ErrHMACVerifyFailed is returned when HMAC verification fails.
	// NEVER returned to API callers -- mapped to ErrAuthFailed externally.
	ErrHMACVerifyFailed = errors.New("hmac verification failed")

	// ErrSeqNumNotMonotonic is returned when a seq_num is not greater than the last accepted.
	ErrSeqNumNotMonotonic = errors.New("seq_num is not monotonically increasing")

	// ErrClockDrift is returned when a message timestamp deviates beyond the allowed tolerance.
	ErrClockDrift = errors.New("message timestamp outside acceptable clock drift window")

	// ErrNonceMismatch is returned when a CHALLENGE_RESP does not match the issued nonce.
	ErrNonceMismatch = errors.New("challenge nonce mismatch")

	// ErrChannelClosed is returned when an operation is attempted on a closed L2 channel.
	ErrChannelClosed = errors.New("L2 channel is closed")

	// ErrHandshakeTimeout is returned when the L2 mTLS handshake does not complete in time.
	ErrHandshakeTimeout = errors.New("L2 handshake timeout")

	// ErrPayloadTampered is returned when payload digest does not match the envelope HMAC.
	ErrPayloadTampered = errors.New("payload tampered: digest mismatch")

	// ErrSessionKeyNotFound is returned when no active session key exists for the agent.
	ErrSessionKeyNotFound = errors.New("session key not found for agent")

	// ErrPolicyApplicationFailed is returned when a POLICY_UPDATE cannot be applied.
	// The Orchestrator must be notified via POLICY_ACK with applied=false.
	ErrPolicyApplicationFailed = errors.New("policy application failed")
)

// Sentinel errors -- Authentication and certificates.
var (
	// ErrAuthFailed is the external-facing authentication failure sentinel.
	// Maps to HTTP 401. Never reveals whether the failure was cert, HMAC, or token.
	ErrAuthFailed = errors.New("authentication failed")

	// ErrCertificateExpired is returned when a peer certificate has expired.
	// Expired certificates are NEVER accepted, even temporarily (Spec s7.3).
	ErrCertificateExpired = errors.New("certificate is expired")

	// ErrSPIFFEIDInvalid is returned when a certificate SPIFFE ID does not match
	// the expected pattern for the peer role.
	ErrSPIFFEIDInvalid = errors.New("invalid SPIFFE ID")

	// ErrInvalidSignature is returned when an Ed25519 quarantine approval signature is invalid.
	ErrInvalidSignature = errors.New("invalid Ed25519 signature")

	// ErrDuplicateApprover is returned when the same operator attempts to sign twice.
	// Each approver for dual-sign quarantine must be a distinct operator (Spec s6.2).
	ErrDuplicateApprover = errors.New("duplicate approver: same operator cannot sign twice")

	// ErrInsufficientApprovals is returned when a quarantine command is submitted
	// without the required number of distinct operator signatures.
	ErrInsufficientApprovals = errors.New("insufficient quarantine approvals: requires 2 distinct operators")
)

// Sentinel errors -- Policy and routing.
var (
	// ErrPolicyDenied is returned when OPA evaluates a request as denied.
	ErrPolicyDenied = errors.New("policy denied")

	// ErrRateLimited is returned when a client exceeds its token bucket allocation.
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrAgentUnavailable is returned when the target agent is QUARANTINED or UNRESPONSIVE.
	ErrAgentUnavailable = errors.New("agent unavailable")

	// ErrAgentDegraded is returned when the target agent is in DEGRADED operational state.
	ErrAgentDegraded = errors.New("agent degraded")

	// ErrPolicyDrift is returned when an agent active policy digest diverges from expected.
	ErrPolicyDrift = errors.New("policy drift: agent policy diverges from orchestrator")

	// ErrUpstreamTimeout is returned when an upstream service does not respond within the SLA.
	ErrUpstreamTimeout = errors.New("upstream timeout")

	// ErrSCCViolation is returned when a data transfer is blocked by SCC policy (Spec s6.1.3).
	// Mapped to HTTP 451 (Unavailable For Legal Reasons).
	ErrSCCViolation = errors.New("SCC violation: transfer blocked by policy")

	// ErrQuarantineActive is returned when a request targets a quarantined agent.
	ErrQuarantineActive = errors.New("agent is quarantined")

	// ErrCircuitBreakerOpen is returned when the upstream circuit breaker is in OPEN state.
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

	// ErrUnknownAgent is returned when an agent ID is not registered with the Orchestrator.
	ErrUnknownAgent = errors.New("unknown agent ID")
)

// Sentinel errors -- Infrastructure.
var (
	// ErrEtcdUnavailable is returned when etcd is unreachable.
	// The system degrades gracefully: seq_num falls back to in-memory with a CRITICAL log.
	ErrEtcdUnavailable = errors.New("etcd unavailable")

	// ErrAuditLedgerUnavailable is returned when the Merkle Audit Ledger is unreachable.
	// This triggers fail-closed behavior: the AASG stops accepting new requests (Spec s8.2).
	ErrAuditLedgerUnavailable = errors.New("audit ledger unavailable: fail-closed")

	// ErrVaultUnavailable is returned when the Vault secret backend is unreachable.
	ErrVaultUnavailable = errors.New("vault unavailable")

	// ErrMasterSecretNotFound is returned when the Vault path for the master secret is empty.
	ErrMasterSecretNotFound = errors.New("master secret not found at configured vault path")
)

// ErrorCode is the machine-readable error code returned in API error responses (Spec s6.1.3).
type ErrorCode string

// API error codes as defined in Spec s6.1.3.
const (
	CodeAuthFailed       ErrorCode = "AUTH_FAILED"
	CodePolicyDenied     ErrorCode = "POLICY_DENIED"
	CodeRateLimited      ErrorCode = "RATE_LIMITED"
	CodeAgentUnavailable ErrorCode = "AGENT_UNAVAILABLE"
	CodeAgentDegraded    ErrorCode = "AGENT_DEGRADED"
	CodePolicyDrift      ErrorCode = "POLICY_DRIFT"
	CodeUpstreamTimeout  ErrorCode = "UPSTREAM_TIMEOUT"
	CodeSCCViolation     ErrorCode = "SCC_VIOLATION"
	CodeInternalError    ErrorCode = "INTERNAL_ERROR"
)

// HTTPStatus returns the HTTP status code for this error code (Spec s6.1.3).
func (c ErrorCode) HTTPStatus() int {
	switch c {
	case CodeAuthFailed:
		return 401
	case CodePolicyDenied:
		return 403
	case CodeRateLimited:
		return 429
	case CodeSCCViolation:
		return 451
	case CodeAgentDegraded:
		return 502
	case CodeAgentUnavailable, CodePolicyDrift:
		return 503
	case CodeUpstreamTimeout:
		return 504
	case CodeInternalError:
		return 500
	}
	return 500
}

// GatewayError is the structured error body returned by all API handlers on failure (Spec s6.1.2).
type GatewayError struct {
	Code       ErrorCode
	Message    string
	RequestID  string
	RetryAfter int
}

// Error implements the error interface.
func (e *GatewayError) Error() string {
	return string(e.Code) + ": " + e.Message
}

// ToGatewayError converts an internal sentinel error to a safe external GatewayError.
// This is the single mapping point between internal errors and external codes.
// Internal details are NEVER included in the result -- they belong in the Audit Ledger.
func ToGatewayError(err error, requestID string) *GatewayError {
	switch {
	case errors.Is(err, ErrAuthFailed),
		errors.Is(err, ErrHMACVerifyFailed),
		errors.Is(err, ErrCertificateExpired),
		errors.Is(err, ErrSPIFFEIDInvalid),
		errors.Is(err, ErrInvalidSignature):
		return &GatewayError{Code: CodeAuthFailed, Message: "authentication failed", RequestID: requestID}
	case errors.Is(err, ErrPolicyDenied):
		return &GatewayError{Code: CodePolicyDenied, Message: "agent not authorized", RequestID: requestID}
	case errors.Is(err, ErrRateLimited):
		return &GatewayError{Code: CodeRateLimited, Message: "too many requests", RequestID: requestID, RetryAfter: 30}
	case errors.Is(err, ErrQuarantineActive),
		errors.Is(err, ErrAgentUnavailable),
		errors.Is(err, ErrCircuitBreakerOpen):
		return &GatewayError{Code: CodeAgentUnavailable, Message: "agent unavailable", RequestID: requestID}
	case errors.Is(err, ErrAgentDegraded):
		return &GatewayError{Code: CodeAgentDegraded, Message: "agent degraded", RequestID: requestID}
	case errors.Is(err, ErrPolicyDrift):
		return &GatewayError{Code: CodePolicyDrift, Message: "policy drift: awaiting resync", RequestID: requestID}
	case errors.Is(err, ErrUpstreamTimeout):
		return &GatewayError{Code: CodeUpstreamTimeout, Message: "upstream timeout", RequestID: requestID}
	case errors.Is(err, ErrSCCViolation):
		return &GatewayError{Code: CodeSCCViolation, Message: "transfer blocked by SCC policy", RequestID: requestID}
	default:
		return &GatewayError{Code: CodeInternalError, Message: "internal server error", RequestID: requestID}
	}
}
