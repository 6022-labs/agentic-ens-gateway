package ens_gateway_mvc

import (
	"strings"

	"github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CcipReadController struct {
	logger               *zap.Logger
	resolveOffchainQuery use_cases.ResolveOffchainQueryInterface
}

func NewCcipReadController(
	logger *zap.Logger,
	resolveOffchainQuery use_cases.ResolveOffchainQueryInterface,
) *CcipReadController {
	return &CcipReadController{logger: logger, resolveOffchainQuery: resolveOffchainQuery}
}

func (c *CcipReadController) RegisterRoutes(app fiber.Router) {
	app.Get("/:sender/:data", c.HandleGet)
	app.Post("/:sender", c.HandlePost)
}

type ccipReadPostBody struct {
	Data string `json:"data"`
}

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
	data := strings.TrimSuffix(ctx.Params("data"), ".json")
	return c.resolve(ctx, ctx.Params("sender"), data)
}

// @Summary CCIP-Read gateway (POST)
// @Tags CCIP-Read
// @Accept json
// @Produce json
// @Param sender path string true "L1 resolver address"
// @Param body body ccipReadPostBody true "resolver call data"
// @Success 200 {object} responses.CcipReadResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.InternalError
// @Router /{sender} [post]
func (c *CcipReadController) HandlePost(ctx *fiber.Ctx) error {
	body := new(ccipReadPostBody)
	if err := ctx.BodyParser(body); err != nil {
		return c.writeError(ctx, errors.NewValidationError("body", "invalid request body"))
	}
	return c.resolve(ctx, ctx.Params("sender"), body.Data)
}

func (c *CcipReadController) resolve(ctx *fiber.Ctx, senderParam, dataHex string) error {
	if !common.IsHexAddress(senderParam) {
		return c.writeError(ctx, errors.NewValidationError("sender", "invalid sender address"))
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return c.writeError(ctx, errors.NewValidationError("data", "invalid call data hex"))
	}

	response, err := c.resolveOffchainQuery.Execute(ctx.UserContext(), &requests.ResolveOffchainQueryRequest{
		Sender: common.HexToAddress(senderParam),
		Data:   data,
	})
	if err != nil {
		c.logger.Error("failed to resolve offchain query", zap.Error(err))
		return c.writeError(ctx, err)
	}

	return ctx.JSON(response)
}

// writeError maps typed errors to statuses; anything unknown collapses to a detail-free 500.
func (c *CcipReadController) writeError(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *errors.ValidationError:
		return ctx.Status(fiber.StatusBadRequest).JSON(e)
	default:
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(errors.NewInternalError("internal_error", "failed to resolve query"))
	}
}
