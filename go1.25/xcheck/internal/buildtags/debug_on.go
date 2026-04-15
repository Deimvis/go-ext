//go:build debug

package buildtags

func OnDebug(fn func()) {
	fn()
}
