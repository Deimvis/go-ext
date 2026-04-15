package xrtpointss

import (
	"slices"

	"github.com/Deimvis/go-ext/go1.25/xiter"
	"github.com/Deimvis/go-ext/go1.25/xiter/xbooliter"
	"github.com/Deimvis/go-ext/go1.25/xruntime/xrtpoint"
	"github.com/Deimvis/go-ext/go1.25/xslices"
)

func MatchPointTags(tags []string) xrtpoint.MatchFn {
	return func(p xrtpoint.PointConst) bool {
		pointTags := xslices.BindComparable(p.Tags())
		return xiter.Reduce(
			xiter.Map(slices.Values(tags), pointTags.Has),
			xbooliter.All,
			true,
		)
	}
}
