package utils

import (
	"os"
	"runtime"
	"syscall"
)

const (
	CREATE_NEW_CONSOLE      = 0x10
	DEBUG_ONLY_THIS_PROCESS = 0x2
	DEBG_PROCESS            = 0x1
)

func LoadProcess(name string, arg ...string) (*os.Process, error) {
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	attr.Sys = &syscall.SysProcAttr{}
    attr.Sys.CreationFlags = DEBUG_ONLY_THIS_PROCESS
	return os.StartProcess(name, arg, &attr)
}
