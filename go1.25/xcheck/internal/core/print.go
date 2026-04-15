package core

import (
	"strings"

	"github.com/Deimvis/go-ext/go1.25/xfmt"
)

// TODO: write default formatter for each standard type (maybe using xstrconv.Formatter)
// NOTE: current default is fmt.Sprintf("%v", v)

type PrintOption func(*printConfig)

func PrintValues() PrintOption {
	return func(c *printConfig) {
		c.showValues = true
	}
}

// TODO: add PrintOption to annotate values (actual, expected, etc)
// 1) xcheck.PrintValues(xcheckfmt.Annotate("act", "exp"))
// 2) xcheck.PrintAnnotatedValues("act", "exp")
// 3) xcheckfmt.AnnotatedValues("act", "exp")

type printConfig struct {
	defaultBaseMsg string
	showValues     bool
	// TODO: pass annotations to fmtValues func (actual, expected, etc)
	fmtValues    func() string
	fmtValues_V2 func(parameters map[predArgPos]parameterFormat) string // parameters allow to parametrize arguments of predicate
}

func FormatMsg(cfg printConfig, msgAndArgsAndOptions ...any) string {
	msgAndArgs := make([]any, 0, len(msgAndArgsAndOptions))
	for _, v := range msgAndArgsAndOptions {
		if opt, ok := v.(PrintOption); ok {
			opt(&cfg)
		} else if opt, ok := v.(func() PrintOption); ok {
			opt()(&cfg)
		} else {
			msgAndArgs = append(msgAndArgs, v)
		}
	}
	baseMsg := xfmt.Sprintfg(msgAndArgs...)
	if baseMsg == "" {
		baseMsg = cfg.defaultBaseMsg
	}
	b := strings.Builder{}
	_, err := b.WriteString(baseMsg)
	if err != nil {
		panic(err)
	}

	if cfg.showValues && cfg.fmtValues != nil {
		// TODO: support NamedArgs ([]string) option for fmtValues
		_, err = b.WriteString(": " + cfg.fmtValues())
		if err != nil {
			panic(err)
		}
	}

	return b.String()
}
