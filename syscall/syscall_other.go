//go:build !windows

package syscall

import "syscall"

func Flock(fd int, how int) (err error) {
	return syscall.Flock(fd, how)
}
