package xembed

import (
	"io/fs"
	"path"
	"strings"
)

type RelativeFs interface {
	Fs
	CWD() string
	// Cd works like ordinary cd (change directory) command,
	// supporting both relative and absolute, and dotted paths.
	// But it does not support special commands `-`,
	// and obviously it does not expand shell shortcuts and variables,
	// such as `~` and `$HOME`.
	// Cd does nothing if argument is empty string.
	Cd(string)
	Clone() RelativeFs
}

func NewRelativeFs(fs Fs) RelativeFs {
	return &relativeFs{
		fs_: fs,
		cwd: "/",
	}
}

type relativeFs struct {
	fs_ Fs
	cwd string
}

func (rfs relativeFs) Open(name string) (fs.File, error) {
	return rfs.fs_.Open(rfs.absPath(name))
}

func (rfs relativeFs) ReadDir(name string) ([]fs.DirEntry, error) {
	return rfs.fs_.ReadDir(rfs.absPath(name))
}

func (rfs relativeFs) ReadFile(name string) ([]byte, error) {
	return rfs.fs_.ReadFile(rfs.absPath(name))
}

func (rfs relativeFs) CWD() string {
	return rfs.cwd
}

func (rfs *relativeFs) Cd(p string) {
	rfs.cwd = cd(rfs.cwd, p)
}

func (rfs *relativeFs) Clone() RelativeFs {
	return &relativeFs{
		fs_: rfs.fs_,
		cwd: rfs.cwd,
	}
}

func (rfs *relativeFs) absPath(name string) string {
	return strings.TrimPrefix(cd(rfs.cwd, name), "/")
}

func cd(cur string, p string) string {
	if len(p) == 0 {
		return cur
	}
	if p[0] == '/' {
		return path.Clean(p)
	}
	return path.Join(cur, p)
}
