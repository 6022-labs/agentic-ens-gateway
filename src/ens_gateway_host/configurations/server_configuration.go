package configurations

import (
	mvc_interfaces "github.com/6022-labs/agentic-ens-gateway/src/common/mvc/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-ens-gateway/src/common/mvc/middlewares"
	"github.com/6022-labs/agentic-ens-gateway/src/common/settings"
	_ "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_host/docs"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/contrib/v3/swaggo"
	fiberzap "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const AppName = "6022-EnsGateway"

func ConfigureServer(container *dig.Container) {
	if err := container.Provide(newHttpServer); err != nil {
		panic(err)
	}

	if err := container.Invoke(registerRoutes); err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:   AppName,
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type registerRoutesParams struct {
	dig.In

	App            *fiber.App
	Logger         *zap.Logger
	RequestMetrics *mvc_middlewares.ApiRequestMetricsMiddleware
	Controllers    []mvc_interfaces.ControllerInterface `group:"controllers"`
}

func registerRoutes(p registerRoutesParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	// Spans only: ApiRequestMetrics owns http.server.*.
	p.App.Use(fiberotel.Middleware(fiberotel.WithoutMetrics(true)))
	p.App.Use(p.RequestMetrics.Handle)
	p.App.Use(fiberzap.New(fiberzap.Config{Logger: p.Logger}))

	p.App.Get("/swagger/*", swaggo.New())

	for _, controller := range p.Controllers {
		controller.RegisterRoutes(p.App)
	}
}
