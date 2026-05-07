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

	if err := getCmd.RunE(getCmd, []string{pidStr, fdStr, "TCP_NODELAY"}); err != nil {
		t.Fatal(err)
	}
	if err := setCmd.RunE(setCmd, []string{pidStr, fdStr, "TCP_NODELAY", "1"}); err != nil {
		t.Fatal(err)
	}
	if err := listCmd.RunE(listCmd, []string{pidStr, fdStr}); err != nil {
		t.Fatal(err)
	}

	// root command execution
	rootCmd.SetArgs([]string{"get", pidStr, fdStr, "TCP_NODELAY"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestCommandsInvalidArgs(t *testing.T) {
	if err := getCmd.RunE(getCmd, []string{"bad", "fd", "TCP_NODELAY"}); err == nil {
		t.Fatal("expected error for non-numeric pid")
	}
	if err := listCmd.RunE(listCmd, []string{"bad", "fd"}); err == nil {
		t.Fatal("expected error for non-numeric pid")
	}
	if err := getCmd.Args(getCmd, []string{"1"}); err == nil {
		t.Fatal("expected ExactArgs(3) to reject 1 arg")
	}
	if err := setCmd.Args(setCmd, []string{"1", "2", "TCP_NODELAY"}); err == nil {
		t.Fatal("expected ExactArgs(4) to reject 3 args")
	}
	if err := listCmd.Args(listCmd, []string{"1"}); err == nil {
		t.Fatal("expected ExactArgs(2) to reject 1 arg")
	}
}

func TestExplainCommand(t *testing.T) {
	if err := explainCmd.RunE(explainCmd, []string{"TCP_NODELAY"}); err != nil {
		t.Fatalf("expected explain on known option to succeed: %v", err)
	}
	if err := explainCmd.RunE(explainCmd, []string{"BOGUS"}); err == nil {
		t.Fatal("expected explain on unknown option to error")
	}
	if err := explainCmd.Args(explainCmd, []string{}); err == nil {
		t.Fatal("expected ExactArgs(1) to reject 0 args")
	}
}
