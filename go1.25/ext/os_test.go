package ext

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShred(t *testing.T) {
	// https://cs.opensource.google/go/go/+/master:src/os/removeall_test.go;l=33;bpv=0;bpt=0?q=TestRemove&ss=go%2Fgo
	t.Run("regular_file", func(t *testing.T) {
		fp := filepath.Join(os.TempDir(), "regular_file")
		fd, err := os.Create(fp)
		if err != nil {
			t.Fatalf("create %q: %s", fp, err)
		}
		fd.Close()
		if _, err = os.Lstat(fp); err != nil {
			t.Fatalf("Lstat %q failed after os.Create: %s", fp, err)
		}
		if err = Shred(fp); err != nil {
			t.Fatalf("Shred %q: %s", fp, err)
		}
		if _, err = os.Lstat(fp); err == nil {
			t.Fatalf("Lstat %q succeeded after Shred", fp)
		}
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		content := []byte("hello")
		fp := filepath.Join(os.TempDir(), "regular_file")
		fd, err := os.Create(fp)
		if err != nil {
			t.Fatalf("create %q: %s", fp, err)
		}
		fd.Write(content)
		fd.Close()

		fpCopy := filepath.Join(os.TempDir(), "regular_file_copy")

		err = CopyFile(fp, fpCopy)
		if err != nil {
			t.Fatalf("copy file %q to %q: %s", fp, fpCopy, err)
		}

		fdCopy, err := os.Open(fpCopy)
		if err != nil {
			t.Fatalf("open %q: %s", fpCopy, err)
		}
		contentCopy, err := io.ReadAll(fdCopy)
		if err != nil {
			t.Fatalf("read all %q: %s", fpCopy, err)
		}

		require.Equal(t, content, contentCopy)
	})
}

// https://cs.opensource.google/go/go/+/refs/tags/go1.23.3:src/os/os_test.go;l=54

type sysDir struct {
	name  string
	files []string
}

var sysdir = func() *sysDir {
	switch runtime.GOOS {
	case "android":
		return &sysDir{
			"/system/lib",
			[]string{
				"libmedia.so",
				"libpowermanager.so",
			},
		}
	case "ios":
		wd, err := syscall.Getwd()
		if err != nil {
			wd = err.Error()
		}
		sd := &sysDir{
			filepath.Join(wd, "..", ".."),
			[]string{
				"ResourceRules.plist",
				"Info.plist",
			},
		}
		found := true
		for _, f := range sd.files {
			path := filepath.Join(sd.name, f)
			if _, err := os.Stat(path); err != nil {
				found = false
				break
			}
		}
		if found {
			return sd
		}
		// In a self-hosted iOS build the above files might
		// not exist. Look for system files instead below.
	case "windows":
		return &sysDir{
			os.Getenv("SystemRoot") + "\\system32\\drivers\\etc",
			[]string{
				"networks",
				"protocol",
				"services",
			},
		}
	case "plan9":
		return &sysDir{
			"/lib/ndb",
			[]string{
				"common",
				"local",
			},
		}
	case "wasip1":
		// wasmtime has issues resolving symbolic links that are often present
		// in directories like /etc/group below (e.g. private/etc/group on OSX).
		// For this reason we use files in the Go source tree instead.
		return &sysDir{
			runtime.GOROOT(),
			[]string{
				"go.env",
				"LICENSE",
				"CONTRIBUTING.md",
			},
		}
	}
	return &sysDir{
		"/etc",
		[]string{
			"group",
			"hosts",
			"passwd",
		},
	}
}()

var sfdir = sysdir.name
var sfname = sysdir.files[0]
