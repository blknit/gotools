package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
	//"reflect"
	sys "golang.org/x/sys/unix"
	//ui "github.com/gizak/termui"
	//"github.com/gizak/termui/widgets"
)

type regs struct {
	regs   [31]uint64
	sp     uint64
	pc     uint64
	pstate uint64
}

func main() {
	if false {
		go func() {
			var procAttr os.ProcAttr
			procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
			procAttr.Sys = &syscall.SysProcAttr{Ptrace: true}
			p, e := os.StartProcess("walua", []string{""}, &procAttr)
			if e != nil {
				fmt.Println(e)
				return
			}
			fmt.Println(p.Pid)
			s, e := p.Wait()
			if e != nil {
				fmt.Println(e)
				return
			}
			fmt.Println(s)
			var r regs
			iov := sys.Iovec{Base: (*byte)(unsafe.Pointer(&r)), Len: uint64(unsafe.Sizeof(r))}
			_, _, e = syscall.Syscall6(syscall.SYS_PTRACE, sys.PTRACE_GETREGSET, uintptr(p.Pid), 1, uintptr(unsafe.Pointer(&iov)), 0, 0)
			if e != syscall.Errno(0) {
				fmt.Println(e)
			}
			fmt.Println(r)
			/*
			   var xstateargs[512]byte
			   iov := sys.Iovec{Base:&xstateargs[0],Len:512}
			   _,_,e=syscall.Syscall6(syscall.SYS_PTRACE,sys.PTRACE_GETREGSET,uintptr(p.Pid),1,uintptr(unsafe.Pointer(&iov)),0,0)
			   if e!=syscall.Errno(0){
			       fmt.Println(e)
			   }
			   fmt.Println(xstateargs)
			*/
			fmt.Println("!!!hello world!!!")
		}()
		time.Sleep(time.Second * 1)
	}
	fmt.Println("hello world")
}
