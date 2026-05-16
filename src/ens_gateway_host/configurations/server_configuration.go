package configurations

import (
	_ "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_host/docs"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	fiberzap "github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const appName = "ENS-Gateway"

func ConfigureServer(container *dig.Container) {
	if err := container.Provide(newHttpServer); err != nil {
		panic(err)
	}

	if err := container.Invoke(registerRoutes); err != nil {
		panic(err)
	}
}

func newHttpServer() *fiber.App {
	return fiber.New(fiber.Config{
		AppName: appName,
		Network: "tcp",
	})
}

type registerRoutesParams struct {
	dig.In

	App         *fiber.App
	Logger      *zap.Logger
	Controllers []ens_gateway_mvc.ControllerInterface `group:"controllers"`
}

func registerRoutes(p registerRoutesParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	p.App.Use(fiberzap.New(fiberzap.Config{Logger: p.Logger}))

	p.App.Get("/swagger/*", swagger.New())

	for _, controller := range p.Controllers {
		controller.RegisterRoutes(p.App)
	}
}
