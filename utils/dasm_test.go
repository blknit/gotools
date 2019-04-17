package utils

import (
	"testing"
	"gotest.tools/assert"
)

func Test_asmx86(t *testing.T) {
	// test x64 + intel
	d, err := NewDasm(X86_64, Intel)
	if err != nil {
		t.Error(err)
	}
	str, err := d.Dasm([]byte{0x48, 0xC7, 0x44, 0x24, 0x20, 0x00, 0x00, 0x00, 0x00}, 0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, str, "mov qword ptr [rsp+0x20], 0x0")

	// test x86 + intel
	d, err = NewDasm(X86, Intel)
	if err != nil {
		t.Error(err)
	}
	str, err = d.Dasm([]byte{0xC7, 0x44, 0x24, 0x20, 0x00, 0x00, 0x00, 0x00}, 0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, str, "mov dword ptr [esp+0x20], 0x0")
}
