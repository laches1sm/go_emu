package cpu

type register interface {
	GetAF() uint16
	SetAF(uint16)
	GetBC() uint16
	SetBC(uint16)
	GetDE() uint16
	SetDE(uint16)
	GetHL()uint16
	SetHL(uint16)

}

type cpu interface{
	Execute(instruction Instruction)
}
type ArithemeticTarget string

type Instruction string

const(
	a ArithemeticTarget = "a"
	b ArithemeticTarget = "b"
	c ArithemeticTarget = "c"
	d ArithemeticTarget = "d"
	e ArithemeticTarget = "e"
	h ArithemeticTarget = "h"
	l ArithemeticTarget = "l"

	ADD Instruction = "ADD"
)


func(c *CPU) Execute(instruction Instruction, target ArithemeticTarget){
	switch instruction {
	case ADD:




	}

}
func (r *Register) GetAF() uint16{
	return uint16(r.a) << 8 | uint16(r.f)
}

func (r *Register) SetAF(value uint16){
	r.a = uint8(value & 0xFF00 >> 8)
	r.f = uint8(value & 0xFF)
}

func (r *Register) GetBC() uint16{
	return uint16(r.b) << 8 | uint16(r.c)
}

func (r *Register) SetBC(value uint16){
	r.b = uint8(value & 0xFF00 >> 8)
	r.c = uint8(value & 0xFF)
}

func (r *Register) GetDE() uint16{
	return uint16(r.d) << 8 | uint16(r.e)
}

func (r *Register) SetDE(value uint16){
	r.d = uint8(value & 0xFF00 >> 8)
	r.e = uint8(value & 0xFF)
}

func (r *Register) GetHL() uint16{
	return uint16(r.h) << 8 | uint16(r.l)
}

func (r *Register) SetHL(value uint16){
	r.h = uint8(value & 0xFF00 >> 8)
	r.l = uint8(value & 0xFF)
}