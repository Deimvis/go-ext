package xerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIs(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{
			errors.New("some error"),
			false,
		},
		{
			customError{Value: "custom error"},
			true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := Is[customError](tc.err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

type customError struct{ Value string }

func (e customError) Error() string { return e.Value }
