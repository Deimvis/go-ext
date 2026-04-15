package xembed

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Deimvis/go-ext/go1.25/xcheck"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

//go:embed test_data
var testData embed.FS

func TestRelativeFs(t *testing.T) {
	open := func(rfs RelativeFs, name string) error {
		_, err := rfs.Open(name)
		return err
	}
	readFile := func(rfs RelativeFs, name string) error {
		_, err := rfs.ReadFile(name)
		return err
	}
	readDir := func(rfs RelativeFs, name string) error {
		_, err := rfs.ReadDir(name)
		return err
	}

	tcs := []struct {
		title string
		op    func(rfs RelativeFs, name string) error
	}{
		{
			"open",
			open,
		},
		{
			"read-file",
			readFile,
		},
		{
			"read-dir",
			readDir,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			op := tc.op
			isDirOp := tc.title == "read-dir"
			fileOrDir := func(file string, dir string) string {
				if isDirOp {
					return dir
				}
				return file
			}
			type cds = []string
			type expCwds = []string

			tcs := []struct {
				title   string
				cds     []string
				expCwds expCwds
				name    string
			}{
				{
					"one-abs-cd",
					cds{"/test_data/dir1"},
					expCwds{"/test_data/dir1"},
					fileOrDir("a", "subdir"),
				},
				{
					"one-abs-cd/relative-name",
					cds{"/test_data/dir1/subdir"},
					expCwds{"/test_data/dir1/subdir"},
					fileOrDir("../a", "../subdir"),
				},
				{
					"one-rel-cd",
					cds{"test_data/dir1"},
					expCwds{"/test_data/dir1"},
					fileOrDir("a", "subdir"),
				},
				{
					"two-rel-cds",
					cds{"test_data", "dir1"},
					expCwds{"/test_data", "/test_data/dir1"},
					fileOrDir("a", "subdir"),
				},
				{
					"many-cds",
					cds{"test_data", "dir1", "..", "/", "/test_data/dir1"},
					expCwds{"/test_data", "/test_data/dir1", "/test_data", "/", "/test_data/dir1"},
					fileOrDir("a", "subdir"),
				},
			}
			for _, tc := range tcs {
				xmust.Eq(len(tc.cds), len(tc.expCwds), "test case not valid", xcheck.PrintWhy())
				t.Run(tc.title, func(t *testing.T) {
					rfs := NewRelativeFs(testData)

					err := op(rfs, tc.name)
					pe := &fs.PathError{}
					require.ErrorAs(t, err, &pe)

					require.Equal(t, "/", rfs.CWD())
					for i := range tc.cds {
						rfs.Cd(tc.cds[i])
						require.Equal(t, tc.expCwds[i], rfs.CWD())
					}

					err = op(rfs, tc.name)
					require.NoError(t, err)
				})
			}
		})
	}
	t.Run("one-abs-cd/open", func(t *testing.T) {
		rfs := NewRelativeFs(testData)

		_, err := rfs.Open("a")
		pe := &fs.PathError{}
		require.ErrorAs(t, err, &pe)

		rfs.Cd("test_data/dir1")
		_, err = rfs.Open("a")
		require.NoError(t, err)
		require.Equal(t, "/test_data/dir1", rfs.CWD())

		_, err = testData.Open("a")
		pe = &fs.PathError{}
		require.ErrorAs(t, err, &pe)
	})
	t.Run("one-abs-cd/read-file", func(t *testing.T) {
		rfs := NewRelativeFs(testData)

		_, err := rfs.ReadFile("a")
		pe := &fs.PathError{}
		require.ErrorAs(t, err, &pe)

		rfs.Cd("test_data/dir1")
		_, err = rfs.ReadFile("a")
		require.NoError(t, err)
		require.Equal(t, "/test_data/dir1", rfs.CWD())

		_, err = testData.ReadFile("a")
		pe = &fs.PathError{}
		require.ErrorAs(t, err, &pe)
	})
	t.Run("one-cd/read-dir", func(t *testing.T) {
		rfs := NewRelativeFs(testData)

		_, err := rfs.ReadDir("dir1")
		pe := &fs.PathError{}
		require.ErrorAs(t, err, &pe)

		rfs.Cd("test_data")
		_, err = rfs.ReadDir("dir1")
		require.NoError(t, err)
		require.Equal(t, "/test_data", rfs.CWD())

		_, err = testData.ReadDir("dir1")
		pe = &fs.PathError{}
		require.ErrorAs(t, err, &pe)
	})
}
