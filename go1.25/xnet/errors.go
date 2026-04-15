package xnet

import (
	"net"
	"os"
	"syscall"
)

func IsAddrInUseError(err error) bool {
	if err == nil {
		return false
	}
	if err, ok := err.(*net.OpError); ok {
		if err, ok := err.Err.(*os.SyscallError); ok {
			return err.Err == syscall.EADDRINUSE
		}
	}
	return false
}
