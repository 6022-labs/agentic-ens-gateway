package controllers_test

import (
	"encoding/json"
	"errors"

	apperrors "github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gateway_responses "github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/responses"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway_mvc"
	use_cases_mocks "github.com/6022-labs/agentic-ens-gateway/tests/ens_gateway_mocks/use_cases_mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	testSender  = "0x1111111111111111111111111111111111111111"
	testDataHex = "0x010203"
)

type WhenHandlingCcipReadRequestTestingSuite struct {
	mockResolveOffchainQuery *use_cases_mocks.MockResolveOffchainQueryInterface
	app                      *fiber.App
}

func WhenHandlingCcipReadRequestBeforeEach(t *testing.T) *WhenHandlingCcipReadRequestTestingSuite {
	ctrl := gomock.NewController(t)
	mockResolveOffchainQuery := use_cases_mocks.NewMockResolveOffchainQueryInterface(ctrl)

	controller := ens_gateway_mvc.NewCcipReadController(zap.NewNop(), mockResolveOffchainQuery)
	app := fiber.New()
	controller.RegisterRoutes(app)

	return &WhenHandlingCcipReadRequestTestingSuite{mockResolveOffchainQuery: mockResolveOffchainQuery, app: app}
}

func TestWhenHandlingCcipReadRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid GET request with a .json suffix", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			suite.mockResolveOffchainQuery.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Return(&gateway_responses.CcipReadResponse{Data: "0xbeef"}, nil)
		}

		t.Run("Should strip .json, call the use case, and return 200 with the data", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, `{"data":"0xbeef"}`, string(body))
		})
	})

	t.Run("Given a valid POST request", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			suite.mockResolveOffchainQuery.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Return(&gateway_responses.CcipReadResponse{Data: "0xbeef"}, nil)
		}

		t.Run("Should call the use case and return 200", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+testSender,
				strings.NewReader(`{"data":"`+testDataHex+`"}`),
			)
			req.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Given malformed call data hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without calling the use case", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/not-hex.json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given an invalid sender address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without calling the use case", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(http.MethodGet, "/not-an-address/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given a POST request with an invalid body", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return 400 without calling the use case", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+testSender,
				strings.NewReader("not json"),
			)
			req.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Given the use case fails with a validation error", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			suite.mockResolveOffchainQuery.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Return(nil, apperrors.NewValidationError("data", "invalid resolver call data"))
		}

		t.Run("Should return 400 with the field and message", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var body map[string]string
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
			assert.Equal(t, "data", body["field"])
			assert.Equal(t, "invalid resolver call data", body["message"])
		})
	})

	t.Run("Given the use case fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenHandlingCcipReadRequestTestingSuite) {
			suite.mockResolveOffchainQuery.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom"))
		}

		t.Run("Should return 500", func(t *testing.T) {
			t.Parallel()

			suite := WhenHandlingCcipReadRequestBeforeEach(t)
			initSuite(suite)

			req := httptest.NewRequest(http.MethodGet, "/"+testSender+"/"+testDataHex+".json", nil)
			resp, err := suite.app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})
	})
}
