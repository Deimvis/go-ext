//go:build !debug

package xrtpoint

import "errors"

func NewInjector() Injector {
	return noopInector
}

func NewInjectionRuleBuilder() InjectionRuleBuilder {
	panic(errUseInRelease)
}

var (
	GlobalInjector Injector = noopInector
)

var (
	errUseInRelease = errors.New("xrtpoint: use in release")
)
