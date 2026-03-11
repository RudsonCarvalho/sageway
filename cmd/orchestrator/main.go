// Command orchestrator is the entrypoint for the AASG Orchestrator.
// It supervises all AASG agents via the L2 Control Plane, manages policy distribution,
// and handles quarantine authorization (dual-sign Ed25519).
// Missing required environment variables cause an immediate exit.
package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/RudsonCarvalho/sageway/internal/config"
)

func main() {
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Str("service", "aasg-orchestrator").
		Logger()

	cfg, err := config.LoadOrchestrator()
	if err != nil {
		logger.Error().Err(err).Msg("configuration validation failed: refusing to start")
		os.Exit(1)
	}

	logger.Info().
		Str("env", cfg.Env).
		Str("grpc_addr", cfg.GRPCAddr).
		Msg("orchestrator configuration loaded successfully")

	// Phase 4: L2 channel listener -- see feat/phase-04-l2-channel
	// Phase 8: supervisor, router, quarantine -- see feat/phase-08-orchestrator

	logger.Info().Msg("orchestrator ready (stub -- full implementation in Phase 8)")
	select {}
}
