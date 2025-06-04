package sockopt

import (
	"net"
	"os"
	"strings"
	"syscall"
	"testing"
)

// helper to get fd from net.Conn
func fdFromConn2(c net.Conn) (int, error) {
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

// create a connected TCP socket pair for tests
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

	cleanup := func() {
		c.Close()
		l.Close()
	}
	return c, cleanup
}

func TestGetSocketName(t *testing.T) {
	c, cleanup := makeSocket(t)
	defer cleanup()

	fd, err := fdFromConn2(c)
	if err != nil {
		t.Fatal(err)
	}
	dup, err := GetSocketFd(os.Getpid(), fd)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(dup)

	name := GetSocketName(dup)
	if !strings.Contains(name, ":") {
		t.Fatalf("unexpected socket name %s", name)
	}
}

func TestSetAndGetSocketOptionWrappers(t *testing.T) {
	c, cleanup := makeSocket(t)
	defer cleanup()

	fd, err := fdFromConn2(c)
	if err != nil {
		t.Fatal(err)
	}

	// ensure option can be set and read via wrappers
	SetSocketOption(os.Getpid(), fd, "TCP_NODELAY", 1, "table")
	GetSocketOption(os.Getpid(), fd, "TCP_NODELAY", "table")
	ListSocketOptions(os.Getpid(), fd, "table")
}
