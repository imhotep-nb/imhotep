package types

import (
	"github.com/antonmedv/expr/vm"
	"github.com/imhotep-nb/units/quantity"
)

type Variable struct {
	/*
		name: Name of variable eg: x
		guess: Current guess (on solved true, value for zero of functions)
			eg: 2.3
		upperlim: Upper limit for iterate value of guess eg: 100
		lowerlim: Lower limit for iterate value of guess eg: -100
		comment: Comment text for the current variable eg: kill me plz
		unit: Dimensional unit for variable (github.com/imhotep-nb/units/quantity)
			eg: "m.kg"
		dimensionality: Exponents for physical dimensionality
			index-> [0    1     2    3    4     5      6      7     8    9       10]
				"m", "kg", "K", "A", "cd", "mol", "rad", "sr", "¤", "byte", "s"
			eg: - m2.kg:            [2 1 0 0 0 0 0 0 0 0 0]
				- km.h-1:           [1 0 0 0 0 0 0 0 0 0 -1]
				- MPa (kg.m-1.s-2): [-1 1 0 0 0 0 0 0 0 0 -2]
		solved: Define if the variable is determinated
	*/
	Name           string
	Guess          float64
	Upperlim       float64
	Lowerlim       float64
	Comment        string
	Unit           quantity.Unit
	Dimensionality []int8
	Solved         bool
}

type Equation struct {
	/*
		index: Equation number in the order that were writed inside
			equations box.
		line: Equation line number of the text writed inside equations box.
		text: Equation text writed by user.
		exec: expr compiled program from equation text.
		vars: pointers's slice to variables
	*/
	Index uint16
	Line  uint16
	Text  string
	Exec  vm.Program
	Vars  []*Variable
}
