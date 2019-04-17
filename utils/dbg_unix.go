package utils

import (
	"errors"
	sys "golang.org/x/sys/unix"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// gitbug.com/tfogal/ptrace/blob/master/ptrace.go

var (
	DbgExited = errors.New("debuger exited")
)

type event interface{}

type Dbg struct {
	proc   *os.Process
	events chan event
	err    chan error
	cmds   chan func()
}

func Debug(name string, args []string) (*Dbg, error) {
	dbg := &Dbg{
		events: make(chan event, 1),
		err:    make(chan error, 1),
		cmds:   make(chan func(), 1),
	}
	err := make(chan error)
	go func() {
		// If the debugger itself is multi-threaded, ptrace calls must come from
		// the same thread that originally attached to the remote thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		p, e := os.StartProcess(name, args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Sys: &syscall.SysProcAttr{
				Ptrace: true,
			},
		})
		err <- e
		if e != nil {
			return
		}
		dbg.proc = p
		go dbg.wait()
		for cmd := range dbg.cmds {
			cmd()
		}
	}()
	e := <-err
	if e != nil {
		return nil, e
	}
	return dbg, nil
}

func Attach(id int) (*Dbg, error) {
	dbg := &Dbg{
		events: make(chan event, 1),
		err:    make(chan error, 1),
		cmds:   make(chan func(), 1),
	}
	err := make(chan error)
	go func() {
		// If the debugger itself is multi-threaded, ptrace calls must come from
		// the same thread that originally attached to the remote thread.
		runtime.LockOSThread()
		e := syscall.PtraceAttach(id)
		if e != nil {
			err <- e
		}
		p, e := os.FindProcess(id)
		err <- e
		if e != nil {
			return
		}
		dbg.proc = p
		go dbg.wait()
		for cmd := range dbg.cmds {
			cmd()
		}
	}()
	e := <-err
	if e != nil {
		return nil, e
	}
	return dbg, nil
}

func (a *Dbg) Close() error {
	close(a.cmds)
	a.cmds = nil
	return syscall.Kill(a.proc.Pid, syscall.SIGKILL)
}

func (a *Dbg) Detach() error {
	err := make(chan error, 1)
	if a.do(func() { err <- syscall.PtraceDetach(a.proc.Pid) }) {
		return <-err
	}
	return DbgExited
}

func (a *Dbg) Events() <-chan event {
	return a.events
}

func (a *Dbg) StepIn() error {
	err := make(chan error, 1)
	if a.do(func() { err <- syscall.PtraceSingleStep(a.proc.Pid) }) {
		return <-err
	}
	return DbgExited
}

func (a *Dbg) Continue() error {
	err := make(chan error, 1)
	if a.do(func() { err <- syscall.PtraceCont(a.proc.Pid, 0) }) {
		return <-err
	}
	return DbgExited
}

func (a *Dbg) GetRegs() (*syscall.PtraceRegs, error) {
	err := make(chan error, 1)
	reg := make(chan *syscall.PtraceRegs, 1)
	if a.do(func() {
		r := &syscall.PtraceRegs{}
		var e error
		if runtime.GOARCH == "arm64" {
			iov := sys.Iovec{Base: (*byte)(unsafe.Pointer(r)), Len: uint64(unsafe.Sizeof(*r))}
			_, _, e = syscall.Syscall6(syscall.SYS_PTRACE, sys.PTRACE_GETREGSET, uintptr(a.proc.Pid), 1, uintptr(unsafe.Pointer(&iov)), 0, 0)
			if e == syscall.Errno(0) {
				e = nil
			}
		} else {
			e = syscall.PtraceGetRegs(a.proc.Pid, r)
		}
		err <- e
		reg <- r
	}) {
		return <-reg, <-err
	}
	return nil, DbgExited
}

func (a *Dbg) do(f func()) bool {
	if a.cmds != nil {
		a.cmds <- f
		return true
	}
	return false
}

func (a *Dbg) wait() {
	defer close(a.err)
	for {
		stat, err := a.proc.Wait()
		if err != nil {
			a.err <- err
			close(a.events)
			return
		}
		if stat.Exited() {
			a.events <- event(stat)
			close(a.events)
			return
		}
		a.events <- event(stat)
	}
}
