//go:build !windows

package syscallw

import "syscall"

func Flock(fd int, how int) (err error) {
	return syscall.Flock(fd, how)
}
