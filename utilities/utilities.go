package utilities

import (
	"os"
	)

type Action struct {
	Type string
	Payload []byte
}

var PermissionsCodes = map[string] int {
	"rwxrr": 744,
	"rwx--": 700,
	"rwrr": 644,
	"rw--": 600,
	"r--": 400,
}

func CWD() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return "./"
}

type Error struct {
	errC chan error
}

// todo ..
func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}


