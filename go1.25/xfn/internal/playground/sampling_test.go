//go:build playground

package playground

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TODO: rewrite

func TestNewSamplingLogFn(t *testing.T) {
	lgCfg := vklog.NewConfig("")
	lgCfg.Outputs = []string{}
	lg, err := vklog.New(
		lgCfg,
		&vklog.AppInfo{
			Name:        "test",
			Environment: "test",
			Release:     "test",
			Instance:    "test",
		},
	)
	require.NoError(t, err)

	w := &bufLogWriter{}
	lg.WithWriter("test_writer", w)

	interval := time.Minute
	logFn, err := newSamplingLogFn(lg.Infom, interval)
	require.NoError(t, err)

	{
		logFn("msg1", logctx{"key": "value"})
		err = lg.Flush()
		require.NoError(t, err)

		require.Len(t, w.buf, 1)
		event := w.buf[0]
		require.Equal(t, vklog.InfoLevel, event.Level)
		require.Equal(t, "msg1", event.Message)
		logFieldValue, ok := event.Fields["key"]
		require.True(t, ok)
		require.Equal(t, "value", logFieldValue)
	}
	{
		logFn("msg2", logctx{"key": "value"})
		err = lg.Flush()
		require.NoError(t, err)

		require.Len(t, w.buf, 1)
	}
	nowFn = func() time.Time {
		return time.Now().Add(interval)
	}
	{
		logFn("msg3", logctx{"key": "value"})
		err = lg.Flush()
		require.NoError(t, err)

		require.Len(t, w.buf, 2)
		event := w.buf[1]
		require.Equal(t, vklog.InfoLevel, event.Level)
		require.Equal(t, "msg3", event.Message)
		logFieldValue, ok := event.Fields["key"]
		require.True(t, ok)
		require.Equal(t, "value", logFieldValue)
	}
}

type bufLogWriter struct {
	buf []vklog.LogEvent

	mu sync.Mutex
}

func (w *bufLogWriter) Preload() {}
func (w *bufLogWriter) Write(e *vklog.LogEvent) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, *e)
	return nil
}
func (w *bufLogWriter) Sync() error {
	return nil
}
