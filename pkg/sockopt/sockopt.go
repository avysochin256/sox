// Package sockopt contains helper functions used by the CLI commands for
// listing and updating socket options.
package sockopt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosuri/uitable"
	"golang.org/x/sys/unix"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// OptionRow represents a single socket option value used for output.
type OptionRow struct {
	Name        string `json:"name" yaml:"name"`
	Value       any    `json:"value" yaml:"value"`
	Description string `json:"description" yaml:"description"`
}

// printOutput prints data in the requested format.
func printOutput(data any, headers []string, format string) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	case "yaml":
		b, err := yaml.Marshal(data)
		if err == nil {
			fmt.Print(string(b))
		}
	default:
		table := uitable.New()
		table.MaxColWidth = 50
		hi := make([]interface{}, len(headers))
		for i, h := range headers {
			hi[i] = h
		}
		table.AddRow(hi...)
		switch v := data.(type) {
		case OptionRow:
			table.AddRow(v.Name, v.Value, v.Description)
		case []OptionRow:
			for _, r := range v {
				table.AddRow(r.Name, r.Value, r.Description)
			}
		}
		fmt.Println(table)
	}
}

// GetSocketName returns the IPv4 address and port of a socket file descriptor.
func GetSocketName(socketFd int) string {
	sn, _ := unix.Getsockname(socketFd)
	sai4 := sn.(*unix.SockaddrInet4)
	socketName := fmt.Sprintf("%d.%d.%d.%d:%d", sai4.Addr[0], sai4.Addr[1], sai4.Addr[2], sai4.Addr[3], sai4.Port)

	return socketName
}

// ListSocketOptions prints all supported options for the given pid/fd pair.
func ListSocketOptions(pid, fd int, format string) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	var rows []OptionRow

	var joinedListErr error
	for _, soname := range OptionsList {
		so := OptionsMap[soname]

		val, err := so.Get(socketFd)

		if err != nil {
			err = fmt.Errorf("unable to get sockopt option %s : %w", so.Name, err)
			errors.Join(joinedListErr, err)

			continue
		}

		display := any(val)
		if so.Unsigned {
			display = fmt.Sprintf("%d", uint32(val))
		}
		rows = append(rows, OptionRow{so.Name, display, so.Description})
	}

	printOutput(rows, []string{"OPTION NAME", "VALUE", "DESCRIPTION"}, format)

	if uw, ok := joinedListErr.(interface{ Unwrap() []error }); ok {
		errs := uw.Unwrap()
		for _, err := range errs {
			slog.Error("unable to get value of sockopt option", slog.Any("error", err))
		}
	}

}

// SetSocketOption changes the option value for the socket defined by pid/fd.
func SetSocketOption(pid, fd int, option string, val int, format string) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	var row OptionRow

	so, ok := OptionsMap[option]
	if !ok {
		err = fmt.Errorf("unsupported socket option %s : %w", option, err)
		os.Exit(1)

	}

	err = so.Set(socketFd, val)
	if err != nil {
		err = fmt.Errorf("unable to set sockopt option %s : %w", so.Name, err)
		os.Exit(1)
	}

	val, err = so.Get(socketFd) // TODO: Remove shadowing

	if err != nil {
		err = fmt.Errorf("unable to get socket option %s  after value was set: %w", so.Name, err)

	}

	display := any(val)
	if so.Unsigned {
		display = fmt.Sprintf("%d", uint32(val))
	}
	row = OptionRow{so.Name, display, so.Description}

	printOutput(row, []string{"SOCKET_OPTION", "VALUE", "DESCRIPTION"}, format)

}

// GetSocketOption prints a single socket option value for the socket defined
// by pid/fd.
func GetSocketOption(pid, fd int, option string, format string) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	var row OptionRow

	so, ok := OptionsMap[option]
	if !ok {
		err = fmt.Errorf("unsupported socket option %s : %w", option, err)
		os.Exit(1)

	}

	val, err := so.Get(socketFd)

	if err != nil {
		err = fmt.Errorf("unable to get socket option %s  after value was set: %w", so.Name, err)

	}

	display := any(val)
	if so.Unsigned {
		display = fmt.Sprintf("%d", uint32(val))
	}
	row = OptionRow{so.Name, display, so.Description}

	printOutput(row, []string{"SOCKET_OPTION", "VALUE", "DESCRIPTION"}, format)

}
