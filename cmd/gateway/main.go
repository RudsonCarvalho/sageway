// Command gateway is the entrypoint for the AASG gateway.
// It loads and validates all configuration from environment variables,
// then starts the HTTP, gRPC, and L2 Control Plane servers.
// Missing required environment variables cause an immediate exit -- no defaults for secrets.
package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/RudsonCarvalho/sageway/internal/config"
)

func main() {
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Str("service", "aasg-gateway").
		Logger()

	cfg, err := config.LoadGateway()
	if err != nil {
		logger.Error().Err(err).Msg("configuration validation failed: refusing to start")
		os.Exit(1)
	}

	logger.Info().
		Str("env", cfg.Env).
		Str("http_addr", cfg.HTTP.Addr).
		Str("grpc_addr", cfg.HTTP.GRPCAddr).
		Msg("gateway configuration loaded successfully")

	// Phase 4: L2 channel -- see feat/phase-04-l2-channel
	// Phase 5: OPA policy engine -- see feat/phase-05-policy
	// Phase 6: audit ledger -- see feat/phase-06-audit
	// Phase 7: HTTP/gRPC servers -- see feat/phase-07-api

	logger.Info().Msg("gateway ready (stub -- servers implemented in Phase 7)")
	select {}
}
