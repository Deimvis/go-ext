package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatmsg_Options(t *testing.T) {
	basecfg := printConfig{
		defaultBaseMsg: "check failed",
		fmtValues:      func() string { return "myvalue" },
	}
	t.Run("PrintValues", func(t *testing.T) {
		t.Run("applied/simple", func(t *testing.T) {
			cfg := basecfg
			cfg.showValues = false
			res := FormatMsg(cfg, "mymsg", PrintValues())
			require.Equal(t, "mymsg: myvalue", res)
		})
		t.Run("applied/already-set", func(t *testing.T) {
			cfg := basecfg
			cfg.showValues = true
			res := FormatMsg(cfg, "mymsg", PrintValues())
			require.Equal(t, "mymsg: myvalue", res)
		})
		t.Run("applied/lazy-opt", func(t *testing.T) {
			cfg := basecfg
			cfg.showValues = false
			res := FormatMsg(cfg, "mymsg", PrintValues)
			require.Equal(t, "mymsg: myvalue", res)
		})
		t.Run("applied/within-args", func(t *testing.T) {
			cfg := basecfg
			cfg.showValues = false
			res := FormatMsg(cfg, "mymsg: %s %s", "start", PrintValues(), "end")
			require.Equal(t, "mymsg: start end: myvalue", res)
		})
	})
}
