//go:build debug

package core

import (
	"fmt"
	"strings"
)

func _validateMsg(msgAndArgs ...any) {
	msg := FormatMsg(printConfig{showValues: false}, msgAndArgs...)
	if strings.Contains(msg, "%!") {
		panic(fmt.Sprintf("formatted msg: %s", msg))
	}
}
