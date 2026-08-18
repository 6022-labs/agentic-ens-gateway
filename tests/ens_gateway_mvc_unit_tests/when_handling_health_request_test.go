package ens_gateway_mvc_unit_tests_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestWhenHandlingHealthRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a GET request to /health", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 200 with status ok", func(t *testing.T) {
			t.Parallel()

			app := fiber.New()
			ens_gateway_mvc.NewHealthController(use_cases.NewGetHealth()).RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, `{"status":"ok"}`, string(body))
		})
	})
}
