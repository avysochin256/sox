// Package sockopt provides helpers for reading and modifying socket options.
package sockopt

import (
	"errors"
	"github.com/oraoto/go-pidfd"
)

var (
	// ErrUnableToGetPidFd is returned when pidfd.Open fails.
	ErrUnableToGetPidFd = errors.New("unable to get fd of pid")
	// ErrUnableToGetSocketFd is returned when the file descriptor cannot be obtained.
	ErrUnableToGetSocketFd = errors.New("unable to get fd of pid")
)

// GetSocketFd returns a duplicate of file descriptor fd from the given process.
// It utilises the pidfd mechanism and thus requires Linux 5.6+.
func GetSocketFd(pid, fd int) (int, error) {
	pidFD, err := pidfd.Open(pid, 0)
	if err != nil {
		return 0, ErrUnableToGetPidFd
	}

	socketFD, err := pidFD.GetFd(fd, 0)
	if err != nil {
		return 0, ErrUnableToGetSocketFd
	}

	return socketFD, nil
}
