package cmd

import (
	"net"
	"os"
	"strconv"
	"syscall"
	"testing"
)

// helper to get fd from net.Conn
func fdFromConn(c net.Conn) (int, error) {
	sc, ok := c.(syscall.Conn)
	if !ok {
		return 0, os.ErrInvalid
	}
	var fd int
	raw, err := sc.SyscallConn()
	if err != nil {
		return 0, err
	}
	err = raw.Control(func(f uintptr) { fd = int(f) })
	return fd, err
}

func makeSocket(t *testing.T) (net.Conn, func()) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ready := make(chan struct{})
	go func() {
		close(ready)
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		select {}
	}()
	<-ready
	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() { c.Close(); l.Close() }
	return c, cleanup
}

func TestCommands(t *testing.T) {
	c, cleanup := makeSocket(t)
	defer cleanup()
	fd, err := fdFromConn(c)
	if err != nil {
		t.Fatal(err)
	}
	pidStr := strconv.Itoa(os.Getpid())
	fdStr := strconv.Itoa(fd)

	getCmd.Run(getCmd, []string{pidStr, fdStr, "TCP_NODELAY"})
	setCmd.Run(setCmd, []string{pidStr, fdStr, "TCP_NODELAY", "1"})
	listCmd.Run(listCmd, []string{pidStr, fdStr})

	// root command execution
	rootCmd.SetArgs([]string{"get", pidStr, fdStr, "TCP_NODELAY"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestCommandsInvalidArgs(t *testing.T) {
	getCmd.Run(getCmd, []string{"bad", "fd", "TCP_NODELAY"})
	listCmd.Run(listCmd, []string{"bad", "fd"})
}
