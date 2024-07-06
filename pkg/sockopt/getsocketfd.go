package sockopt

import (
	"errors"
	"github.com/oraoto/go-pidfd"
)

var (
	ErrUnableToGetPidFd    = errors.New("unable to get fd of pid")
	ErrUnableToGetSocketFd = errors.New("unable to get fd of pid")
)

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
