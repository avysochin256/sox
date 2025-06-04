package sockopt

import (
	"net"
	"os"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// helper to obtain fd from net.Conn
func fdFromConn(c net.Conn) (int, error) {
	sc, ok := c.(syscall.Conn)
	if !ok {
		return 0, os.ErrInvalid
	}
	raw, err := sc.SyscallConn()
	if err != nil {
		return 0, err
	}
	var fd int
	err = raw.Control(func(f uintptr) {
		fd = int(f)
	})
	return fd, err
}

func TestOptionsListInMap(t *testing.T) {
	for _, name := range OptionsList {
		if _, ok := OptionsMap[name]; !ok {
			t.Errorf("option %s missing in OptionsMap", name)
		}
	}
}

func TestSetGetOption(t *testing.T) {
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

	// duplicate fd via pidfd
	newfd, err := GetSocketFd(os.Getpid(), fd)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(newfd)

	opt := OptionsMap["TCP_NODELAY"]

	if err := opt.Set(newfd, 1); err != nil {
		t.Fatal(err)
	}

	val, err := opt.Get(newfd)
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.Fatalf("expected 1 got %d", val)
	}
}
