package requests_test

import (
	"testing"

	apperrors "github.com/6022-labs/agentic-ens-gateway/src/common/errors"
	"github.com/6022-labs/agentic-ens-gateway/src/ens_gateway/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

const validSender = "0x1111111111111111111111111111111111111111"

func TestWhenValidatingAResolveOffchainQueryRequest(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid sender and call data", func(t *testing.T) {
		t.Parallel()

		t.Run("Should expose the parsed address and decoded call data", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: validSender, Data: "0x010203"}

			assert.NoError(t, request.ValidateAndSanitize())
			assert.Equal(t, validSender, request.SenderAddress().Hex())
			assert.Equal(t, []byte{0x01, 0x02, 0x03}, request.CallData())
		})
	})

	t.Run("Given a sender that is not a hex address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail on the sender field with the message the gateway has always returned", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: "not-an-address", Data: "0x010203"}

			err := request.ValidateAndSanitize()

			var validationError *apperrors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "sender", validationError.Field)
			assert.Equal(t, "invalid sender address", validationError.Message)
		})
	})

	t.Run("Given call data that is not hex", func(t *testing.T) {
		t.Parallel()

		t.Run("Should fail on the data field with the message the gateway has always returned", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: validSender, Data: "not-hex"}

			err := request.ValidateAndSanitize()

			var validationError *apperrors.ValidationError
			assert.ErrorAs(t, err, &validationError)
			assert.Equal(t, "data", validationError.Field)
			assert.Equal(t, "invalid call data hex", validationError.Message)
		})
	})

	t.Run("Given a sender padded with whitespace", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it, because main never trimmed the inputs", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: " " + validSender, Data: "0x010203"}

			assert.Error(t, request.ValidateAndSanitize())
		})
	})

	t.Run("Given call data carrying a .json suffix", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it, because only the GET route strips the suffix", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: validSender, Data: "0x010203.json"}

			assert.Error(t, request.ValidateAndSanitize())
		})
	})

	t.Run("Given the sender is validated before the call data", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report the sender when both are invalid, as main did", func(t *testing.T) {
			t.Parallel()

			request := &requests.ResolveOffchainQueryRequest{Sender: "nope", Data: "also-nope"}

			var validationError *apperrors.ValidationError
			assert.ErrorAs(t, request.ValidateAndSanitize(), &validationError)
			assert.Equal(t, "sender", validationError.Field)
		})
	})
}
