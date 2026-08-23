package ens_gateway_mvc

import (
	"github.com/6022-labs/agentic-ens-gateway/src/common/mvc"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/gofiber/fiber/v3"
)

const HealthEndpoint = "/health"

type HealthController struct {
	getHealth *use_cases.GetHealth
}

func NewHealthController(getHealth *use_cases.GetHealth) *HealthController {
	return &HealthController{getHealth: getHealth}
}

func (c *HealthController) RegisterRoutes(app fiber.Router) {
	app.Get(HealthEndpoint, c.GetHealth)
}

// @Summary Health check
// @Description Liveness probe for the ENS gateway.
// @Tags Health
// @Produce json
// @Success 200 {object} responses.GetHealthResponse
// @Router /health [get]
func (c *HealthController) GetHealth(ctx fiber.Ctx) error {
	response, err := c.getHealth.Execute(ctx.Context())
	if err != nil {
		return mvc.WriteError(ctx, err)
	}

	return ctx.JSON(response)
}
