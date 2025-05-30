# SOX - Socket Options eXtractor

SOX is a command-line tool that allows you to view and modify TCP socket options for any process. It's particularly useful for debugging network issues and tuning network performance without requiring application restarts.

## Features

- List all available socket options for a given socket
- Get specific socket option values
- Set socket option values
- Supports a wide range of TCP and socket options
- Validates option values against allowed ranges
- No application restart required
- Clear, tabulated output

## Requirements

- Linux kernel 5.6+
- `CAP_SYS_PTRACE` capability (run with sudo)
- Go 1.20 or later (for building from source)

## Installation

### Using Go Install
```bash
go install github.com/valexz/sox@latest
```

### Building from Source
```bash
git clone https://github.com/valexz/sox.git
cd sox
go build
```

## Usage

### List All Socket Options
```bash
# First, find the process and socket file descriptor using ss
sudo ss -ntpa
State       Recv-Q       Send-Q     Local Address:Port     Peer Address:Port     Process                                               
LISTEN      0            128        0.0.0.0:22             0.0.0.0:*            users:(("sshd",pid=1062,fd=3))

# List all socket options for the process
sudo sox list 1062 3
OPTION NAME             VALUE           DESCRIPTION                            
SO_KEEPALIVE            1               Enable TCP keepalive                   
TCP_KEEPIDLE            7200            Time before sending keepalive probes   
TCP_KEEPINTVL          75              Time between keepalive probes          
TCP_KEEPCNT             9               Number of keepalive probes             
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
TCP_TIMESTAMP           19100429        Enable TCP timestamps            
```

### Get Specific Option
```bash
sudo sox get 1062 3 TCP_NODELAY
SOCKET_OPTION   VALUE   DESCRIPTION                    
TCP_NODELAY     0       Disable Nagle's algorithm
```

### Set Option Value
```bash
sudo sox set 1062 3 TCP_NODELAY 1
SOCKET_OPTION   VALUE   DESCRIPTION                    
TCP_NODELAY     1       Disable Nagle's algorithm
```

## Supported Socket Options

### TCP Options
- `TCP_NODELAY`: Disable Nagle's algorithm (0/1)
- `TCP_MAXSEG`: Maximum segment size (536-65535 bytes)
- `TCP_KEEPIDLE`: Time before sending keepalive probes (1-32767 seconds)
- `TCP_KEEPINTVL`: Time between keepalive probes (1-32767 seconds)
- `TCP_KEEPCNT`: Number of keepalive probes (1-127)
- `TCP_CORK`: Don't send partial frames (0/1)
- `TCP_DEFER_ACCEPT`: Delay accept() until data arrives (0-32767 seconds)
- `TCP_WINDOW_CLAMP`: Bound advertised window
- `TCP_FASTOPEN`: Enable TCP Fast Open
- `TCP_QUICKACK`: Enable quickack mode (0/1)

### Socket Options
- `SO_KEEPALIVE`: Enable TCP keepalive (0/1)
- `SO_RCVBUF`: Socket receive buffer size (≥2048 bytes)
- `SO_SNDBUF`: Socket send buffer size (≥2048 bytes)
- `SO_RCVLOWAT`: Minimum receive buffer space available
- `SO_SNDLOWAT`: Minimum send buffer space available
- `SO_REUSEADDR`: Allow reuse of local addresses (0/1)
- `SO_REUSEPORT`: Allow multiple sockets to bind to same address/port (0/1)

## Common Use Cases

1. **Debugging Connection Issues**
   ```bash
   # Check keepalive settings
   sudo sox get 1234 3 SO_KEEPALIVE
   sudo sox get 1234 3 TCP_KEEPIDLE
   sudo sox get 1234 3 TCP_KEEPINTVL
   sudo sox get 1234 3 TCP_KEEPCNT
   ```

2. **Optimizing for Low Latency**
   ```bash
   # Enable TCP_NODELAY (disable Nagle's algorithm)
   sudo sox set 1234 3 TCP_NODELAY 1
   # Enable quick acknowledgments
   sudo sox set 1234 3 TCP_QUICKACK 1
   ```

3. **Tuning Buffer Sizes**
   ```bash
   # Increase receive buffer
   sudo sox set 1234 3 SO_RCVBUF 262144
   # Increase send buffer
   sudo sox set 1234 3 SO_SNDBUF 262144
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Author

Alexander Vysochin (<avyssochin@gmail.com>)