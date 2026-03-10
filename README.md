# sageway

> Governance-aware security gateway for autonomous AI agents — behavioral supervision, cryptographic integrity, and SCC/GDPR enforcement at the decision boundary.

---

## The problem

Market gateways — Kong, Istio, AWS API Gateway — answer one question: *is this request authorized?*

They do not answer: *is this agent still the same agent that was authorized? Is this action proportionate to the mandate that originated it? Does this data transfer have a valid Transfer Impact Assessment for this destination country?*

An autonomous AI agent can have zero error rate and normal latency while operating outside its mandate. Traffic metrics don't catch that. Sageway does.

---

## What sageway does differently

| Capability | Market gateway | Sageway |
|---|---|---|
| Auth, rate limiting, proxying | ✅ | ✅ (via Envoy data plane) |
| Behavioral supervision | ❌ | ✅ L2 Control Plane |
| Cryptographic agent identity | ❌ | ✅ HMAC challenge/response |
| Policy drift detection | ❌ | ✅ digest-based |
| Anomaly-based quarantine | ❌ | ✅ dual-sign authorization |
| SCC/GDPR enforcement with TIA context | ❌ | ✅ HTTP 451 with regulatory reason |
| Legally defensible audit trail | ❌ | ✅ Merkle Audit Ledger |

---

## Architecture

Sageway operates on two planes:

**Data Plane** — standard HTTP/gRPC reverse proxy. Handles auth, rate limiting, and request routing. Designed to sit on top of Envoy — not compete with it.

**Control Plane (L2)** — persistent supervision channel between the gateway and the orchestrator. Monitors agent state, detects behavioral anomalies, enforces policy consistency, and initiates cryptographic challenges when an agent goes silent or diverges from its known policy.

```
┌─────────────────────────────────────────────────┐
│                  WORLD EXTERNAL                 │
└──────────────────────┬──────────────────────────┘
                       │ HTTPS / gRPC
                       ▼
┌──────────────────────────────────────────────────┐
│                   SAGEWAY                        │
│  Auth · Rate Limit · Audit · Policy Enforcement  │
└────────┬─────────────────────────────┬───────────┘
         │ Control Plane (L2)          │ Data Plane
         │ gRPC bidirectional stream   │ HTTP/gRPC proxy
         ▼                             ▼
┌─────────────────┐       ┌────────────────────────┐
│  ORCHESTRATOR   │       │    INTERNAL SERVICES   │
│  L2 supervision │◄─────►│  [AI Agents]  [APIs]   │
│  Policy distrib │       └────────────────────────┘
└─────────────────┘
```

The L2 state machine tracks every agent through 7 operational states:

```
CONNECTED → STALE → CHALLENGED → UNRESPONSIVE → QUARANTINED → RECOVERING → CONNECTED
```

A compromised agent that stops responding to cryptographic challenges is quarantined — not just rate-limited.

---

## Theoretical foundation

Sageway is the reference implementation of the **Action Claim** governance model, introduced in:

> Carvalho, R.K.S. *Toward an Operational Ontology of Agentic Action*. Zenodo, 2026.
> DOI: [10.5281/zenodo.18930044](https://doi.org/10.5281/zenodo.18930044)

The paper argues that the correct pre-execution governance object for agentic systems is not an access control decision but a structured claim — carrying declared intent, derived impact, and delegation chain — that a governance layer can evaluate for proportionality before any world-state change occurs.

Sageway implements that governance layer.

---

## Status

🚧 **Active development — not production ready**

Implemented via phased PRs. See [CHANGELOG.md](CHANGELOG.md) for current state.

---

## Getting started

```bash
git clone https://github.com/RudsonCarvalho/sageway
cd sageway
make dev        # starts etcd, Elasticsearch, OPA via docker-compose
make test       # runs all tests with race detector
make build      # builds gateway and orchestrator binaries
```

Requires: Go 1.22+, Docker, docker-compose.

Full setup guide: [docs/getting-started.md](docs/getting-started.md)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). All PRs must pass the full CI pipeline — build, tests with `-race`, golangci-lint, gosec, and govulncheck — before review.

---

## License

[Apache 2.0](LICENSE)
