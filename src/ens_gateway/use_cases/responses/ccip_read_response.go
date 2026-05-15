package responses

// CcipReadResponse is the ERC-3668 gateway response body: the 0x-hex-encoded
// data the CCIP-Read client passes to the resolver's callback.
type CcipReadResponse struct {
	Data string `json:"data"`
}
