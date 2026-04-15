package xcheck

import "github.com/Deimvis/go-ext/go1.25/xcheck/internal/core"

// NOTE: very experimental, likely name and effect will be changed in future

// PrintValues returns PrintOption
// which adds values to result message.
// It works in both ways:
// - as PrintOption: `xmust.Eq(1, 2, xcheck.PrintValues())`
// - as func() PrintOption: `xmust.Eq(1, 2, xcheck.PrintValues)`
var PrintValues = core.PrintValues

// PrintWhy was renamed to PrintValues
var PrintWhy = PrintValues
