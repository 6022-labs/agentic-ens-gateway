package mvc

import (
	"github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"github.com/gofiber/fiber/v3"
)

// WriteError maps a use-case error to its HTTP status. The concrete error
// types carry the JSON body ({code,message} or {field,message}); an unknown
// error collapses to 500. Shared by every *_mvc controller package.
func WriteError(ctx fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *errors.ValidationError:
		return ctx.Status(fiber.StatusBadRequest).JSON(e)
	case *errors.UnauthorizedError:
		return ctx.Status(fiber.StatusUnauthorized).JSON(e)
	case *errors.NotFoundError:
		return ctx.Status(fiber.StatusNotFound).JSON(e)
	case *errors.ConflictError:
		return ctx.Status(fiber.StatusConflict).JSON(e)
	case *errors.UnavailableError:
		return ctx.Status(fiber.StatusBadGateway).JSON(e)
	case *errors.InternalError:
		return ctx.Status(fiber.StatusInternalServerError).JSON(e)
	default:
		// These services are publicly reachable: an unrecognized error must not leak its text.
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(errors.NewInternalError("internal_error", "internal error"))
	}
}
