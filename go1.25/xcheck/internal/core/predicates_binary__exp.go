package core

import (
	"fmt"
	"strings"
)

// GOEXPERIMENT=aliastypeparams
// type BinaryPredicate[T any] = func(T, T) bool

// TODO: maybe keep predicate a function type, but have interface and check for implementation (e.g. BinaryPredicateObservable)
// TODO: maybe rename this to DetailedBinaryPredicate and return struct with info on call()?
// TODO: maybe support trivial binary predicate, that implements only Call() and check in runtime
// whether it is a detailed predicate with print config options customized.
// And add NewBinaryPred() function to simply create BiaryPredicate from just func(T, T) bool:
//
//	eqlen := NewBinaryPred(func(v1, v2 string) bool { return len(v1) == len(v2) })
//	AnyOf(eqlen, "123", []string{"1", "22", "333"}) -> true
type BinaryPredicate_[T any] interface {
	// HIGHLY EXPERIMENTAL STRAIGHTFORWARD DUMB INTERFACE
	call(T, T) bool
	defaultPrintConfig(T, T) printConfig
}

// type BinaryPredicate[T any] interface {
// 	Call(T, T) bool
// 	// would show either 1 != 2 or <x> != <y>, where ...
// 	// TODO: rename so that it shows that format is for unsatisfied.
//  // need to return what exact value made check to fail
// 	FormatValues(xcheckfmt.BinaryPredicateConfig) string
// 	// TODO: add method like CallWithReport() (bool, formattedReport)
// }

type EqPred[T comparable] struct{}

// interface guards
var _ BinaryPredicate_[int] = EqPred[int]{}

func (eq EqPred[T]) call(v1, v2 T) bool {
	return v1 == v2
}

func (eq EqPred[T]) defaultPrintConfig(v1, v2 T) printConfig {
	return printConfig{
		defaultBaseMsg: "not equal",
		fmtValues: func() string {
			// TODO: show diff when values are complex (slices, structs, etc)
			return fmt.Sprintf("%v != %v", v1, v2)
		},
		showValues: false,

		// VERY DUMB IMPL, replace map with more convenient type
		fmtValues_V2: func(p map[predArgPos]parameterFormat) string {
			restrictions := []string{}
			var v1String string
			if v1Fmt, ok := p[predArg1]; ok {
				v1String = v1Fmt.literal()
				restrictions = append(restrictions, v1String+" "+v1Fmt.restrictions())
			} else {
				v1String = fmt.Sprintf("%v", v1)
			}
			var v2String string
			if v2Fmt, ok := p[predArg2]; ok {
				v2String = v2Fmt.literal()
				restrictions = append(restrictions, v2String+" "+v2Fmt.restrictions())
			} else {
				v2String = fmt.Sprintf("%v", v2)
			}

			msg := fmt.Sprintf("%s != %s", v1String, v2String)
			if len(restrictions) > 0 {
				msg += ", where " + strings.Join(restrictions, ", ")
			}
			return msg
		},
	}
}
