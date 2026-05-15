package ens_gateway_mvc

import "github.com/gofiber/fiber/v2"

// ControllerInterface registers a controller's HTTP routes. The gateway serves
// public ERC-3668 endpoints, so there is no authenticated/unauthenticated split.
type ControllerInterface interface {
	RegisterRoutes(app fiber.Router)
}
