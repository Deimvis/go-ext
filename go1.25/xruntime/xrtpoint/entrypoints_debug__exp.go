//go:build debug

package xrtpoint

func NewInjector() Injector {
	return &injector{
		rules: &orderedInjectionRules{},
	}
}

func NewInjectionRuleBuilder() InjectionRuleBuilder {
	return &injectionRuleBuilder{rule: &injectionRule{}}
}

var (
	GlobalInjector Injector = NewInjector()
)
