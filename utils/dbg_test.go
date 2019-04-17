package utils

import (
	"fmt"
	"testing"
	//"os"
	//"syscall"
	"time"
)

func Test_dbg(t *testing.T) {
	dbg, err := Debug("../../../bin/golua", []string{""})
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
