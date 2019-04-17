package utils

import (
	"fmt"
	"testing"
	//"os"
	//"syscall"
)

func Test_dbg(t *testing.T) {
	dbg, err := Debug("../../../bin/golua", []string{""})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbg.Close()
	for _ = range dbg.Events() {
		regs, err := dbg.GetRegs()
		if err != nil {
			fmt.Println(err)
		}
		if regs != nil {
			fmt.Printf("%+v\n", regs)
		}
		var txt = make([]byte, 20)
		_, err = dbg.PeekText(uintptr(regs.Pc), txt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%x:%v\n", regs.Pc, txt)
		dbg.Continue()
	}
	fmt.Println("exit")
}
