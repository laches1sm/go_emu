package cpu

import (
	"math"
	"math/bits"
)

type register interface {
	GetAF() uint16
	SetAF(uint16)
	GetBC() uint16
	SetBC(uint16)
	GetDE() uint16
	SetDE(uint16)
	GetHL() uint16
	SetHL(uint16)
}

type cpu interface {
	Execute(instruction Instruction)
	Add(value uint8)
	AddHL(value uint8)
	Subtract(value uint8)
}

const (
	aReg = "a"
	bReg = "b"
	cReg = "c"
	dReg = "d"
	eReg = "e"
	hReg = "h"
	lReg = "l"

	ADD   = "ADD"
	ADDHL = "ADDHL"
	ADDC  = "ADDC"
	SUB   = "SUB"
	SBC   = "SBC"
	AND   = "AND"
	OR    = "OR"
	XOR   = "XOR"
	CP    = "CP"
	INC   = "INC"
	DEC   = "DEC"
	CCF   = "CCF"
	SCF   = "SCF"
	RRA   = "RRA"
	RLA   = "RLA"
	RRCA  = "RRCA"
	RRLA  = "RRLA"
	CPL   = "CPL"
	BIT   = "BIT"
	RESET = "RESET"
	SET   = "SET"
	SRL   = "SRL"
	RR    = "RR"
	RL    = "RL"
	RRC   = "RRC"
	RLC   = "RLC"
	SRA   = "SRA"
	SLA   = "SLA"
	SWAP  = "SWAP"

	zeroFlagBytePosition      uint8 = 7
	subtractFlagBytePosition  uint8 = 6
	halfCarryFlagBytePosition uint8 = 5
	carryFlagBytePosition     uint8 = 4
)

func convertFlagRegUint8(flagRegister FlagRegister) uint8 {
	if flagRegister.zero {
		return 1 << zeroFlagBytePosition
	}
	if flagRegister.subtract {
		return 1 << subtractFlagBytePosition
	}
	if flagRegister.halfCarry {
		return 1 << halfCarryFlagBytePosition
	}
	if flagRegister.carry {
		return 1 << carryFlagBytePosition
	}
	return 0
}

func convertUint8FlagReg(val uint8) FlagRegister {
	zero := (val >> zeroFlagBytePosition & 0b1) != 0
	subtract := (val >> subtractFlagBytePosition & 0b1) != 0
	halfCarry := (val >> halfCarryFlagBytePosition & 0b1) != 0
	carry := (val >> carryFlagBytePosition & 0b1) != 0

	return FlagRegister{
		zero:      zero,
		subtract:  subtract,
		halfCarry: halfCarry,
		carry:     carry,
	}

}
func (cpu *CPU) Add(value uint8) uint8 {
	// first, check if value has overflow
	value = value + cpu.Registers.a
	newValue, overflow := overflowCheck(value)
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = overflow
	cpu.Registers.f.halfCarry = (cpu.Registers.a&0xF)+(value&0xF) > 0xF
	return newValue

}

func (cpu *CPU) AddHL(value uint8) uint8 {
	value = value + cpu.Registers.h + cpu.Registers.l
	newValue, overflow := overflowCheck(value)
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = overflow
	cpu.Registers.f.halfCarry = (cpu.Registers.h&0xF)+(cpu.Registers.l&0xF)+(value&0xF) > 0xF
	return newValue
}

func (cpu *CPU) AddC(value uint8) uint8 {
	var carryVal uint8
	value = value + cpu.Registers.a
	newValue, overflow := overflowCheck(value)
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = overflow
	if cpu.Registers.f.carry {
		carryVal = 1 << carryFlagBytePosition
	}
	cpu.Registers.f.halfCarry = (cpu.Registers.a&0xF)+(value&0xF)+(carryVal&0xF) > 0xF

	newValue = newValue + carryVal
	return newValue
}

// TODO: Note to self: CP is similar to this, only that the result does not get stored back into the a register after completion.
// Just use this function, only afterwards don't store the result in a.
func (cpu *CPU) Subtract(value uint8) uint8 {
	value = value&0xF - cpu.Registers.a&0xF
	newValue, overflow := overflowCheck(value)
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = true
	cpu.Registers.f.carry = overflow
	cpu.Registers.f.halfCarry = (cpu.Registers.a&0xF)-(value&0xF) > 0xF
	return newValue
}

func (cpu *CPU) SubtractC(value uint8) uint8 {
	var carryVal uint8
	value = value&0xF - cpu.Registers.a&0xF
	newValue, overflow := overflowCheck(value)
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = true
	cpu.Registers.f.carry = overflow
	if cpu.Registers.f.carry {
		carryVal = 1 << carryFlagBytePosition
	}
	cpu.Registers.f.halfCarry = (cpu.Registers.a&0xF)-(value&0xF)-(carryVal) > 0xF
	return newValue
}

func (cpu *CPU) And(value uint8) {
	andValue := cpu.Registers.a & value
	cpu.Registers.a = andValue
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = false
	cpu.Registers.f.halfCarry = false
}

func (cpu *CPU) Or(value uint8) {
	orValue := cpu.Registers.a | value
	cpu.Registers.a = orValue
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = false
	cpu.Registers.f.halfCarry = false
}

func (cpu *CPU) Xor(value uint8) {
	xorValue := cpu.Registers.a ^ value
	cpu.Registers.a = xorValue
	cpu.Registers.f.zero = false
	cpu.Registers.f.subtract = false
	cpu.Registers.f.carry = false
	cpu.Registers.f.halfCarry = false
}

func (cpu *CPU) Inc(value uint8) uint8 {
	return value + 1
}

func (cpu *CPU) Dec(value uint8) uint8 {
	return value - 1
}

func (cpu *CPU) ComplementCarryFlag(flagVal bool) {
	cpu.Registers.f.carry = flagVal
}

func (cpu *CPU) SetCarryFlag() {
	cpu.Registers.f.carry = true
}
func (cpu *CPU) RotateRight(regVal uint8) uint8 {
	regVal = bits.RotateLeft8(regVal, int(-carryFlagBytePosition))
	return regVal
}
func (cpu *CPU) RotateLeft(regVal uint8) uint8 {
	regVal = bits.RotateLeft8(regVal, int(carryFlagBytePosition))
	return regVal
}
func (cpu *CPU) RotateRightRegisterNotCarry(rotVal, regVal uint8) uint8 {
	regVal = bits.RotateLeft8(regVal, int(-rotVal))
	return regVal
}

func (cpu *CPU) RotateLeftRegisterNotCarry(rotVal, regVal uint8) uint8 {
	regVal = bits.RotateLeft8(regVal, int(rotVal))
	return regVal
}

func (cpu *CPU) RotateRightARegister() {
	// Golang is weird, in that the bits package lets you rotate left, but to rotate right you have to supply the the carry as a negative.
	// There's a good reason behind this, right? I guess it makes sense, but lazy old me wants an actual function, goddamnit!!
	newAVal := bits.RotateLeft8(cpu.Registers.a, int(-carryFlagBytePosition))
	cpu.Registers.a = newAVal
}

func (cpu *CPU) RotateLeftARegister() {
	newAVal := bits.RotateLeft8(cpu.Registers.a, int(carryFlagBytePosition))
	cpu.Registers.a = newAVal
}

func (cpu *CPU) RotateRightARegisterNotCarry(rotVal uint8) {
	newAVal := bits.RotateLeft8(cpu.Registers.a, int(-rotVal))
	cpu.Registers.a = newAVal
}

func (cpu *CPU) RotateLeftARegisterNotCarry(rotVal uint8) {
	newAVal := bits.RotateLeft8(cpu.Registers.a, int(rotVal))
	cpu.Registers.a = newAVal
}

func (cpu *CPU) Complement(val uint8) {
	cpu.Registers.a = ^val
}
func (cpu *CPU) BitSet(regVal uint8, pos uint8) bool {
	val := regVal & (1 << pos)
	return (val > 0)
}

func (cpu *CPU) Reset(regVal uint8, pos uint8) uint8 {
	var mask uint8
	mask = ^(1 << pos)
	regVal &= mask
	return regVal
}
func (cpu *CPU) Set(regVal uint8, pos uint8) uint8 {
	regVal |= (1 << pos)
	return regVal
}

func (cpu *CPU) ShiftRightLogical (regVal uint8) uint8{
	return regVal >> 1
}

func (cpu *CPU) ShiftLeftArithmetic(regVal uint8) uint8{
	return regVal << 1
}
func overflowCheck(value uint8) (uint8, bool) {
	if value > math.MaxUint8 {
		return value, true
	} else {
		return value, false
	}

}

func (cpu *CPU) Execute(instruction Instruction) {
	var value, newVal uint8
	switch instruction.opt {
	case ADD:
		switch instruction.targetReg {
		case aReg:
			value = cpu.Registers.a
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case bReg:
			// TODO: Refactor this!!
			value = cpu.Registers.b
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case cReg:
			value = cpu.Registers.c
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case dReg:
			value = cpu.Registers.d
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case eReg:
			value = cpu.Registers.e
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case hReg:
			value = cpu.Registers.h
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		case lReg:
			value = cpu.Registers.l
			newVal = cpu.Add(value)
			cpu.Registers.a = newVal
		default:
			break

		}
	case ADDHL:
		switch instruction.targetReg {
		case aReg:
			value = cpu.Registers.a
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case bReg:
			// TODO: Refactor this!!
			value = cpu.Registers.b
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case cReg:
			value = cpu.Registers.c
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case dReg:
			value = cpu.Registers.d
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case eReg:
			value = cpu.Registers.e
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case hReg:
			value = cpu.Registers.h
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		case lReg:
			value = cpu.Registers.l
			newVal = cpu.AddHL(value)
			cpu.Registers.SetHL(uint16(newVal))
		default:
			break
		}

	}

}
func (r *Register) GetAF() uint16 {
	fRegVal := convertFlagRegUint8(r.f)
	return uint16(r.a)<<8 | uint16(fRegVal)
}

func (r *Register) SetAF(value uint16) {
	fRegVal := convertFlagRegUint8(r.f)
	r.a = uint8(value & 0xFF00 >> 8)
	fRegVal = uint8(value & 0xFF)
	r.f = convertUint8FlagReg(fRegVal)
}

func (r *Register) GetBC() uint16 {
	return uint16(r.b)<<8 | uint16(r.c)
}

func (r *Register) SetBC(value uint16) {
	r.b = uint8(value & 0xFF00 >> 8)
	r.c = uint8(value & 0xFF)
}

func (r *Register) GetDE() uint16 {
	return uint16(r.d)<<8 | uint16(r.e)
}

func (r *Register) SetDE(value uint16) {
	r.d = uint8(value & 0xFF00 >> 8)
	r.e = uint8(value & 0xFF)
}

func (r *Register) GetHL() uint16 {
	return uint16(r.h)<<8 | uint16(r.l)
}

func (r *Register) SetHL(value uint16) {
	r.h = uint8(value & 0xFF00 >> 8)
	r.l = uint8(value & 0xFF)
}
