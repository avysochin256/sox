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

## Getting started
1. Locate the target socket using `ss` or a similar tool to obtain the PID and
   file descriptor.

2. List all supported options for that socket:
   ```bash
   sudo sox list <pid> <fd>
   ```

3. Get the value of a single option:
   ```bash
   sudo sox get <pid> <fd> TCP_NODELAY
   ```

4. Change an option value:
   ```bash
   sudo sox set <pid> <fd> TCP_NODELAY 1
   ```

See the built in help (`sox --help`) for more commands and options.
