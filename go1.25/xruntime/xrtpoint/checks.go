package xrtpoint

import "github.com/Deimvis/go-ext/go1.25/xcheck/xmust"

func MustEnabled() {
	xmust.True(Enabled(), "xrtpoint: not enabled")
}
