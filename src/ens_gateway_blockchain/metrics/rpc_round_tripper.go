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

	// JSON-RPC errors return HTTP 200 and are not caught here.
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}
	t.metrics.RecordRequest(req.Context(), t.chainId, req.URL.String(), time.Since(start), statusCode, err)

	return resp, err
}
