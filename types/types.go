package types

import (
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
	name           string
	guess          float64
	upperlim       float64
	lowerlim       float64
	comment        string
	unit           quantity.Unit
	dimensionality []int8
	solved         bool
}

type Equation struct {
	index uint16
	line  uint16
	text  string
	vars  []*Variable
}
