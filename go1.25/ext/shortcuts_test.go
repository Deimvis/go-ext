package ext

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUntilFirstErr(t *testing.T) {
	testCases := []struct {
		fns      []func() error
		expecetd error
	}{
		{
			[]func() error{
				func() error { return nil },
			},
			nil,
		},
		{
			[]func() error{
				func() error { return nil },
				func() error { return errors.New("123") },
				func() error { return nil },
			},
			errors.New("123"),
		},
		{
			[]func() error{
				func() error { return errors.New("111") },
				func() error { return errors.New("222") },
			},
			errors.New("111"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := UntilFirstErr(tc.fns...)
			require.Equal(t, tc.expecetd, err)
		})
	}
}
