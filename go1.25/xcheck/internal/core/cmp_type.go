package core

import "fmt"

func TypeImplements[T any, I any](msgAndArgsAndOpts ...any) (bool, string) {
	var v T
	_, ok := any(v).(I)
	if !ok {
		pcfg := printConfig{
			defaultBaseMsg: "type not implements",
			fmtValues: func() string {
				var i I
				return fmt.Sprintf("%T not implements %T", v, i)
			},
			showValues: false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}
