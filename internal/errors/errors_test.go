package errors_test

import (
	"errors"
	"fmt"
	"testing"

	aerrors "github.com/RudsonCarvalho/sageway/internal/errors"
)

func TestSentinelErrorsAreDistinct(t *testing.T) {
	t.Parallel()
	sentinels := []error{
		aerrors.ErrReplayDetected,
		aerrors.ErrHMACVerifyFailed,
		aerrors.ErrSeqNumNotMonotonic,
		aerrors.ErrClockDrift,
		aerrors.ErrNonceMismatch,
		aerrors.ErrChannelClosed,
		aerrors.ErrHandshakeTimeout,
		aerrors.ErrPayloadTampered,
		aerrors.ErrSessionKeyNotFound,
		aerrors.ErrPolicyApplicationFailed,
		aerrors.ErrAuthFailed,
		aerrors.ErrCertificateExpired,
		aerrors.ErrSPIFFEIDInvalid,
		aerrors.ErrInvalidSignature,
		aerrors.ErrDuplicateApprover,
		aerrors.ErrInsufficientApprovals,
		aerrors.ErrPolicyDenied,
		aerrors.ErrRateLimited,
		aerrors.ErrAgentUnavailable,
		aerrors.ErrAgentDegraded,
		aerrors.ErrPolicyDrift,
		aerrors.ErrUpstreamTimeout,
		aerrors.ErrSCCViolation,
		aerrors.ErrQuarantineActive,
		aerrors.ErrCircuitBreakerOpen,
		aerrors.ErrUnknownAgent,
		aerrors.ErrEtcdUnavailable,
		aerrors.ErrAuditLedgerUnavailable,
		aerrors.ErrVaultUnavailable,
		aerrors.ErrMasterSecretNotFound,
	}
	for i, a := range sentinels {
		for j, b := range sentinels {
			if i != j && errors.Is(a, b) {
				t.Errorf("sentinel %v incorrectly matches sentinel %v", a, b)
			}
		}
	}
}

func TestToGatewayError_AuthErrors(t *testing.T) {
	t.Parallel()
	authErrors := []error{
		aerrors.ErrAuthFailed,
		aerrors.ErrHMACVerifyFailed,
		aerrors.ErrCertificateExpired,
		aerrors.ErrSPIFFEIDInvalid,
		aerrors.ErrInvalidSignature,
	}
	for _, err := range authErrors {
		ge := aerrors.ToGatewayError(err, "req-123")
		if ge.Code != aerrors.CodeAuthFailed {
			t.Errorf("expected %s for %v, got %s", aerrors.CodeAuthFailed, err, ge.Code)
		}
		if ge.RequestID != "req-123" {
			t.Errorf("expected RequestID req-123, got %s", ge.RequestID)
		}
		if ge.Code.HTTPStatus() != 401 {
			t.Errorf("expected HTTP 401 for %v, got %d", ge.Code, ge.Code.HTTPStatus())
		}
	}
}

func TestToGatewayError_SCCViolation(t *testing.T) {
	t.Parallel()
	ge := aerrors.ToGatewayError(aerrors.ErrSCCViolation, "req-scc")
	if ge.Code != aerrors.CodeSCCViolation {
		t.Errorf("expected SCC_VIOLATION, got %s", ge.Code)
	}
	if ge.Code.HTTPStatus() != 451 {
		t.Errorf("expected HTTP 451, got %d", ge.Code.HTTPStatus())
	}
}

func TestToGatewayError_RateLimited(t *testing.T) {
	t.Parallel()
	ge := aerrors.ToGatewayError(aerrors.ErrRateLimited, "req-rl")
	if ge.Code != aerrors.CodeRateLimited {
		t.Errorf("expected RATE_LIMITED, got %s", ge.Code)
	}
	if ge.RetryAfter != 30 {
		t.Errorf("expected RetryAfter=30, got %d", ge.RetryAfter)
	}
	if ge.Code.HTTPStatus() != 429 {
		t.Errorf("expected HTTP 429, got %d", ge.Code.HTTPStatus())
	}
}

func TestToGatewayError_InternalErrorHidesDetails(t *testing.T) {
	t.Parallel()
	internalErr := errors.New("secret internal detail")
	ge := aerrors.ToGatewayError(internalErr, "req-int")
	if ge.Code != aerrors.CodeInternalError {
		t.Errorf("expected INTERNAL_ERROR, got %s", ge.Code)
	}
	if ge.Message != "internal server error" {
		t.Errorf("internal error detail leaked to caller: %s", ge.Message)
	}
}

func TestToGatewayError_WrappedErrorsResolve(t *testing.T) {
	t.Parallel()
	wrapped := fmt.Errorf("outer context: %w", aerrors.ErrPolicyDenied)
	ge := aerrors.ToGatewayError(wrapped, "req-wrap")
	if ge.Code != aerrors.CodePolicyDenied {
		t.Errorf("expected POLICY_DENIED for wrapped error, got %s", ge.Code)
	}
}

func TestHTTPStatus_AllCodesHaveMapping(t *testing.T) {
	t.Parallel()
	codes := []struct {
		code     aerrors.ErrorCode
		wantHTTP int
	}{
		{aerrors.CodeAuthFailed, 401},
		{aerrors.CodePolicyDenied, 403},
		{aerrors.CodeRateLimited, 429},
		{aerrors.CodeSCCViolation, 451},
		{aerrors.CodeAgentDegraded, 502},
		{aerrors.CodeAgentUnavailable, 503},
		{aerrors.CodePolicyDrift, 503},
		{aerrors.CodeUpstreamTimeout, 504},
		{aerrors.CodeInternalError, 500},
	}
	for _, tc := range codes {
		if got := tc.code.HTTPStatus(); got != tc.wantHTTP {
			t.Errorf("HTTPStatus(%s) = %d, want %d", tc.code, got, tc.wantHTTP)
		}
	}
}

func TestGatewayError_ErrorString(t *testing.T) {
	t.Parallel()
	ge := &aerrors.GatewayError{Code: aerrors.CodeAuthFailed, Message: "authentication failed"}
	want := "AUTH_FAILED: authentication failed"
	if got := ge.Error(); got != want {
		t.Errorf("GatewayError.Error() = %q, want %q", got, want)
	}
}

func TestToGatewayError_AgentUnavailableGroup(t *testing.T) {
	t.Parallel()
	unavailableErrors := []error{
		aerrors.ErrAgentUnavailable,
		aerrors.ErrQuarantineActive,
		aerrors.ErrCircuitBreakerOpen,
	}
	for _, err := range unavailableErrors {
		ge := aerrors.ToGatewayError(err, "req-unavail")
		if ge.Code != aerrors.CodeAgentUnavailable {
			t.Errorf("expected %s for %v, got %s", aerrors.CodeAgentUnavailable, err, ge.Code)
		}
		if ge.Code.HTTPStatus() != 503 {
			t.Errorf("expected HTTP 503 for %v, got %d", ge.Code, ge.Code.HTTPStatus())
		}
	}
}

func TestToGatewayError_AgentDegraded(t *testing.T) {
	t.Parallel()
	ge := aerrors.ToGatewayError(aerrors.ErrAgentDegraded, "req-deg")
	if ge.Code != aerrors.CodeAgentDegraded {
		t.Errorf("expected AGENT_DEGRADED, got %s", ge.Code)
	}
	if ge.Code.HTTPStatus() != 502 {
		t.Errorf("expected HTTP 502, got %d", ge.Code.HTTPStatus())
	}
}

func TestToGatewayError_PolicyDrift(t *testing.T) {
	t.Parallel()
	ge := aerrors.ToGatewayError(aerrors.ErrPolicyDrift, "req-drift")
	if ge.Code != aerrors.CodePolicyDrift {
		t.Errorf("expected POLICY_DRIFT, got %s", ge.Code)
	}
	if ge.Code.HTTPStatus() != 503 {
		t.Errorf("expected HTTP 503, got %d", ge.Code.HTTPStatus())
	}
}

func TestToGatewayError_UpstreamTimeout(t *testing.T) {
	t.Parallel()
	ge := aerrors.ToGatewayError(aerrors.ErrUpstreamTimeout, "req-timeout")
	if ge.Code != aerrors.CodeUpstreamTimeout {
		t.Errorf("expected UPSTREAM_TIMEOUT, got %s", ge.Code)
	}
	if ge.Code.HTTPStatus() != 504 {
		t.Errorf("expected HTTP 504, got %d", ge.Code.HTTPStatus())
	}
}
