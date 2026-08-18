package ens_gateway_mvc

import (
	"github.com/6022-labs/agentic-ens-gateway/src/common/mvc"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	CcipReadGetEndpoint  = "/:sender/:data"
	CcipReadPostEndpoint = "/:sender"
)

type CcipReadController struct {
	logger               *zap.Logger
	resolveOffchainQuery *use_cases.ResolveOffchainQuery
}

func NewCcipReadController(
	logger *zap.Logger,
	resolveOffchainQuery *use_cases.ResolveOffchainQuery,
) *CcipReadController {
	return &CcipReadController{logger: logger, resolveOffchainQuery: resolveOffchainQuery}
}

func (c *CcipReadController) RegisterRoutes(app fiber.Router) {
	app.Get(CcipReadGetEndpoint, c.HandleGet)
	app.Post(CcipReadPostEndpoint, c.HandlePost)
}

// HandleGet resolves a CCIP-Read query passed entirely in the path.
// @Summary CCIP-Read gateway (GET)
// @Tags CCIP-Read
// @Produce json
// @Param sender path string true "L1 resolver address"
// @Param data path string true "0x-hex resolver call data (optional `.json` suffix)"
// @Success 200 {object} responses.CcipReadResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.InternalError
// @Router /{sender}/{data} [get]
func (c *CcipReadController) HandleGet(ctx *fiber.Ctx) error {
	request := &requests.ResolveOffchainQueryRequest{}
	if err := ctx.ParamsParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"code": "invalid_params", "message": err.Error()})
	}

	return c.resolve(ctx, request)
}

// HandlePost resolves a CCIP-Read query whose call data arrives in the body.
// @Summary CCIP-Read gateway (POST)
// @Tags CCIP-Read
// @Accept json
// @Produce json
// @Param sender path string true "L1 resolver address"
// @Param body body requests.ResolveOffchainQueryRequest true "resolver call data"
// @Success 200 {object} responses.CcipReadResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.InternalError
// @Router /{sender} [post]
func (c *CcipReadController) HandlePost(ctx *fiber.Ctx) error {
	request := &requests.ResolveOffchainQueryRequest{}
	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"code": "invalid_body", "message": err.Error()})
	}

	request.Sender = ctx.Params("sender")

	return c.resolve(ctx, request)
}

func (c *CcipReadController) resolve(ctx *fiber.Ctx, request *requests.ResolveOffchainQueryRequest) error {
	response, err := c.resolveOffchainQuery.Execute(ctx.UserContext(), request)
	if err != nil {
		c.logger.Error("failed to resolve offchain query", zap.Error(err))
		return mvc.WriteError(ctx, err)
	}

	return ctx.JSON(response)
}
