// Package sockopt contains low level wrappers around various TCP socket options.
package sockopt

import (
	"fmt"
	"golang.org/x/sys/unix"
)

// SocketOption describes a single socket option.
// MinVal and MaxVal are used for basic range validation when setting values.
type SocketOption struct {
	Name        string
	Option      int
	Level       int
	MinVal      int
	MaxVal      int
	Description string
}

// Set changes the value of the socket option for the given socket file descriptor.
func (so SocketOption) Set(socketFD int, value int) error {
	err := unix.SetsockoptInt(socketFD, so.Level, so.Option, value)
	if err != nil {
		err = fmt.Errorf("unable to get value of sockopt option %s", so.Name)
	}

	return err

}

// Get returns the current value of the socket option for the given socket file descriptor.
func (so SocketOption) Get(socketFD int) (int, error) {
	val, err := unix.GetsockoptInt(socketFD, so.Level, so.Option)
	if err != nil {
		err = fmt.Errorf("unable to get value of sockopt option %s", so.Name)
	}

	return val, err
}

// OptionsList provides a stable order for the list command output.
var OptionsList = []string{
	"SO_KEEPALIVE",
	"TCP_KEEPIDLE",
	"TCP_KEEPINTVL",
	"TCP_KEEPCNT",
	"TCP_USER_TIMEOUT",
	"TCP_NODELAY",
	"TCP_MAXSEG",
	"TCP_CORK",
	"TCP_SYNCNT",
	"TCP_LINGER2",
	"TCP_DEFER_ACCEPT",
	"TCP_WINDOW_CLAMP",
	"TCP_INFO",
	"TCP_QUICKACK",
	"TCP_CONGESTION",
	"TCP_REPAIR",
	"TCP_REPAIR_QUEUE",
	"TCP_QUEUE_SEQ",
	"TCP_REPAIR_OPTIONS",
	"TCP_FASTOPEN",
	"TCP_TIMESTAMP",
}

// OptionsMap maps the option name to its description and numeric identifiers.
var OptionsMap = map[string]SocketOption{
	"SO_KEEPALIVE": {
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_KEEPALIVE,
		Name:        "SO_KEEPALIVE",
		MinVal:      0,
		MaxVal:      1,
		Description: "Enable or disable TCP keepalive",
	},
	"TCP_KEEPIDLE": {
		Name:        "TCP_KEEPIDLE",
		Option:      unix.TCP_KEEPIDLE,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Start keepalives after this period",
	},
	"TCP_KEEPINTVL": {
		Name:        "TCP_KEEPINTVL",
		Option:      unix.TCP_KEEPINTVL,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Interval between keepalives",
	},
	"TCP_KEEPCNT": {
		Name:        "TCP_KEEPCNT",
		Option:      unix.TCP_KEEPCNT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Number of keepalives before death",
	},
	"TCP_USER_TIMEOUT": {
		Name:        "TCP_USER_TIMEOUT",
		Option:      unix.TCP_USER_TIMEOUT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      0xFFFFFFFF,
		Description: "Time to wait for peer response",
	},
	"TCP_NODELAY": {
		Name:        "TCP_NODELAY",
		Option:      unix.TCP_NODELAY,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Disable Nagle's algorithm",
	},
	"TCP_MAXSEG": {
		Name:        "TCP_MAXSEG",
		Option:      unix.TCP_MAXSEG,
		Level:       unix.IPPROTO_TCP,
		MinVal:      536,
		MaxVal:      65535,
		Description: "Maximum segment size",
	},
	"TCP_CORK": {
		Name:        "TCP_CORK",
		Option:      unix.TCP_CORK,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Control sending of partial frames",
	},
	"TCP_SYNCNT": {
		Name:        "TCP_SYNCNT",
		Option:      unix.TCP_SYNCNT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      255,
		Description: "Number of SYN retransmits",
	},
	"TCP_LINGER2": {
		Name:        "TCP_LINGER2",
		Option:      unix.TCP_LINGER2,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      32767,
		Description: "Lifetime of orphaned FIN-WAIT-2 state",
	},
	"TCP_DEFER_ACCEPT": {
		Name:        "TCP_DEFER_ACCEPT",
		Option:      unix.TCP_DEFER_ACCEPT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      32767,
		Description: "Wake up listener only when data arrives",
	},
	"TCP_WINDOW_CLAMP": {
		Name:        "TCP_WINDOW_CLAMP",
		Option:      unix.TCP_WINDOW_CLAMP,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      1073725440,
		Description: "Set maximum window size",
	},
	"TCP_INFO": {
		Name:        "TCP_INFO",
		Option:      unix.TCP_INFO,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Information about this socket",
	},
	"TCP_QUICKACK": {
		Name:        "TCP_QUICKACK",
		Option:      unix.TCP_QUICKACK,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Enable quick ACK",
	},
	"TCP_CONGESTION": {
		Name:        "TCP_CONGESTION",
		Option:      unix.TCP_CONGESTION,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Get/Set congestion control algorithm",
	},
	"TCP_REPAIR": {
		Name:        "TCP_REPAIR",
		Option:      unix.TCP_REPAIR,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "TCP repair mode",
	},
	"TCP_REPAIR_QUEUE": {
		Name:        "TCP_REPAIR_QUEUE",
		Option:      unix.TCP_REPAIR_QUEUE,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      3,
		Description: "Repair queue (0: NONE, 1: RECV, 2: SEND)",
	},
	"TCP_QUEUE_SEQ": {
		Name:        "TCP_QUEUE_SEQ",
		Option:      unix.TCP_QUEUE_SEQ,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Set/get queue sequence",
	},
	"TCP_REPAIR_OPTIONS": {
		Name:        "TCP_REPAIR_OPTIONS",
		Option:      unix.TCP_REPAIR_OPTIONS,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Repair options",
	},
	"TCP_FASTOPEN": {
		Name:        "TCP_FASTOPEN",
		Option:      unix.TCP_FASTOPEN,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Enable TCP Fast Open",
	},
	"TCP_TIMESTAMP": {
		Name:        "TCP_TIMESTAMP",
		Option:      unix.TCP_TIMESTAMP,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Enable TCP timestamps",
	},
}
