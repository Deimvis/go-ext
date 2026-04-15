package xrtpoint

import "context"

// Injector injects runtime points
// and configures what they do on Call().
// Injector and all subsequent types are NOT thread-safe.
type Injector interface {
	NewPoint() Point
	Rules() OrderedInjectionRules
}

// OrderedInjectionRules represents ordered list of rules.
// On each Point Call() rules are traversed
// in the original order and all matched rules execute their Action().
type OrderedInjectionRules interface {
	PushFront(...InjectionRule)
	PushBack(...InjectionRule)
	Clear()

	ToSlice() []InjectionRule
}

type injector struct {
	rules OrderedInjectionRules
}

func (i *injector) NewPoint() Point {
	callFn := func(p PointConst) {
		c := newContext(context.Background(), p)
		for _, r := range i.rules.ToSlice() {
			if r.Match(p) {
				r.Action()(c)
			}
		}
	}
	return newPoint(i, callFn)
}

func (i *injector) Rules() OrderedInjectionRules {
	return i.rules
}

type orderedInjectionRules struct {
	s []InjectionRule
}

func (oir *orderedInjectionRules) PushFront(rules ...InjectionRule) {
	oir.s = append(rules, oir.s...)
}

func (oir *orderedInjectionRules) PushBack(rules ...InjectionRule) {
	oir.s = append(oir.s, rules...)
}

func (oir *orderedInjectionRules) Clear() {
	oir.s = oir.s[:0]
}

func (oir *orderedInjectionRules) ToSlice() []InjectionRule {
	return oir.s
}
