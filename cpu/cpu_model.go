package cpu

type Register struct {
	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
	f FlagRegister
	h uint8
	l uint8
}

type CPU struct {
	Registers Register
}

type FlagRegister struct {
	zero      bool
	subtract  bool
	halfCarry bool
	carry     bool
}

type Instruction struct {
	opt       string
	targetReg string
}
