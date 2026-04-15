package xruntime

import (
	"errors"
	"runtime"
	"strings"
	"sync"
)

type CodePosition interface {
	// File is an absolute path to source file on fs used during build time.
	File() string
	Line() int
}

type CallPosition interface {
	CodePosition
	// PackagePath returns "" if no module information
	// Package path has a format: {module path}/{package folder path} (see https://go.dev/ref/mod#modules-overview).
	// Note that "main" package doesn't have a distinct, universal package path,
	// it always equals "main" without including module path and subfolder.
	PackagePath() string
}

func XCaller(skip int) (CallPosition, error) {
	pc, file, line, ok := runtime.Caller(1 + skip)
	if !ok {
		return nil, errors.New("not possible recover caller info")
	}
	cp := &callPosition{
		pc:   pc,
		file: file,
		line: line,

		fn_lazy: sync.OnceValue(func() *runtime.Func {
			return runtime.FuncForPC(pc)
		}),
	}
	return cp, nil
}

type callPosition struct {
	pc   uintptr
	file string
	line int

	fn_lazy func() *runtime.Func
}

func (cp *callPosition) File() string {
	return cp.file
}

func (cp *callPosition) Line() int {
	return cp.line
}

func (cp *callPosition) PackagePath() string {
	fn := cp.fn_lazy()
	fnName := fn.Name()
	slashInd := strings.LastIndexByte(fnName, '/')
	ind := strings.IndexByte(fnName[slashInd+1:], '.')
	if ind == -1 {
		return ""
	}
	return fnName[:slashInd+1+ind]
}
