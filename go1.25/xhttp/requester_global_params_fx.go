package xhttp

import (
	"go.uber.org/fx"
)

// TODO: move to xhttpfx

type RequesterGlobalParamsInput struct {
	fx.In

	// TODO: rename group to "global"
	RoundTripWraps []RoundTripWrapFn `group:"round_trip_wraps"`
}

// shortcut for fx to initialize RequesterGlobalParams
func NewRequesterGlobalParams(inp RequesterGlobalParamsInput) *RequesterGlobalParams {
	return &RequesterGlobalParams{
		RoundTripWraps: inp.RoundTripWraps,
	}
}
