package core

import "reflect"

func ReflectNil(v any, msgAndArgsAndOpts ...any) (bool, string) {
	rv := reflect.ValueOf(v)
	if !rv.IsNil() {
		pcfg := printConfig{
			defaultBaseMsg: "reflect not nil",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}

func ReflectNotNil(v any, msgAndArgsAndOpts ...any) (bool, string) {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		pcfg := printConfig{
			defaultBaseMsg: "reflect nil",
			fmtValues:      nil,
			showValues:     false,
		}
		msg := FormatMsg(pcfg, msgAndArgsAndOpts...)
		return false, msg
	}
	return true, ""
}
