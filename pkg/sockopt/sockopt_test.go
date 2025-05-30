/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/

package sockopt

import (
	"os"
	"testing"
)

func requireRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Fatal("This test requires root privileges (CAP_SYS_PTRACE capability). Please run with sudo.")
	}
}

func TestGetOptionByName(t *testing.T) {
	tests := []struct {
		name       string
		optionName string
		wantErr    bool
	}{
		{
			name:       "valid option TCP_NODELAY",
			optionName: "TCP_NODELAY",
			wantErr:    false,
		},
		{
			name:       "valid option SO_KEEPALIVE",
			optionName: "SO_KEEPALIVE",
			wantErr:    false,
		},
		{
			name:       "invalid option",
			optionName: "INVALID_OPTION",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := GetOptionByName(tt.optionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOptionByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && opt == nil {
				t.Errorf("GetOptionByName() returned nil option for valid name")
			}
		})
	}
}

func TestValidateSocketOption(t *testing.T) {
	tests := []struct {
		name       string
		optionName string
		value      int
		wantErr    bool
	}{
		{
			name:       "valid TCP_NODELAY value 0",
			optionName: "TCP_NODELAY",
			value:      0,
			wantErr:    false,
		},
		{
			name:       "valid TCP_NODELAY value 1",
			optionName: "TCP_NODELAY",
			value:      1,
			wantErr:    false,
		},
		{
			name:       "invalid TCP_NODELAY value",
			optionName: "TCP_NODELAY",
			value:      2,
			wantErr:    true,
		},
		{
			name:       "valid TCP_MAXSEG value",
			optionName: "TCP_MAXSEG",
			value:      1460,
			wantErr:    false,
		},
		{
			name:       "invalid TCP_MAXSEG value below min",
			optionName: "TCP_MAXSEG",
			value:      100,
			wantErr:    true,
		},
		{
			name:       "invalid TCP_MAXSEG value above max",
			optionName: "TCP_MAXSEG",
			value:      70000,
			wantErr:    true,
		},
		{
			name:       "invalid option name",
			optionName: "INVALID_OPTION",
			value:      1,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSocketOption(tt.optionName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSocketOption() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateProcess(t *testing.T) {
	requireRoot(t)

	tests := []struct {
		name    string
		pid     int
		wantErr bool
	}{
		{
			name:    "invalid pid 0",
			pid:     0,
			wantErr: true,
		},
		{
			name:    "invalid pid negative",
			pid:     -1,
			wantErr: true,
		},
		{
			name:    "valid pid 1 (init)",
			pid:     1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProcess(tt.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFD(t *testing.T) {
	tests := []struct {
		name    string
		fd      int
		wantErr bool
	}{
		{
			name:    "valid fd 0",
			fd:      0,
			wantErr: false,
		},
		{
			name:    "valid fd 3",
			fd:      3,
			wantErr: false,
		},
		{
			name:    "invalid fd negative",
			fd:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFD(tt.fd)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFD() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
