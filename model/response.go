package model

const (
	ErrDecodeParams  = "invalid query parameter - currency must be provided"
	ErrConvert       = "error converting currency"
	ErrRouteNotFound = "route not found"
)

// ConvertResp is the response struct for XE Service
type ConvertResp struct {
	From           string  `json:"from,omitempty"`
	To             string  `json:"to,omitempty"`
	Rate           float64 `json:"rate,omitempty"`
	Recommendation string  `json:"recommendation,omitempty"`
	Error          string  `json:"error,omitempty"`
}
