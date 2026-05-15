# agentic-ens-gateway

ERC-3668 / CCIP-Read gateway for agent ENS records.

Agent text records (`agent.runtime-url`, `agent.runtime-kind`) are stored
on-chain in `AgentEnsRegistry` contracts on one or more EVM chains. On Ethereum
L1, the `AgentOffchainResolver` contract reverts with `OffchainLookup` per
ERC-3668, and a CCIP-Read client calls this gateway. The gateway:

1. Decodes the resolver's `IResolverService.resolve(bytes name, bytes data)`
   call data, extracts the chain id from a fixed label of the DNS-encoded
   name, then decodes the inner `text(bytes32 node, string key)` query.
2. Reads `AgentEnsRegistry.text(node, key)` on the matching chain.
3. Returns an ENS-format **signed** response — `abi.encode(result, expires,
   sig)` — that the L1 resolver's `SignatureVerifier` authenticates.

Architecture mirrors `agentic-agent-node` (hexagonal modules, `go.uber.org/dig`
DI, `koanf` config, Fiber host, `when_*_test.go` tests).

## Layout

- `src/ens_gateway` — core domain (codec, signer, response builder, use case).
- `src/ens_gateway_blockchain` — per-chain RPC clients + `AgentEnsRegistry` reader.
- `src/ens_gateway_mvc` — Fiber controllers (CCIP-Read endpoint + health).
- `src/ens_gateway_host` — HTTP host entrypoint, DI wiring, server config.

## Run

```sh
cp src/ens_gateway_host/appsettings.example.json src/ens_gateway_host/appsettings.json
# fill chains.<chain_id>.rpc_http_url and signed_response.private_key
cd src/ens_gateway_host && go run .
```
