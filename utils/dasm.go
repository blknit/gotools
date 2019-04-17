package utils

import (
	"errors"
	"golang.org/x/arch/arm/armasm"
	"golang.org/x/arch/arm64/arm64asm"
	//"golang.org/x/arch/ppc64/ppc64asm"
	"golang.org/x/arch/x86/x86asm"
)

type Syntax int

const (
	_ Syntax = iota
	Go
	GNU
	Intel
)

type CPU int

const (
	_ CPU = iota
	X86
	X86_64
	ARM
	ARM_64
)

var (
	ErrInvalidCPU    = errors.New("invalidat cpu")
	ErrInvalidMode   = errors.New("invalidat mode")
	ErrInvalidSyntax = errors.New("invalidat syntax")
)

type IDasm interface {
	Dasm(bin []byte, pc uint64) (string, error)
}

func NewDasm(cpu CPU, syntax Syntax) (IDasm, error) {
	switch cpu {
	case X86:
		return NewDasmX86(32, syntax)
	case X86_64:
		return NewDasmX86(64, syntax)
	case ARM:
		return NewDasmArm(syntax)
	}
	return nil, ErrInvalidCPU
}

/////////////////////////////////////////////////////////////////////////
type DasmX86 struct {
	mode   int
	syntax Syntax
}

func NewDasmX86(mode int, syntax Syntax) (*DasmX86, error) {
	if mode != 16 && mode != 32 && mode != 64 {
		return nil, ErrInvalidMode
	}
	if (syntax != Go) && (syntax != GNU) && (syntax != Intel) {
		return nil, ErrInvalidSyntax
	}
	d := DasmX86{
		mode:   mode,
		syntax: syntax,
	}
	return &d, nil
}

func (a *DasmX86) Dasm(bin []byte, pc uint64) (string, error) {
	inst, err := x86asm.Decode(bin, a.mode)
	if err != nil {
		return "", err
	}
	switch a.syntax {
	case Go:
		return x86asm.GoSyntax(inst, pc, nil), nil
	case GNU:
		return x86asm.GNUSyntax(inst, pc, nil), nil
	case Intel:
		return x86asm.IntelSyntax(inst, pc, nil), nil
	}
	return "", ErrInvalidSyntax
}

/////////////////////////////////////////////////////////////////////////
type DasmArm struct {
	mode   int
	syntax Syntax
}

func NewDasmArm(syntax Syntax) (*DasmArm, error) {
	if syntax != Go && syntax != GNU {
		return nil, ErrInvalidSyntax
	}
	d := DasmArm{
		mode:   int(armasm.ModeARM),
		syntax: syntax,
	}
	return &d, nil
}

func (a *DasmArm) Dasm(bin []byte, pc uint64) (string, error) {
	inst, err := armasm.Decode(bin, armasm.Mode(a.mode))
	if err != nil {
		return "", err
	}
	switch a.syntax {
	case Go:
		return armasm.GoSyntax(inst, pc, nil, nil), nil
	case GNU:
		return armasm.GNUSyntax(inst), nil
	}
	return "", ErrInvalidSyntax
}

/////////////////////////////////////////////////////////////////////////
type DasmArm64 struct {
	syntax Syntax
}

func NewDasmArm64(syntax Syntax) (*DasmArm64, error) {
	if syntax != Go && syntax != GNU {
		return nil, ErrInvalidSyntax
	}
	d := DasmArm64{
		syntax: syntax,
	}
	return &d, nil
}

func (a *DasmArm64) Dasm(bin []byte, pc uint64) (string, error) {
	inst, err := arm64asm.Decode(bin)
	if err != nil {
		return "", err
	}
	switch a.syntax {
	case Go:
		return arm64asm.GoSyntax(inst, pc, nil, nil), nil
	case GNU:
		return arm64asm.GNUSyntax(inst), nil
	}
	return "", ErrInvalidSyntax
}
