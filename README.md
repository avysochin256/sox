# SOX - socket options extractor/setter.

Sox allows to get or modify any tcp socket option for any process. Sometimes may be usefull for debugging.
Also sox may be used as quick fix as it allows to change socket options of any app without restart and downtime.
---


## Requirements:
- Linux kernel 5.7+
- `CAP_SYS_PTRACE` capability

## Instalation:
```
go install github.com/valexz/sox@latest
```

## Usage example


### Get pid and fd of socket via ss:
```
❯ sudo ss -ntpa
State       Recv-Q       Send-Q     Local Address:Port           Peer Address:Port        Process                                               
LISTEN      0            128        0.0.0.0:22                   0.0.0.0:*                users:(("sshd",pid=1062,fd=3))                       
```

### List all socket options for sshd process with PID=1062 and  tcp socket with FD=3
```
❯ sudo ./sox list 1062 3
SOCKET:                 0.0.0.0:22

OPTIONS:                NAME            VALUE                                   DESCRIPTION
SO_KEEPALIVE            0               Enable or disable TCP keepalive        
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
TCP_TIMESTAMP           2677060         Enable TCP timestamps                  
```

### Set socket option SO_KEEPALIVE to 1 for sshd process with PID=1062 and  tcp socket with FD=3
```
❯ sudo ./sox set 1062 3 SO_KEEPALIVE 1
SOCKET_OPTION   VALUE   DESCRIPTION                    
SO_KEEPALIVE    1       Enable or disable TCP keepalive
```

### Get value of socket option SO_KEEPALIVE for sshd process with PID=1062 and  tcp socket with FD=3
```
❯ sudo ./sox get 1062 3 SO_KEEPALIVE
SOCKET_OPTION   VALUE   DESCRIPTION                    
SO_KEEPALIVE    0       Enable or disable TCP keepalive
``
