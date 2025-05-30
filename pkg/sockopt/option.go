/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/

package sockopt

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// SocketOption represents a socket option with its name, level, option number, and description
type SocketOption struct {
	Name        string
	Level       int
	Option      int
	Description string
	MinValue    int // Minimum allowed value
	MaxValue    int // Maximum allowed value, 0 means no upper limit
	Default     int // Default value
}

// GetOptionByName returns a socket option by its name
func GetOptionByName(name string) (*SocketOption, error) {
	opt, ok := socketOptions[name]
	if !ok {
		return nil, fmt.Errorf("unknown socket option: %s", name)
	}
	return &opt, nil
}

// socketOptions is a map of supported socket options
var socketOptions = map[string]SocketOption{
	// TCP level options
	"TCP_NODELAY": {
		Name:        "TCP_NODELAY",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_NODELAY,
		Description: "Disable Nagle's algorithm (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     0,
	},
	"TCP_MAXSEG": {
		Name:        "TCP_MAXSEG",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_MAXSEG,
		Description: "Maximum segment size (bytes)",
		MinValue:    536,
		MaxValue:    65535,
		Default:     1460,
	},
	"TCP_CORK": {
		Name:        "TCP_CORK",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_CORK,
		Description: "Don't send partial frames (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     0,
	},
	"TCP_KEEPIDLE": {
		Name:        "TCP_KEEPIDLE",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_KEEPIDLE,
		Description: "Time (in seconds) before sending keepalive probes",
		MinValue:    1,
		MaxValue:    32767,
		Default:     7200,
	},
	"TCP_KEEPINTVL": {
		Name:        "TCP_KEEPINTVL",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_KEEPINTVL,
		Description: "Time (in seconds) between keepalive probes",
		MinValue:    1,
		MaxValue:    32767,
		Default:     75,
	},
	"TCP_KEEPCNT": {
		Name:        "TCP_KEEPCNT",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_KEEPCNT,
		Description: "Number of keepalive probes before connection drop",
		MinValue:    1,
		MaxValue:    127,
		Default:     9,
	},
	"TCP_SYNCNT": {
		Name:        "TCP_SYNCNT",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_SYNCNT,
		Description: "Number of SYN retransmits before giving up",
		MinValue:    1,
		MaxValue:    255,
		Default:     6,
	},
	"TCP_LINGER2": {
		Name:        "TCP_LINGER2",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_LINGER2,
		Description: "Lifetime of orphaned FIN-WAIT-2 state (seconds)",
		MinValue:    -1,
		MaxValue:    32767,
		Default:     60,
	},
	"TCP_DEFER_ACCEPT": {
		Name:        "TCP_DEFER_ACCEPT",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_DEFER_ACCEPT,
		Description: "Delay accept() until data arrives (seconds)",
		MinValue:    0,
		MaxValue:    32767,
		Default:     0,
	},
	"TCP_WINDOW_CLAMP": {
		Name:        "TCP_WINDOW_CLAMP",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_WINDOW_CLAMP,
		Description: "Bound advertised window (bytes)",
		MinValue:    0,
		MaxValue:    0, // No upper limit
		Default:     0,
	},
	"TCP_QUICKACK": {
		Name:        "TCP_QUICKACK",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_QUICKACK,
		Description: "Enable quickack mode (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     1,
	},
	"TCP_FASTOPEN": {
		Name:        "TCP_FASTOPEN",
		Level:       unix.IPPROTO_TCP,
		Option:      unix.TCP_FASTOPEN,
		Description: "Enable TCP Fast Open (0: disabled, others: queue length)",
		MinValue:    0,
		MaxValue:    0, // No upper limit
		Default:     0,
	},

	// Socket level options
	"SO_KEEPALIVE": {
		Name:        "SO_KEEPALIVE",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_KEEPALIVE,
		Description: "Enable TCP keepalive (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     0,
	},
	"SO_RCVBUF": {
		Name:        "SO_RCVBUF",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_RCVBUF,
		Description: "Socket receive buffer size (bytes)",
		MinValue:    2048,
		MaxValue:    0, // System dependent
		Default:     87380,
	},
	"SO_SNDBUF": {
		Name:        "SO_SNDBUF",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_SNDBUF,
		Description: "Socket send buffer size (bytes)",
		MinValue:    2048,
		MaxValue:    0, // System dependent
		Default:     87380,
	},
	"SO_RCVLOWAT": {
		Name:        "SO_RCVLOWAT",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_RCVLOWAT,
		Description: "Minimum receive buffer space available (bytes)",
		MinValue:    1,
		MaxValue:    0, // System dependent
		Default:     1,
	},
	"SO_SNDLOWAT": {
		Name:        "SO_SNDLOWAT",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_SNDLOWAT,
		Description: "Minimum send buffer space available (bytes)",
		MinValue:    1,
		MaxValue:    0, // System dependent
		Default:     1,
	},
	"SO_REUSEADDR": {
		Name:        "SO_REUSEADDR",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_REUSEADDR,
		Description: "Allow reuse of local addresses (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     0,
	},
	"SO_REUSEPORT": {
		Name:        "SO_REUSEPORT",
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_REUSEPORT,
		Description: "Allow multiple sockets to bind to same address/port (0: disabled, 1: enabled)",
		MinValue:    0,
		MaxValue:    1,
		Default:     0,
	},
}

func (so SocketOption) Set(socketFD int, value int) error {
	err := unix.SetsockoptInt(socketFD, so.Level, so.Option, value)
	if err != nil {
		err = fmt.Errorf("unable to get value of sockopt option %s", so.Name)
	}

	return err

}

func (so SocketOption) Get(socketFD int) (int, error) {
	val, err := unix.GetsockoptInt(socketFD, so.Level, so.Option)
	if err != nil {
		err = fmt.Errorf("unable to get value of sockopt option %s", so.Name)
	}

	return val, err
}

// OptionsList Provides stable order for output of list command
var OptionsList = []string{
	"SO_KEEPALIVE",
	"TCP_KEEPALIVE",
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
	"TCP_FASTOPEN",
	"TCP_TIMESTAMP",
}

var OptionsMap = map[string]SocketOption{
	"SO_KEEPALIVE": {
		Level:       unix.SOL_SOCKET,
		Option:      unix.SO_KEEPALIVE,
		Name:        "SO_KEEPALIVE",
		MinValue:    0,
		MaxValue:    1,
		Description: "Enable or disable TCP keepalive",
		Default:     0,
	},
	"TCP_KEEPIDLE": {
		Name:        "TCP_KEEPIDLE",
		Option:      unix.TCP_KEEPIDLE,
		Level:       unix.IPPROTO_TCP,
		MinValue:    1,
		MaxValue:    32767,
		Description: "Start keepalives after this period",
		Default:     7200,
	},
	"TCP_KEEPINTVL": {
		Name:        "TCP_KEEPINTVL",
		Option:      unix.TCP_KEEPINTVL,
		Level:       unix.IPPROTO_TCP,
		MinValue:    1,
		MaxValue:    32767,
		Description: "Interval between keepalives",
		Default:     75,
	},
	"TCP_KEEPCNT": {
		Name:        "TCP_KEEPCNT",
		Option:      unix.TCP_KEEPCNT,
		Level:       unix.IPPROTO_TCP,
		MinValue:    1,
		MaxValue:    32767,
		Description: "Number of keepalives before death",
		Default:     9,
	},
	"TCP_NODELAY": {
		Name:        "TCP_NODELAY",
		Option:      unix.TCP_NODELAY,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    1,
		Description: "Disable Nagle's algorithm",
		Default:     0,
	},
	"TCP_MAXSEG": {
		Name:        "TCP_MAXSEG",
		Option:      unix.TCP_MAXSEG,
		Level:       unix.IPPROTO_TCP,
		MinValue:    536,
		MaxValue:    65535,
		Description: "Maximum segment size",
		Default:     1460,
	},
	"TCP_CORK": {
		Name:        "TCP_CORK",
		Option:      unix.TCP_CORK,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    1,
		Description: "Control sending of partial frames",
		Default:     0,
	},
	"TCP_SYNCNT": {
		Name:        "TCP_SYNCNT",
		Option:      unix.TCP_SYNCNT,
		Level:       unix.IPPROTO_TCP,
		MinValue:    1,
		MaxValue:    255,
		Description: "Number of SYN retransmits",
		Default:     6,
	},
	"TCP_LINGER2": {
		Name:        "TCP_LINGER2",
		Option:      unix.TCP_LINGER2,
		Level:       unix.IPPROTO_TCP,
		MinValue:    -1,
		MaxValue:    32767,
		Description: "Lifetime of orphaned FIN-WAIT-2 state",
		Default:     60,
	},
	"TCP_DEFER_ACCEPT": {
		Name:        "TCP_DEFER_ACCEPT",
		Option:      unix.TCP_DEFER_ACCEPT,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    32767,
		Description: "Wake up listener only when data arrives",
		Default:     0,
	},
	"TCP_WINDOW_CLAMP": {
		Name:        "TCP_WINDOW_CLAMP",
		Option:      unix.TCP_WINDOW_CLAMP,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Set maximum window size",
		Default:     0,
	},
	"TCP_INFO": {
		Name:        "TCP_INFO",
		Option:      unix.TCP_INFO,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Information about this socket",
	},
	"TCP_QUICKACK": {
		Name:        "TCP_QUICKACK",
		Option:      unix.TCP_QUICKACK,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    1,
		Description: "Enable quick ACK",
		Default:     1,
	},
	"TCP_CONGESTION": {
		Name:        "TCP_CONGESTION",
		Option:      unix.TCP_CONGESTION,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Get/Set congestion control algorithm",
	},
	"TCP_REPAIR": {
		Name:        "TCP_REPAIR",
		Option:      unix.TCP_REPAIR,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    1,
		Description: "TCP repair mode",
	},
	"TCP_REPAIR_QUEUE": {
		Name:        "TCP_REPAIR_QUEUE",
		Option:      unix.TCP_REPAIR_QUEUE,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    3,
		Description: "Repair queue (0: NONE, 1: RECV, 2: SEND)",
	},
	"TCP_QUEUE_SEQ": {
		Name:        "TCP_QUEUE_SEQ",
		Option:      unix.TCP_QUEUE_SEQ,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Set/get queue sequence",
	},
	"TCP_REPAIR_OPTIONS": {
		Name:        "TCP_REPAIR_OPTIONS",
		Option:      unix.TCP_REPAIR_OPTIONS,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Repair options",
	},
	"TCP_FASTOPEN": {
		Name:        "TCP_FASTOPEN",
		Option:      unix.TCP_FASTOPEN,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    0,
		Description: "Enable TCP Fast Open",
		Default:     0,
	},
	"TCP_TIMESTAMP": {
		Name:        "TCP_TIMESTAMP",
		Option:      unix.TCP_TIMESTAMP,
		Level:       unix.IPPROTO_TCP,
		MinValue:    0,
		MaxValue:    1,
		Description: "Enable TCP timestamps",
	},
}
