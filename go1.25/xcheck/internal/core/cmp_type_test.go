package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeImplements(t *testing.T) {
	type C struct {
		context.Context
	}
	ok, _ := TypeImplements[C, context.Context]()
	require.True(t, ok)

	ok, _ = TypeImplements[int, context.Context]()
	require.False(t, ok)
}
