package ens_gateway_mvc

import (
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/responses"
	"github.com/gofiber/fiber/v2"
)

const HealthEndpoint = "/health"

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) RegisterRoutes(app fiber.Router) {
	app.Get(HealthEndpoint, c.GetHealth)
}

// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} responses.GetHealthResponse
// @Router /health [get]
func (c *HealthController) GetHealth(ctx *fiber.Ctx) error {
	return ctx.JSON(&responses.GetHealthResponse{Status: "ok"})
}
