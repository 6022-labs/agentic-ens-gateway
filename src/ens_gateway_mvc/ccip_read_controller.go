package ens_gateway_mvc

import (
	"strings"

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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{sender} [post]
func (c *CcipReadController) HandlePost(ctx *fiber.Ctx) error {
	body := new(ccipReadPostBody)
	if err := ctx.BodyParser(body); err != nil {
		return badRequest(ctx, "invalid request body")
	}
	return c.resolve(ctx, ctx.Params("sender"), body.Data)
}

func (c *CcipReadController) resolve(ctx *fiber.Ctx, senderParam, dataHex string) error {
	if !common.IsHexAddress(senderParam) {
		return badRequest(ctx, "invalid sender address")
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return badRequest(ctx, "invalid call data hex")
	}

	response, err := c.resolveOffchainQuery.Execute(ctx.UserContext(), &requests.ResolveOffchainQueryRequest{
		Sender: common.HexToAddress(senderParam),
		Data:   data,
	})
	if err != nil {
		c.logger.Error("failed to resolve offchain query", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "failed to resolve query"})
	}

	return ctx.JSON(response)
}

func badRequest(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": message})
}
