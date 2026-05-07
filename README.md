# Sox

Sox (Socket Options eXplorer) is a small command line tool that allows you to
inspect and modify TCP socket options of any running process. It can be helpful
for debugging or tuning network applications without restarting them.

## Requirements
- Linux kernel 5.6 or newer (for the pidfd API)
- `CAP_SYS_PTRACE` capability to operate on foreign processes

## Installation
```bash
go install github.com/avysochin256/sox@latest
```

## Quickstart

Follow the steps below to locate a socket and inspect or modify its options.

### 1. Get PID and file descriptor via `ss`
Run `ss -ntpa` and look at the `users` column:

```bash
sudo ss -ntpa

State  Recv-Q Send-Q Local Address:Port Peer Address:Port Process
LISTEN 0      128    0.0.0.0:22         0.0.0.0:*       users:(("sshd",pid=1062,fd=3))
```

Here the PID is `1062` and the file descriptor is `3`.

### 2. List all socket options
```bash
sudo sox list 1062 3

OPTION NAME         VALUE       HINT    DESCRIPTION
SO_KEEPALIVE        1                   Enable (1) / disable (0) periodic keepalive probes for dead-peer detection.
SO_BINDTOIFINDEX    0                   Bind socket to a network interface by ifindex; 0 unbinds. Needs CAP_NET_RAW.
TCP_KEEPIDLE        7200                Idle seconds before keepalive probes start. Requires SO_KEEPALIVE.
TCP_KEEPINTVL       75                  Seconds between successive keepalive probes.
TCP_KEEPCNT         9                   Max unacked keepalive probes before the connection is dropped.
TCP_NODELAY         0                   Disable Nagle's algorithm; small writes are sent immediately.
...
TCP_REPAIR_QUEUE    n/a                 Active repair queue: 0=NONE, 1=RECV, 2=SEND. Requires TCP_REPAIR=1.
```

`n/a` means the option is supported but the kernel refuses `getsockopt` in the
current socket state (e.g. the `TCP_REPAIR_*` options outside repair mode); the
option is still settable. The `HINT` column is populated only when sox can
derive a friendlier label from the raw value — see the `SO_BINDTOIFINDEX`
example below.

### 3. Set a socket option
```bash
sudo sox set 1062 3 SO_KEEPALIVE 1

SOCKET_OPTION   VALUE   HINT    DESCRIPTION
SO_KEEPALIVE    1               Enable (1) / disable (0) periodic keepalive probes for dead-peer detection.
```

### 4. Get a socket option
```bash
sudo sox get 1062 3 SO_KEEPALIVE

SOCKET_OPTION   VALUE   HINT    DESCRIPTION
SO_KEEPALIVE    0               Enable (1) / disable (0) periodic keepalive probes for dead-peer detection.
```

`-o json` and `-o yaml` are also accepted on every subcommand for structured
output.

## Bind a socket to a specific network interface

`SO_BINDTOIFINDEX` lets you pin an existing socket to a NIC by its kernel
ifindex (useful for split-horizon proxying, VRFs, namespace-bound upstreams,
etc.). Look up the ifindex with `ip -o link` — the leading integer on each
line is the ifindex:

```bash
ip -o link

1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 ...
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 ...
3: wg0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 ...
```

For a single interface, `cat /sys/class/net/eth0/ifindex` returns the same
value. Bind the socket to it:

```bash
sudo sox set 1062 3 SO_BINDTOIFINDEX 2

SOCKET_OPTION       VALUE   HINT    DESCRIPTION
SO_BINDTOIFINDEX    2       eth0    Bind socket to a network interface by ifindex; 0 unbinds. Needs CAP_NET_RAW.
```

The `HINT` column resolves the ifindex back to a name so you can verify the
binding. Pass `0` to remove the binding:

```bash
sudo sox set 1062 3 SO_BINDTOIFINDEX 0
```

Setting this option requires `CAP_NET_RAW` in addition to the usual
`CAP_SYS_PTRACE` needed to reach a foreign process.

## Show option details with `sox explain`

`sox explain` prints the level, value range, summary and a long-form
description for a single option, without touching any process:

```bash
sox explain SO_BINDTOIFINDEX

SO_BINDTOIFINDEX
================

Level: SOL_SOCKET
Range: [0, 2147483647]

Summary:
  Bind socket to a network interface by ifindex; 0 unbinds. Needs CAP_NET_RAW.

Details:
  Restricts a socket to a specific network interface, identified by its
  kernel ifindex. Look up an interface's ifindex with "ip -o link" or
  "cat /sys/class/net/<iface>/ifindex".
  ...
```

Pass `-o json` or `-o yaml` for a structured form suitable for scripting:

```bash
sox -o yaml explain TCP_NODELAY
```

See the built-in help (`sox --help`) for the full command list.

## Known problems

- On Linux kernels older than 5.7, `getsockopt(SO_BINDTOIFINDEX)` returns
  `ENOPROTOOPT` even though `setsockopt` works (the read side was added later
  than the write side). `sox list` and `sox get` will show `n/a` for this
  option on such kernels; `sox set SO_BINDTOIFINDEX <ifindex>` still works.
