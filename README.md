# agentic-ens-gateway

ERC-3668 / CCIP-Read gateway for agent ENS records.

Agent text records (`agent.runtime-url`, `agent.runtime-kind`) are stored
on-chain in the `AgentEnsRegistry` contract on the Polygon-side chain. On
Ethereum L1, the `AgentOffchainResolver` contract reverts with `OffchainLookup`
per ERC-3668, and a CCIP-Read client calls this gateway. The gateway:

1. Decodes the resolver's `IResolverService.resolve(bytes name, bytes data)`
   call data, then the inner `text(bytes32 node, string key)` query.
2. Reads `AgentEnsRegistry.text(node, key)` on the Polygon-side chain.
3. Returns an ENS-format **signed** response — `abi.encode(result, expires,
   sig)` — that the L1 resolver's `SignatureVerifier` authenticates.

The signing key lives in HashiCorp Vault. Architecture mirrors
`agentic-agent-node` (hexagonal modules, `go.uber.org/dig` DI, `koanf` config,
Fiber host, `when_*_test.go` tests).

## Layout

- `src/ens_gateway` — core domain (codec, signer, response builder, use case).
- `src/ens_gateway_blockchain` — Polygon RPC client + `AgentEnsRegistry` reader.
- `src/ens_gateway_vault` — Vault-backed signer key storage.
- `src/ens_gateway_mvc` — Fiber controllers (CCIP-Read endpoint + health).
- `src/ens_gateway_host` — HTTP host entrypoint, DI wiring, server config.

## Run

```sh
cp src/ens_gateway_host/appsettings.example.json src/ens_gateway_host/appsettings.json
# fill blockchain.rpc_http_url, vault.*, and the registry address for your chain
cd src/ens_gateway_host && go run .
```
