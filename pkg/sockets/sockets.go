package sockets

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SocketInfo holds information about a network socket.
type SocketInfo struct {
	Protocol   string
	LocalAddr  string
	RemoteAddr string
	State      string
	Inode      string
	PID        string
	FD         string
}

// Function to read and parse /proc/net/tcp or /proc/net/tcp6
func parseProcNet(protocol string) ([]SocketInfo, error) {
	file, err := os.Open(fmt.Sprintf("/proc/net/%s", protocol))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip the first line (header)
	scanner.Scan()

	var connections []SocketInfo
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localAddr := parseAddress(fields[1])
		remoteAddr := parseAddress(fields[2])
		state := parseState(fields[3])
		inode := fields[9]

		connection := SocketInfo{
			Protocol:   protocol,
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			State:      state,
			Inode:      inode,
		}
		connections = append(connections, connection)
	}

	return connections, scanner.Err()
}

// Function to parse IP address and port
func parseAddress(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return ""
	}
	ip := parts[0]
	port := parts[1]
	parsedIP := parseHexIP(ip)
	parsedPort, _ := strconv.ParseInt(port, 16, 64)
	return fmt.Sprintf("%s:%d", parsedIP, parsedPort)
}

// Function to parse hexadecimal IP address
func parseHexIP(hexIP string) string {
	bytes := make([]byte, len(hexIP)/2)
	for i := 0; i < len(bytes); i++ {
		b, _ := strconv.ParseUint(hexIP[i*2:i*2+2], 16, 8)
		bytes[len(bytes)-1-i] = byte(b)
	}
	return fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
}

// Function to parse TCP state
func parseState(state string) string {
	states := map[string]string{
		"01": "ESTABLISHED",
		"02": "SYN_SENT",
		"03": "SYN_RECV",
		"04": "FIN_WAIT1",
		"05": "FIN_WAIT2",
		"06": "TIME_WAIT",
		"07": "CLOSE",
		"08": "CLOSE_WAIT",
		"09": "LAST_ACK",
		"0A": "LISTEN",
		"0B": "CLOSING",
	}
	return states[state]
}

// Function to find the PID and FD from inode
func findPidFdFromInode(inode string) (string, string, error) {
	procDirs, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		return "", "", err
	}

	for _, procFd := range procDirs {
		link, err := os.Readlink(procFd)
		if err != nil {
			continue
		}
		if strings.Contains(link, fmt.Sprintf("socket:[%s]", inode)) {
			// Extract PID and FD
			parts := strings.Split(procFd, "/")
			if len(parts) >= 3 {
				pid := parts[2]
				fd := parts[4]
				return pid, fd, nil
			}
		}
	}

	return "", "", fmt.Errorf("inode %s not found", inode)
}

func getConnections() {
	tcp4Connections, err := parseProcNet("tcp")
	if err != nil {
		fmt.Println(err)
		return
	}

	tcp6Connections, err := parseProcNet("tcp6")
	if err != nil {
		fmt.Println(err)
		return
	}

	allConnections := append(tcp4Connections, tcp6Connections...)

	for i, connection := range allConnections {
		pid, fd, err := findPidFdFromInode(connection.Inode)
		if err != nil {
			fmt.Println(err)
			continue
		}
		allConnections[i].PID = pid
		allConnections[i].FD = fd
	}

	for _, conn := range allConnections {
		fmt.Printf("Protocol: %s, Local: %s, Remote: %s, State: %s, Inode: %s, PID: %s, FD: %s\n",
			conn.Protocol, conn.LocalAddr, conn.RemoteAddr, conn.State, conn.Inode, conn.PID, conn.FD)
	}
}
