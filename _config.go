package main

import "flag"

type Config struct {
	targetPid int
	targetFD  int
}

func GetConfig() *Config {

	targetPid := flag.Int("pid", -1, "target PID")
	targetFD := flag.Int("fd", -1, "target FD")

	flag.Parse()

	c := &Config{
		targetPid: *targetPid,
		targetFD:  *targetFD,
	}

	return c
}
