# Sox

Sox (Socket Options eXplorer) is a small command line tool that allows you to
inspect and modify TCP socket options of any running process. It can be helpful
for debugging or tuning network applications without restarting them.

## Requirements
- Linux kernel 5.6 or newer (for the pidfd API)
- `CAP_SYS_PTRACE` capability to operate on foreign processes

## Installation
```bash
go install github.com/valexz/sox@latest
```

## Quickstart

Follow the steps below to locate a socket and inspect or modify its options.

### 1. Get PID and file descriptor via `ss`
Run `ss -ntpa` and look at the `users` column:

```bash
sudo ss -ntpa
```

Example output:

```
State  Recv-Q Send-Q Local Address:Port Peer Address:Port Process
LISTEN 0      128    0.0.0.0:22         0.0.0.0:*       users:(("sshd",pid=1062,fd=3))
```

Here the PID is `1062` and the file descriptor is `3`.

### 2. List all socket options
```bash
sudo sox list 1062 3
OPTION NAME             VALUE           DESCRIPTION
SO_KEEPALIVE            1               Enable or disable TCP keepalive
TCP_KEEPIDLE            7200            Start keepalives after this period
TCP_KEEPINTVL           75              Interval between keepalives
TCP_KEEPCNT             9               Number of keepalives before death
TCP_NODELAY             0               Disable Nagle's algorithm
TCP_MAXSEG              536             Maximum segment size
TCP_CORK                0               Control sending of partial frames
TCP_SYNCNT              6               Number of SYN retransmits
TCP_LINGER2             60              Lifetime of orphaned FIN-WAIT-2 state
TCP_DEFER_ACCEPT        0               Wake up listener only when data arrives
TCP_WINDOW_CLAMP        0               Set maximum window size
TCP_INFO                10              Information about this socket
TCP_QUICKACK            1               Enable quick ACK
TCP_CONGESTION          1768060259      Get/Set congestion control algorithm
TCP_REPAIR              0               TCP repair mode
TCP_FASTOPEN            0               Enable TCP Fast Open
TCP_TIMESTAMP           19100429        Initial TCP timestamp value
```

### 3. Set a socket option
```bash
sudo sox set 1062 3 SO_KEEPALIVE 1
SOCKET_OPTION   VALUE   DESCRIPTION
SO_KEEPALIVE    1       Enable or disable TCP keepalive
```

### 4. Get a socket option
```bash
sudo sox get 1062 3 SO_KEEPALIVE
SOCKET_OPTION   VALUE   DESCRIPTION
SO_KEEPALIVE    0       Enable or disable TCP keepalive
```

See the built-in help (`sox --help`) for more commands and options.
