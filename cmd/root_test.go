/*
Copyright © 2024 Alexander Vysochin <avyssochin@gmail.com>
Licensed under the Apache License, Version 2.0
*/

package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func requireRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Fatal("This test requires root privileges (CAP_SYS_PTRACE capability). Please run with sudo.")
	}
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantOutput string
	}{
		{
			name:    "no args shows help",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
		},
		{
			name:    "version flag",
			args:    []string{"--version"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("root command error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(output) == 0 && !tt.wantErr {
				t.Error("expected output for help/version, got none")
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	requireRoot(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"get"},
			wantErr: true,
		},
		{
			name:    "missing fd",
			args:    []string{"get", "1234"},
			wantErr: true,
		},
		{
			name:    "missing option",
			args:    []string{"get", "1234", "3"},
			wantErr: true,
		},
		{
			name:    "invalid pid",
			args:    []string{"get", "invalid", "3", "TCP_NODELAY"},
			wantErr: true,
		},
		{
			name:    "invalid fd",
			args:    []string{"get", "1234", "invalid", "TCP_NODELAY"},
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"get", "1234", "3", "TCP_NODELAY", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("get command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCommand(t *testing.T) {
	requireRoot(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"set"},
			wantErr: true,
		},
		{
			name:    "missing fd",
			args:    []string{"set", "1234"},
			wantErr: true,
		},
		{
			name:    "missing option",
			args:    []string{"set", "1234", "3"},
			wantErr: true,
		},
		{
			name:    "missing value",
			args:    []string{"set", "1234", "3", "TCP_NODELAY"},
			wantErr: true,
		},
		{
			name:    "invalid pid",
			args:    []string{"set", "invalid", "3", "TCP_NODELAY", "1"},
			wantErr: true,
		},
		{
			name:    "invalid fd",
			args:    []string{"set", "1234", "invalid", "TCP_NODELAY", "1"},
			wantErr: true,
		},
		{
			name:    "invalid value",
			args:    []string{"set", "1234", "3", "TCP_NODELAY", "invalid"},
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"set", "1234", "3", "TCP_NODELAY", "1", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("set command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListCommand(t *testing.T) {
	requireRoot(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{"list"},
			wantErr: true,
		},
		{
			name:    "missing fd",
			args:    []string{"list", "1234"},
			wantErr: true,
		},
		{
			name:    "invalid pid",
			args:    []string{"list", "invalid", "3"},
			wantErr: true,
		},
		{
			name:    "invalid fd",
			args:    []string{"list", "1234", "invalid"},
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"list", "1234", "3", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(rootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("list command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
