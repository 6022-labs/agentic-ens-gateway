package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/responses"
)

// GetHealth reports that the gateway process is serving requests.
type GetHealth struct{}

func NewGetHealth() *GetHealth {
	return &GetHealth{}
}

func (u *GetHealth) Execute(_ context.Context) (*responses.GetHealthResponse, error) {
	return &responses.GetHealthResponse{Status: "ok"}, nil
}
