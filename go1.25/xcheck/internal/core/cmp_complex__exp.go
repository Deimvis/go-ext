package core

import "fmt"

// TODO: AllOf, NoneOf
// TODO: consider adding shortcuts like EqAnyOf, NilAnyOf, NotNilAnyOf, ...

// TODO: reimplement for UnaryPredicate and implement binding to cast binary predicates to unary (more likely with method in BinaryPredicate)
func AnyOf[T any](pred BinaryPredicate_[T], v1 T, v2opts []T, msgAndArgsAndOpts ...any) (bool, string) {
	for _, v2 := range v2opts {
		if pred.call(v1, v2) {
			return true, ""
		}
	}
	var placeholder T
	subpcfg := pred.defaultPrintConfig(v1, placeholder)
	if subpcfg.fmtValues_V2 == nil {
		panic("not implemented")
	}
	pcfg := printConfig{
		defaultBaseMsg: `any_of not met`,
		fmtValues: func() string {
			return subpcfg.fmtValues_V2(map[predArgPos]parameterFormat{
				predArg2: newParameterFormat(`<x>`, fmt.Sprintf("belongs to %v", v2opts)),
			})
		},
		showValues: false,
	}
	msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
	return false, msg
}
