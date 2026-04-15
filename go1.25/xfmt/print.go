package xfmt

import (
	"fmt"
	"strings"
)

// Sprintfg = Sprintf + generalized to work with zero arguments (returns empty string).
// First argument must be a string as it is for Sprintf.
// It helps implement functions with optional message arguments,
// e.g. Assert(bool, msgAndArgs... any) -> calls Sprintfg(msgAndArgs...).
func Sprintfg(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	msg, ok := msgAndArgs[0].(string)
	if !ok {
		panic(fmt.Errorf("Sprintfg first argument is not string: `%+v`)", msgAndArgs[0]))
	}
	if len(msgAndArgs) == 1 {
		return msg
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msg, msgAndArgs[1:]...)
	}
	panic("bug: unreachable code has been reached")
}

// Sprintfkv formats key-value pairs into string.
// There must be no dangling key in the end (with no value), panics otherwise.
// Format string should include two %v entries (with any formatting flags)
// when key-value types aren't consistent.
// Example:
// Sprintfkv("(%v, %v)", "\t", "key", "value", "key2", 2)
// -> (key, value)\t(key2, 2)
func Sprintfkv(kvFormat string, sep string, keysAndValues ...any) string {
	if (len(keysAndValues) & 1) == 1 {
		panic(fmt.Errorf("Sprintfkv arguments have dangling key: `%v`", keysAndValues[len(keysAndValues)-1]))
	}
	s := make([]string, len(keysAndValues)/2)
	i := 0
	for i < len(keysAndValues) {
		k := keysAndValues[i]
		v := keysAndValues[i+1]
		s[i/2] = fmt.Sprintf(kvFormat, k, v)
		i += 2
	}
	return strings.Join(s, sep)
}
