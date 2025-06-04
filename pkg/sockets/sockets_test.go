package sockets

import (
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

// helper similar to sockopt tests
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

func TestParseAddressAndState(t *testing.T) {
	if parseHexIP("0100007F") != "127.0.0.1" {
		t.Fatalf("unexpected hex ip")
	}
	addr := parseAddress("0100007F:0016")
	if addr != "127.0.0.1:22" {
		t.Fatalf("got %s", addr)
	}
	if parseState("01") != "ESTABLISHED" {
		t.Fatalf("unexpected state")
	}
}

func TestFindPidFdFromInode(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

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
	defer c.Close()

	fd, err := fdFromConn(c)
	if err != nil {
		t.Fatal(err)
	}

	link, err := os.Readlink("/proc/" + strconv.Itoa(os.Getpid()) + "/fd/" + strconv.Itoa(fd))
	if err != nil {
		t.Skip("cannot readlink")
	}
	inode := strings.TrimSuffix(strings.TrimPrefix(link, "socket:["), "]")
	pid, newfd, err := findPidFdFromInode(inode)
	if err != nil {
		t.Fatal(err)
	}
	if pid != strconv.Itoa(os.Getpid()) || newfd != strconv.Itoa(fd) {
		t.Fatalf("unexpected pid/fd %s %s", pid, newfd)
	}
}

func TestGetConnections(t *testing.T) {
	getConnections()
}
