package utils

import (
	"fmt"
	"testing"
    //"os"
    //"syscall"
    "time"
)
func Test_dbg(t *testing.T) {
	dbg, err := Debug("../../../bin/walua",[]string{""})
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 1)
	regs, err := dbg.GetRegs()
	if err != nil {
		fmt.Println(err)
	}
	if regs != nil {
		fmt.Println(regs)
	}
}

/*
func Test_dbg(t *testing.T) {
	dbg, err := Debug("../../../bin/walua")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbg.Close()
	for _ = range dbg.Events() {
		time.Sleep(time.Second * 1)
		regs, err := dbg.GetRegs()
		if err != nil {
			fmt.Println(err)
		}
		if regs != nil {
			fmt.Println(regs)
		}
		dbg.Continue()
	}
	fmt.Println("exit")
}

func Test_dbg2(t *testing.T) {
	go func() {
		fmt.Println("start")
		p, e := os.StartProcess("walua", []string{""}, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Sys: &syscall.SysProcAttr{
				Ptrace: true,
			},
		})
		if e != nil {
			fmt.Println(e)
		}

		s, e := p.Wait()
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(s)

		go func() {
			var regs syscall.PtraceRegs
			pid := p.Pid
			if e = syscall.PtraceGetRegs(pid, &regs); e != nil {
				fmt.Println(e)
			}
			fmt.Println(regs)
		}()
		fmt.Println("end")
	}()
	time.Sleep(time.Second * 1)
}
*/
