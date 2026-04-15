package xrtpoint

import (
	"github.com/Deimvis/go-ext/go1.25/xchaincall/xwrapcall"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

type InjectionRule interface {
	Match(PointConst) bool
	Action() Action
	// NOTE: Action() looks weird, because it will always be called as Action()(c).
	//       But it looks even more weird when __rule__ has Action(c) method.
	//       It's more natural that rule knows when to match point with Match()
	//       and stores Action() that can be called: Action()(c).
}

type MatchFn func(PointConst) bool

type InjectionRuleBuilder interface {
	Match(MatchFn) InjectionRuleBuilder
	With(...Middleware) InjectionRuleBuilder
	Do(Action) InjectionRule
}

type injectionRule struct {
	match func(PointConst) bool
	a     Action
}

func (ir *injectionRule) Match(p PointConst) bool {
	return ir.match(p)
}

func (ir *injectionRule) Action() Action {
	return ir.a
}

type injectionRuleBuilder struct {
	rule *injectionRule
	mws  []Middleware
}

func (irb *injectionRuleBuilder) Match(fn MatchFn) InjectionRuleBuilder {
	xmust.True(irb.rule.match == nil, "match already set")
	irb.rule.match = fn
	return irb
}

func (irb *injectionRuleBuilder) With(mws ...Middleware) InjectionRuleBuilder {
	irb.mws = append(irb.mws, mws...)
	return irb
}

func (irb *injectionRuleBuilder) Do(a Action) InjectionRule {
	xmust.True(irb.rule.a == nil, "action already set")
	irb.rule.a = a
	return irb.build()
}

func (irb *injectionRuleBuilder) build() InjectionRule {
	xmust.True(irb.rule.match != nil, "match not set")
	xmust.True(irb.rule.a != nil, "action not set")
	if len(irb.mws) > 0 {
		gmws := make([]xwrapcall.Middleware[Context], len(irb.mws))
		for i, mw := range irb.mws {
			gmws[i] = mw.Generalize()
		}
		ga := irb.rule.a.Generalize()
		ga = xwrapcall.New[Context]().With(gmws...).Do(ga)
		irb.rule.a = ga.Silence()
	}
	return irb.rule
}
