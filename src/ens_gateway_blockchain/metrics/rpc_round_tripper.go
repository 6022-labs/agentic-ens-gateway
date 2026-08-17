package metrics

import (
	"net/http"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_blockchain/metrics/interfaces"
)

// RpcRoundTripper records rpc.client.* per JSON-RPC call, aggregated per chain (the method is in the body).
type RpcRoundTripper struct {
	base    http.RoundTripper
	metrics metrics_interfaces.BlockchainRpcMetricsInterface
	chainId string
}

func NewRpcRoundTripper(
	base http.RoundTripper,
	metrics metrics_interfaces.BlockchainRpcMetricsInterface,
	chainId string,
) *RpcRoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &RpcRoundTripper{
		base:    base,
		metrics: metrics,
		chainId: chainId,
	}
}

func (t *RpcRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.base.RoundTrip(req)

	// Count transport errors and 4xx/5xx; JSON-RPC errors return HTTP 200 and are not caught here.
	failed := err != nil || (resp != nil && resp.StatusCode >= 400)
	t.metrics.RecordRequest(req.Context(), t.chainId, time.Since(start), failed)

	return resp, err
}
