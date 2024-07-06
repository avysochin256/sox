package sockopt

import (
	"errors"
	"fmt"
	"github.com/gosuri/uitable"
	"golang.org/x/sys/unix"
	"log/slog"
	"os"
)

func GetSocketName(socketFd int) string {
	sn, _ := unix.Getsockname(socketFd)
	sai4 := sn.(*unix.SockaddrInet4)
	socketName := fmt.Sprintf("%d.%d.%d.%d:%d", sai4.Addr[0], sai4.Addr[1], sai4.Addr[2], sai4.Addr[3], sai4.Port)

	return socketName
}

func ListSocketOptions(pid, fd int) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	table := uitable.New()

	table.MaxColWidth = 50
	table.AddRow("OPTION NAME", "VALUE", "DESCRIPTION")

	var joinedListErr error
	for _, soname := range OptionsList {
		so := OptionsMap[soname]

		val, err := so.Get(socketFd)

		if err != nil {
			err = fmt.Errorf("unable to get sockopt option %s : %w", so.Name, err)
			errors.Join(joinedListErr, err)

			continue
		}

		table.AddRow(so.Name, val, so.Description)
	}

	fmt.Println(table)

	if uw, ok := joinedListErr.(interface{ Unwrap() []error }); ok {
		errs := uw.Unwrap()
		for _, err := range errs {
			slog.Error("unable to get value of sockopt option", slog.Any("error", err))
		}
	}

}

func SetSocketOption(pid, fd int, option string, val int) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("SOCKET_OPTION", "VALUE", "DESCRIPTION")

	so, ok := OptionsMap[option]
	if !ok {
		err = fmt.Errorf("unsupported socket option %s : %w", option, err)
		os.Exit(1)

	}

	err = so.Set(socketFd, val)
	if err != nil {
		err = fmt.Errorf("unable to set sockopt option %s : %w", so.Name, err)
		os.Exit(1)
	}

	val, err = so.Get(socketFd) // TODO: Remove shadowing

	if err != nil {
		err = fmt.Errorf("unable to get socket option %s  after value was set: %w", so.Name, err)

	}

	table.AddRow(so.Name, val, so.Description)

	fmt.Println(table)

}


func GetSocketOption(pid, fd int, option string) {

	socketFd, err := GetSocketFd(pid, fd)
	if err != nil {
		slog.Error("unable to get sockopt fd", slog.Any("error", err))
	}

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("SOCKET_OPTION", "VALUE", "DESCRIPTION")

	so, ok := OptionsMap[option]
	if !ok {
		err = fmt.Errorf("unsupported socket option %s : %w", option, err)
		os.Exit(1)

	}

	val, err := so.Get(socketFd)

	if err != nil {
		err = fmt.Errorf("unable to get socket option %s  after value was set: %w", so.Name, err)

	}

	table.AddRow(so.Name, val, so.Description)

	fmt.Println(table)

}