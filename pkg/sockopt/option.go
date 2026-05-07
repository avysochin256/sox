// Package sockopt contains low level wrappers around various TCP socket options.
package sockopt

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

// SocketOption describes a single socket option.
// MinVal and MaxVal are used for basic range validation when setting values.
// Description is a one-or-two-line summary used by `sox list`.
// Details is the long-form explanation used by `sox explain`.
type SocketOption struct {
	Name        string
	Option      int
	Level       int
	MinVal      int
	MaxVal      int
	Unsigned    bool
	Description string
	Details     string
}

// LevelName returns the symbolic name of the getsockopt level (SOL_SOCKET,
// IPPROTO_TCP, ...) for display in `sox explain`.
func (so SocketOption) LevelName() string {
	switch so.Level {
	case unix.SOL_SOCKET:
		return "SOL_SOCKET"
	case unix.IPPROTO_TCP:
		return "IPPROTO_TCP"
	default:
		return fmt.Sprintf("level=%d", so.Level)
	}
}

// Set changes the value of the socket option for the given socket file descriptor.
func (so SocketOption) Set(socketFD int, value int) error {
	if so.MaxVal != so.MinVal && (value < so.MinVal || value > so.MaxVal) {
		return fmt.Errorf("value %d out of range [%d,%d] for %s", value, so.MinVal, so.MaxVal, so.Name)
	}

	err := unix.SetsockoptInt(socketFD, so.Level, so.Option, value)
	if err != nil {
		err = fmt.Errorf("unable to set sockopt option %s: %w", so.Name, err)
	}

	return err

}

// Hint returns a short human-friendly annotation for a value that cannot be
// inferred from the integer alone — e.g. the interface name behind a
// SO_BINDTOIFINDEX ifindex. Empty string means no hint.
func (so SocketOption) Hint(val int) string {
	switch so.Name {
	case "SO_BINDTOIFINDEX":
		if val <= 0 {
			return ""
		}
		if iface, err := net.InterfaceByIndex(val); err == nil {
			return iface.Name
		}
	}
	return ""
}

// Get returns the current value of the socket option for the given socket file descriptor.
func (so SocketOption) Get(socketFD int) (int, error) {
	val, err := unix.GetsockoptInt(socketFD, so.Level, so.Option)
	if err != nil {
		err = fmt.Errorf("unable to get value of sockopt option %s: %w", so.Name, err)
	}

	return val, err
}

// OptionsList provides a stable order for the list command output.
var OptionsList = []string{
	"SO_KEEPALIVE",
	"SO_BINDTOIFINDEX",
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
		Description: "Enable (1) / disable (0) periodic keepalive probes for dead-peer detection.",
		Details: `When enabled, the kernel sends periodic ACK probes on a connection that has been idle, and tears the connection down if the peer fails to respond. Useful for detecting half-open connections behind NAT or after a peer crash.

The probe schedule is controlled by three companion options:
  TCP_KEEPIDLE  — idle time before the first probe.
  TCP_KEEPINTVL — interval between probes.
  TCP_KEEPCNT   — number of unacknowledged probes before disconnect.

Setting SO_KEEPALIVE alone is not enough; the timers default to 2 hours of idle time, which is rarely what you want.`,
	},
	"SO_BINDTOIFINDEX": {
		Name:        "SO_BINDTOIFINDEX",
		Option:      unix.SO_BINDTOIFINDEX,
		Level:       unix.SOL_SOCKET,
		MinVal:      0,
		MaxVal:      1<<31 - 1,
		Description: "Bind socket to a network interface by ifindex; 0 unbinds. Needs CAP_NET_RAW.",
		Details: `Restricts a socket to a specific network interface, identified by its kernel ifindex. Look up an interface's ifindex with "ip -br link" or "ip address show" (the integer next to the interface name).

Compared to SO_BINDTODEVICE, this variant takes a numeric index rather than a device name, avoiding a name-to-index lookup on every send and surviving renames.

Common use cases:
  - Source-binding outgoing connections to a specific NIC or VRF.
  - Selecting an interface inside a different network namespace.
  - Setting up split-horizon proxying (e.g. nginx accepts on one
    interface, proxy_pass to upstreams via another).

Set 0 to remove an existing binding. Requires CAP_NET_RAW. On kernels older than 5.7, getsockopt for this option returns ENOPROTOOPT — sox displays "n/a" in that case, but set still works.`,
	},
	"TCP_KEEPIDLE": {
		Name:        "TCP_KEEPIDLE",
		Option:      unix.TCP_KEEPIDLE,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Idle seconds before keepalive probes start. Requires SO_KEEPALIVE.",
		Details: `Number of seconds the connection must be idle (no traffic in either direction) before the kernel begins sending keepalive probes. The system-wide default is 7200 seconds (2 hours), set via /proc/sys/net/ipv4/tcp_keepalive_time.

Has no effect unless SO_KEEPALIVE is also enabled on the socket. Pair with TCP_KEEPINTVL and TCP_KEEPCNT for a complete probing policy.`,
	},
	"TCP_KEEPINTVL": {
		Name:        "TCP_KEEPINTVL",
		Option:      unix.TCP_KEEPINTVL,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Seconds between successive keepalive probes.",
		Details: `Interval in seconds between successive keepalive probes after the first one has been sent. System-wide default is 75 seconds (/proc/sys/net/ipv4/tcp_keepalive_intvl).

Total dead-detection time is roughly TCP_KEEPIDLE + TCP_KEEPCNT * TCP_KEEPINTVL.`,
	},
	"TCP_KEEPCNT": {
		Name:        "TCP_KEEPCNT",
		Option:      unix.TCP_KEEPCNT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      32767,
		Description: "Max unacked keepalive probes before the connection is dropped.",
		Details: `Number of keepalive probes that may go unanswered before the kernel concludes the peer is dead and closes the socket with ETIMEDOUT. System-wide default is 9 (/proc/sys/net/ipv4/tcp_keepalive_probes).`,
	},
	"TCP_USER_TIMEOUT": {
		Name:        "TCP_USER_TIMEOUT",
		Option:      unix.TCP_USER_TIMEOUT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      0xFFFFFFFF,
		Description: "Max ms data may stay unacknowledged before TCP closes the connection.",
		Details: `Independent of the keepalive timers, TCP_USER_TIMEOUT places an upper bound on how long data sitting in the retransmission queue may go unacknowledged. Once the timeout fires the connection is aborted with ETIMEDOUT, regardless of the standard exponential backoff.

Useful when you need bounded failure detection time on an active connection — much faster than waiting for the keepalive schedule. Set to 0 to use the kernel default (effectively unlimited).`,
	},
	"TCP_NODELAY": {
		Name:        "TCP_NODELAY",
		Option:      unix.TCP_NODELAY,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Disable Nagle's algorithm; small writes are sent immediately.",
		Details: `Nagle's algorithm holds small writes in the kernel until the previous segment is acknowledged, in order to reduce header overhead on chatty workloads. Setting TCP_NODELAY=1 turns this off, sending data immediately at the cost of more (and smaller) packets.

Almost always wanted for interactive or RPC-style protocols where round-trip latency matters more than per-byte efficiency. Mutually exclusive in spirit with TCP_CORK.`,
	},
	"TCP_MAXSEG": {
		Name:        "TCP_MAXSEG",
		Option:      unix.TCP_MAXSEG,
		Level:       unix.IPPROTO_TCP,
		MinVal:      536,
		MaxVal:      65535,
		Description: "Maximum outgoing TCP segment size in bytes; capped by path MTU.",
		Details: `Caps the size of TCP segments this socket will emit. The kernel may use a smaller value if path MTU discovery or peer MSS announcement requires it. The minimum 536 corresponds to the IPv4 minimum MSS.

Lowering MSS can be useful when working over tunnels with smaller MTUs (IPIP, GRE, WireGuard) but is rarely needed; the kernel discovers a sane value automatically.`,
	},
	"TCP_CORK": {
		Name:        "TCP_CORK",
		Option:      unix.TCP_CORK,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Coalesce writes into full frames until cleared.",
		Details: `Acts as the conceptual opposite of TCP_NODELAY: while corked, the kernel holds back partial segments and only emits packets that fill an MSS, capping the wait at ~200 ms. Clearing the option flushes any pending data immediately.

Typical pattern: cork before a sendfile()/write batch, uncork to flush. Useful for serving HTTP responses where headers and body are written separately but should ideally land in the same packet.`,
	},
	"TCP_SYNCNT": {
		Name:        "TCP_SYNCNT",
		Option:      unix.TCP_SYNCNT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      1,
		MaxVal:      255,
		Description: "Number of SYN retransmits before connect() aborts (default 6).",
		Details: `Caps how many times the kernel will retransmit the initial SYN on connect() before giving up with ETIMEDOUT. Lowering this is the fastest way to shorten connect() timeouts on unreachable hosts; the system-wide default lives in /proc/sys/net/ipv4/tcp_syn_retries.`,
	},
	"TCP_LINGER2": {
		Name:        "TCP_LINGER2",
		Option:      unix.TCP_LINGER2,
		Level:       unix.IPPROTO_TCP,
		MinVal:      -1,
		MaxVal:      32767,
		Description: "Seconds an orphaned FIN_WAIT2 socket lingers; -1 disables the timer.",
		Details: `Per-socket override for net.ipv4.tcp_fin_timeout. Once the local side has closed and the socket is orphaned (no longer referenced by a userspace fd), this controls how long the kernel keeps the connection in FIN_WAIT2 waiting for the peer's FIN.

Set to -1 to disable the timer (matches RFC 793 behavior — the socket stays around until ACKed). The default of 60 seconds is usually fine.`,
	},
	"TCP_DEFER_ACCEPT": {
		Name:        "TCP_DEFER_ACCEPT",
		Option:      unix.TCP_DEFER_ACCEPT,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      32767,
		Description: "Listen-only: defer accept() until data arrives, up to N seconds.",
		Details: `On a listening socket, postpones the return of accept() until the first byte of data is available on the new connection. Reduces wakeups for protocols where the client always speaks first (HTTP, TLS), avoiding a context switch for connections that arrive but never send.

The value is the maximum number of seconds the kernel will wait before completing the handshake regardless of data arrival. Setting 0 disables the deferral.`,
	},
	"TCP_WINDOW_CLAMP": {
		Name:        "TCP_WINDOW_CLAMP",
		Option:      unix.TCP_WINDOW_CLAMP,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1073725440,
		Description: "Upper bound (bytes) on receive window; 0 = kernel autotune.",
		Details: `Caps the maximum receive window the kernel will advertise to the peer for this socket. Useful for capping per-connection memory consumption or for testing how applications behave on bandwidth-limited links.

Setting 0 (or a value below the minimum) disables clamping and re-enables autotuning via /proc/sys/net/ipv4/tcp_rmem.`,
	},
	"TCP_INFO": {
		Name:        "TCP_INFO",
		Option:      unix.TCP_INFO,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Read-only struct tcp_info; sox shows the first byte (TCP state).",
		Details: `Returns a struct tcp_info describing the live state of the connection: TCP state, RTT, RTT variance, retransmits, congestion window, packet counters, and so on. The full structure is defined in <linux/tcp.h>.

sox uses GetsockoptInt and only surfaces the first int, which corresponds to tcpi_state — the TCP connection state (1=ESTABLISHED, 10=LISTEN, etc.). Read-only; trying to set it returns EPERM.`,
	},
	"TCP_QUICKACK": {
		Name:        "TCP_QUICKACK",
		Option:      unix.TCP_QUICKACK,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Force immediate ACKs for the next few packets; not sticky.",
		Details: `Forces the kernel out of delayed-ACK mode briefly: subsequent ACKs are sent right away rather than being coalesced. The flag is not sticky — the kernel may turn it off again as soon as it decides delayed ACKs are appropriate, so applications that want quick ACKs persistently must re-enable it after each receive.

Useful for cutting latency in request/response protocols where the peer is waiting on the ACK to advance.`,
	},
	"TCP_CONGESTION": {
		Name:        "TCP_CONGESTION",
		Option:      unix.TCP_CONGESTION,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Per-socket congestion control algorithm (cubic, bbr, reno, ...).",
		Details: `Selects the congestion control algorithm for this socket independently of the system-wide default (net.ipv4.tcp_congestion_control). Algorithms must be loaded as kernel modules; check /proc/sys/net/ipv4/tcp_available_congestion_control for the active list.

NOTE: The kernel's getsockopt for this option returns the algorithm name as a string. sox calls GetsockoptInt and shows the leading bytes interpreted as an int, so the displayed value is not directly meaningful — treat it as "the algorithm is set" rather than a comparable integer.`,
	},
	"TCP_REPAIR": {
		Name:        "TCP_REPAIR",
		Option:      unix.TCP_REPAIR,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Toggle TCP repair mode for CRIU-style migration. Needs CAP_NET_ADMIN.",
		Details: `Puts the socket into a special "repair" mode that lets userspace inspect and inject internal TCP state — sequence numbers, queue contents, options — without affecting the wire. Used by CRIU and similar tools to migrate live connections across hosts.

Companion options that only work while repair mode is active:
  TCP_REPAIR_QUEUE   — select RECV or SEND queue.
  TCP_QUEUE_SEQ      — read/write the head sequence number.
  TCP_REPAIR_OPTIONS — inject saved TCP options.

Requires CAP_NET_ADMIN. Should be turned off again before normal operation resumes.`,
	},
	"TCP_REPAIR_QUEUE": {
		Name:        "TCP_REPAIR_QUEUE",
		Option:      unix.TCP_REPAIR_QUEUE,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      3,
		Description: "Active repair queue: 0=NONE, 1=RECV, 2=SEND. Requires TCP_REPAIR=1.",
		Details: `While the socket is in repair mode (TCP_REPAIR=1), this picks which side of the connection the next peek/inject operation acts on:
  0  TCP_NO_QUEUE   — no queue selected (default)
  1  TCP_RECV_QUEUE — receive queue (incoming, not yet read)
  2  TCP_SEND_QUEUE — send queue (outgoing, not yet acked)

Outside of repair mode, getsockopt returns EINVAL — sox shows "n/a" then.`,
	},
	"TCP_QUEUE_SEQ": {
		Name:        "TCP_QUEUE_SEQ",
		Option:      unix.TCP_QUEUE_SEQ,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Head sequence number of the selected repair queue. Requires TCP_REPAIR=1.",
		Details: `Reads or writes the TCP sequence number at the head of the currently selected repair queue. Used during checkpoint/restore to capture and replay queue state across a migration.

Outside of repair mode (or with no queue selected), getsockopt returns EINVAL — sox shows "n/a" then.`,
	},
	"TCP_REPAIR_OPTIONS": {
		Name:        "TCP_REPAIR_OPTIONS",
		Option:      unix.TCP_REPAIR_OPTIONS,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Description: "Write-only: inject saved TCP options during socket restore.",
		Details: `Used during restore to push previously captured per-connection TCP options back into the kernel: timestamps (TSval/TSecr), SACK enable, MSS, window scaling. The format is an array of struct tcp_repair_opt (see <linux/tcp.h>).

Write-only: the kernel has no getsockopt handler, so reads always return ENOPROTOOPT. sox always shows "n/a" for this option.`,
	},
	"TCP_FASTOPEN": {
		Name:        "TCP_FASTOPEN",
		Option:      unix.TCP_FASTOPEN,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      1,
		Description: "Enable TCP Fast Open (0-RTT data on SYN).",
		Details: `Lets a client send data in the SYN packet of a reconnection, eliminating the first round trip when a valid TFO cookie is cached.

On a listening socket, the value is the maximum length of the pending TFO request queue. On a connecting socket, it controls cookie acceptance behavior. System-wide policy is governed by /proc/sys/net/ipv4/tcp_fastopen.

Both peers and the network path must support TFO for it to take effect; middleboxes occasionally strip the option.`,
	},
	"TCP_TIMESTAMP": {
		Name:        "TCP_TIMESTAMP",
		Option:      unix.TCP_TIMESTAMP,
		Level:       unix.IPPROTO_TCP,
		MinVal:      0,
		MaxVal:      0,
		Unsigned:    true,
		Description: "Initial TSval; can be set during repair to align timestamps.",
		Details: `Reads or sets the initial value of the TCP timestamp (TSval) used in the TCP timestamps option (RFC 7323). Mostly useful as a diagnostic counter, or during checkpoint/restore where timestamps must be re-aligned so the restored peer doesn't see them go backwards.

The kernel returns this as an unsigned 32-bit value, so sox displays it as uint32 to avoid showing negative numbers.`,
	},
}
