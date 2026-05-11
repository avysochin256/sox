// Package sockopt contains helper functions used by the CLI commands for
// listing and updating socket options.
package sockopt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"syscall"

	"github.com/gosuri/uitable"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v3"
)

// OptionRow represents a single socket option value used for output.
type OptionRow struct {
	Name        string `json:"name" yaml:"name"`
	Value       any    `json:"value" yaml:"value"`
	Hint        string `json:"hint,omitempty" yaml:"hint,omitempty"`
	Description string `json:"description" yaml:"description"`
}

// printOutput prints data in the requested format.
func printOutput(data any, headers []string, format string) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			slog.Error("unable to marshal output to JSON", slog.Any("error", err))
			return
		}
		fmt.Println(string(b))
	case "yaml":
		b, err := yaml.Marshal(data)
		if err != nil {
			slog.Error("unable to marshal output to YAML", slog.Any("error", err))
			return
		}
		fmt.Print(string(b))
	default:
		table := uitable.New()
		table.MaxColWidth = 80
		table.Wrap = true
		hi := make([]interface{}, len(headers))
		for i, h := range headers {
			hi[i] = h
		}
		table.AddRow(hi...)
		switch v := data.(type) {
		case OptionRow:
			table.AddRow(v.Name, v.Value, v.Hint, v.Description)
		case []OptionRow:
			for _, r := range v {
				table.AddRow(r.Name, r.Value, r.Hint, r.Description)
			}
		}
		fmt.Println(table)
	}
}

// GetSocketName returns the address:port of the socket bound to socketFd.
// Supports both IPv4 and IPv6 sockets.
func GetSocketName(socketFd int) (string, error) {
	sn, err := unix.Getsockname(socketFd)
	if err != nil {
		return "", fmt.Errorf("getsockname: %w", err)
	}
	var (
		addr []byte
		port int
	)
	switch sa := sn.(type) {
	case *unix.SockaddrInet4:
		addr, port = sa.Addr[:], sa.Port
	case *unix.SockaddrInet6:
		addr, port = sa.Addr[:], sa.Port
	default:
		return "", fmt.Errorf("unsupported sockaddr type %T", sn)
	}
	return net.JoinHostPort(net.IP(addr).String(), strconv.Itoa(port)), nil
}

// ListSocketOptions prints all supported options for the given pid/fd pair.
func ListSocketOptions(pid, fd int, format string) error {
	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		return fmt.Errorf("unable to get sockopt fd: %w", err)
	}
	defer syscall.Close(socketFd)

	rows := make([]OptionRow, 0, len(OptionsList))
	for _, soname := range OptionsList {
		so := OptionsMap[soname]

		val, err := so.Get(socketFd)
		if err != nil {
			// Some options (e.g. TCP_REPAIR_QUEUE, TCP_QUEUE_SEQ,
			// TCP_REPAIR_OPTIONS) are only readable while TCP_REPAIR is
			// enabled, and SO_BINDTOIFINDEX is set-only on kernels < 5.7.
			// The "n/a" cell is the right signal; running `sox get` on the
			// specific option will surface the underlying error.
			rows = append(rows, OptionRow{Name: so.Name, Value: "n/a", Description: so.Description})
			continue
		}

		display := any(val)
		if so.Unsigned {
			display = fmt.Sprintf("%d", uint32(val))
		}
		rows = append(rows, OptionRow{Name: so.Name, Value: display, Hint: so.Hint(val), Description: so.Description})
	}

	printOutput(rows, []string{"OPTION NAME", "VALUE", "HINT", "DESCRIPTION"}, format)
	return nil
}

// SetSocketOption changes the option value for the socket defined by pid/fd.
func SetSocketOption(pid, fd int, option string, val int, format string) error {
	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		return fmt.Errorf("unable to get sockopt fd: %w", err)
	}
	defer syscall.Close(socketFd)

	so, ok := OptionsMap[option]
	if !ok {
		return fmt.Errorf("unsupported socket option %q", option)
	}

	if err := so.Set(socketFd, val); err != nil {
		return fmt.Errorf("unable to set sockopt option %s: %w", so.Name, err)
	}

	// Set succeeded; render the current value. If the post-set read-back
	// fails (e.g. write-only options like TCP_REPAIR_OPTIONS), report
	// "n/a" rather than fabricating a 0 from the unix.GetsockoptInt
	// `(0, err)` convention.
	var (
		display any
		hint    string
	)
	got, getErr := so.Get(socketFd)
	if getErr != nil {
		slog.Warn("set succeeded but post-set read-back failed",
			slog.String("option", so.Name),
			slog.Any("error", getErr))
		display = "n/a"
	} else {
		display = got
		if so.Unsigned {
			display = fmt.Sprintf("%d", uint32(got))
		}
		hint = so.Hint(got)
	}

	row := OptionRow{Name: so.Name, Value: display, Hint: hint, Description: so.Description}
	printOutput(row, []string{"SOCKET_OPTION", "VALUE", "HINT", "DESCRIPTION"}, format)
	return nil
}

// ExplainData carries all metadata about a single socket option for the
// `sox explain` command.
type ExplainData struct {
	Name        string `json:"name" yaml:"name"`
	Level       string `json:"level" yaml:"level"`
	Min         int    `json:"min" yaml:"min"`
	Max         int    `json:"max" yaml:"max"`
	ReadOnly    bool   `json:"read_only" yaml:"read_only"`
	Description string `json:"description" yaml:"description"`
	Details     string `json:"details" yaml:"details"`
}

// ExplainSocketOption prints the full description of a socket option, in the
// requested output format.
func ExplainSocketOption(option, format string) error {
	so, ok := OptionsMap[option]
	if !ok {
		return fmt.Errorf("unsupported socket option %q (run `sox list <pid> <fd>` for the supported set)", option)
	}

	data := ExplainData{
		Name:        so.Name,
		Level:       so.LevelName(),
		Min:         so.MinVal,
		Max:         so.MaxVal,
		ReadOnly:    so.MinVal == so.MaxVal,
		Description: so.Description,
		Details:     so.Details,
	}

	switch format {
	case "json":
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal explain output to JSON: %w", err)
		}
		fmt.Println(string(b))
	case "yaml":
		b, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("unable to marshal explain output to YAML: %w", err)
		}
		fmt.Print(string(b))
	default:
		rangeStr := fmt.Sprintf("[%d, %d]", so.MinVal, so.MaxVal)
		if data.ReadOnly {
			rangeStr = "(read-only / no range validation)"
		}
		fmt.Printf("%s\n", so.Name)
		fmt.Printf("%s\n\n", strings.Repeat("=", len(so.Name)))
		fmt.Printf("Level: %s\n", data.Level)
		fmt.Printf("Range: %s\n\n", rangeStr)
		fmt.Printf("Summary:\n  %s\n\n", so.Description)
		if so.Details != "" {
			fmt.Printf("Details:\n%s\n", indent(so.Details, "  "))
		}
	}
	return nil
}

// indent prefixes every line of s with the given prefix.
func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if l == "" {
			continue
		}
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}

// GetSocketOption prints a single socket option value for the socket defined
// by pid/fd.
func GetSocketOption(pid, fd int, option string, format string) error {
	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		return fmt.Errorf("unable to get sockopt fd: %w", err)
	}
	defer syscall.Close(socketFd)

	so, ok := OptionsMap[option]
	if !ok {
		return fmt.Errorf("unsupported socket option %q", option)
	}

	val, err := so.Get(socketFd)
	if err != nil {
		return fmt.Errorf("unable to get sockopt option %s: %w", so.Name, err)
	}

	display := any(val)
	if so.Unsigned {
		display = fmt.Sprintf("%d", uint32(val))
	}
	row := OptionRow{Name: so.Name, Value: display, Hint: so.Hint(val), Description: so.Description}

	printOutput(row, []string{"SOCKET_OPTION", "VALUE", "HINT", "DESCRIPTION"}, format)
	return nil
}
