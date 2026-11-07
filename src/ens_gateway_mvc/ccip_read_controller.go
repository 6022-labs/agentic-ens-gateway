package ens_gateway_mvc

import (
	"strings"

	"github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"github.com/6022-labs/agentic-ens-gateway/src/common/mvc"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const (
	CcipReadGetEndpoint  = "/:sender/:data"
	CcipReadPostEndpoint = "/:sender"
)

type CcipReadController struct {
	errorWriter          *mvc.ErrorWriter
	resolveOffchainQuery *use_cases.ResolveOffchainQuery
}

func NewCcipReadController(
	errorWriter *mvc.ErrorWriter,
	resolveOffchainQuery *use_cases.ResolveOffchainQuery,
) *CcipReadController {
	return &CcipReadController{errorWriter: errorWriter, resolveOffchainQuery: resolveOffchainQuery}
}

func (c *CcipReadController) RegisterRoutes(app fiber.Router) {
	app.Get(CcipReadGetEndpoint, c.HandleGet)
	app.Post(CcipReadPostEndpoint, c.HandlePost)
}

// HandleGet resolves a CCIP-Read query passed entirely in the path.
// @Summary CCIP-Read gateway (GET)
// @Description Resolves an ERC-3668 OffchainLookup. `sender` is the L1 resolver address;
// @Description `data` is the 0x-hex resolver call data, optionally suffixed with `.json`.
// @Tags CCIP-Read
// @Produce json
// @Param sender path string true "L1 resolver address"
// @Param data path string true "0x-hex resolver call data (optional `.json` suffix)"
// @Success 200 {object} responses.CcipReadResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.InternalError
// @Router /{sender}/{data} [get]
func (c *CcipReadController) HandleGet(ctx fiber.Ctx) error {
	request := &requests.ResolveOffchainQueryRequest{}
	if err := ctx.Bind().URI(request); err != nil {
		return c.errorWriter.Write(ctx, errors.NewValidationError("data", "invalid call data hex"))
	}

	// The GET route carries the call data in the path, where an optional `.json` suffix is allowed.
	request.Data = strings.TrimSuffix(request.Data, ".json")

	return c.resolve(ctx, request)
}

// HandlePost resolves a CCIP-Read query whose call data arrives in the body.
// @Summary CCIP-Read gateway (POST)
// @Description Resolves an ERC-3668 OffchainLookup. `sender` is the L1 resolver address;
// @Description the body carries the 0x-hex resolver call data.
// @Tags CCIP-Read
// @Accept json
// @Produce json
// @Param sender path string true "L1 resolver address"
// @Param body body requests.ResolveOffchainQueryRequest true "resolver call data"
// @Success 200 {object} responses.CcipReadResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.InternalError
// @Router /{sender} [post]
func (c *CcipReadController) HandlePost(ctx fiber.Ctx) error {
	request := &requests.ResolveOffchainQueryRequest{}
	if err := ctx.Bind().Body(request); err != nil {
		return c.errorWriter.Write(ctx, errors.NewValidationError("body", "invalid request body"))
	}

	request.Sender = ctx.Params("sender")

	return c.resolve(ctx, request)
}

func (c *CcipReadController) resolve(ctx fiber.Ctx, request *requests.ResolveOffchainQueryRequest) error {
	response, err := c.resolveOffchainQuery.Execute(ctx.Context(), request)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}

	return ctx.JSON(response)
}
