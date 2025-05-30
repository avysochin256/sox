/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/

package sockopt

import (
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/sys/unix"
)

// Error types for socket operations
var (
	ErrInvalidPID      = fmt.Errorf("invalid process ID")
	ErrInvalidFD       = fmt.Errorf("invalid file descriptor")
	ErrInvalidOption   = fmt.Errorf("invalid socket option")
	ErrInvalidValue    = fmt.Errorf("invalid option value")
	ErrProcessAccess   = fmt.Errorf("cannot access process")
	ErrSocketOperation = fmt.Errorf("socket operation failed")
)

func GetSocketName(socketFd int) string {
	sn, _ := unix.Getsockname(socketFd)
	sai4 := sn.(*unix.SockaddrInet4)
	socketName := fmt.Sprintf("%d.%d.%d.%d:%d", sai4.Addr[0], sai4.Addr[1], sai4.Addr[2], sai4.Addr[3], sai4.Port)

	return socketName
}

// validateSocketOption checks if the option exists and the value is within allowed range
func validateSocketOption(optionName string, value int) error {
	opt, err := GetOptionByName(optionName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidOption, err)
	}

	if opt.MaxValue != 0 && value > opt.MaxValue {
		return fmt.Errorf("%w: value %d exceeds maximum %d for option %s",
			ErrInvalidValue, value, opt.MaxValue, optionName)
	}

	if value < opt.MinValue {
		return fmt.Errorf("%w: value %d is below minimum %d for option %s",
			ErrInvalidValue, value, opt.MinValue, optionName)
	}

	return nil
}

// validateProcess checks if the process exists and is accessible
func validateProcess(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("%w: PID must be positive", ErrInvalidPID)
	}

	// Check if process exists
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("%w: process %d not found", ErrInvalidPID, pid)
	}

	// Try to signal the process to check if we have access
	err = proc.Signal(unix.Signal(0))
	if err != nil {
		return fmt.Errorf("%w: cannot access process %d", ErrProcessAccess, pid)
	}

	return nil
}

// validateFD checks if the file descriptor is valid
func validateFD(fd int) error {
	if fd < 0 {
		return fmt.Errorf("%w: FD must be non-negative", ErrInvalidFD)
	}
	return nil
}

// GetSocketOption retrieves the value of a socket option
func GetSocketOption(pid int, fd int, optionName string) error {
	if err := validateProcess(pid); err != nil {
		return err
	}

	if err := validateFD(fd); err != nil {
		return err
	}

	opt, err := GetOptionByName(optionName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidOption, err)
	}

	value, err := unix.GetsockoptInt(fd, opt.Level, opt.Option)
	if err != nil {
		return fmt.Errorf("%w: failed to get option %s: %v", ErrSocketOperation, optionName, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SOCKET_OPTION\tVALUE\tDESCRIPTION")
	fmt.Fprintf(w, "%s\t%d\t%s\n", opt.Name, value, opt.Description)
	return w.Flush()
}

// SetSocketOption sets the value of a socket option
func SetSocketOption(pid int, fd int, optionName string, value int) error {
	if err := validateProcess(pid); err != nil {
		return err
	}

	if err := validateFD(fd); err != nil {
		return err
	}

	if err := validateSocketOption(optionName, value); err != nil {
		return err
	}

	opt, _ := GetOptionByName(optionName) // Error already checked in validateSocketOption

	err := unix.SetsockoptInt(fd, opt.Level, opt.Option, value)
	if err != nil {
		return fmt.Errorf("%w: failed to set option %s to %d: %v",
			ErrSocketOperation, optionName, value, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SOCKET_OPTION\tVALUE\tDESCRIPTION")
	fmt.Fprintf(w, "%s\t%d\t%s\n", opt.Name, value, opt.Description)
	return w.Flush()
}

// ListSocketOptions lists all available socket options and their current values
func ListSocketOptions(pid int, fd int) error {
	if err := validateProcess(pid); err != nil {
		return err
	}

	if err := validateFD(fd); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "OPTION NAME\tVALUE\tDESCRIPTION")

	for _, opt := range socketOptions {
		value, err := unix.GetsockoptInt(fd, opt.Level, opt.Option)
		if err != nil {
			// Skip options that cannot be read
			continue
		}
		fmt.Fprintf(w, "%s\t%d\t%s\n", opt.Name, value, opt.Description)
	}

	return w.Flush()
}
