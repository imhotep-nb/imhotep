package types

import (
	"log"
	"math"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/imhotep-nb/units/quantity"
	"gonum.org/v1/gonum/optimize"
)

type APIInput struct {
	/*
		Parse variables and Equations data from API in JSON format
	*/
	Equations []EquationJSON `json:"eqns"`
	Variables []VariableJSON `json:"vars"`
	Settings  SolverSettings `json:"settings"`
	Debug     bool           `json:"debug"`
}

type EquationJSON struct {
	/*
		Parse equations data from API in JSON format
	*/
	Text            string `json:"text"`
	Line            int    `json:"Line"`
	UnitsParsedText string
}

type VariableJSON struct {
	/*
		Parse variables data from API in JSON format
	*/
	Name     string  `json:"name"`
	Guess    float64 `json:"guess"`
	Upperlim float64 `json:"upperlim"`
	Lowerlim float64 `json:"lowerlim"`
	Comment  string  `json:"comment"`
	Unit     string  `json:"unit"`
}

type SolverSettings struct {
	/*
		Set of settings for equation solver

	*/
	// InitValues
	GradientThreshold float64 `json:"gradientThreshold"`
	// Converger
	MajorIterations int `json:"majorIterations"`
	Runtime         int `json:"runtime"`
	FuncEvaluations int `json:"funcEvaluations"`
	GradEvaluations int `json:"gradEvaluations"`
	// HessEvaluations
	Concurrent int `json:"concurrent"`
}

type Info struct {
	Logs     []string
	Errors   []string
	Warnings []string
	Msgs     []string
	Graph    map[interface{}][]interface{}
}

type BlockInfo struct {
	Stats  optimize.Stats
	Status string
}

type Stats struct {
	Blocks []BlockInfo
	Global optimize.Stats
}

type APIOutput struct {
	/*
		Struct for output of API
	*/
	Eqns     []EquationJSON
	Vars     []VariableJSON
	Settings SolverSettings
	Stats    Stats
	Info     Info
}

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
	Index          uint16
	Guess          float64
	TempValue      float64
	Upperlim       float64
	Lowerlim       float64
	Comment        string
	Unit           quantity.Quantity
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
		indexVars: Integer array with index of variables on equation
		vars: pointers's slice to ALL variables
	*/
	Index     uint16
	Line      uint16
	Text      string
	Exec      *vm.Program
	Env       map[string]interface{}
	IndexVars []int
	Vars      []*Variable
}

func (e *Equation) UpdateEnv(guessUpdate bool) {
	/*
		Update guesses from variables in the environment
	*/
	for _, v := range e.Vars {
		if guessUpdate {
			// Equations system are evaluated in SI
			valueSI := v.Guess * v.Unit.ToSI().Value()
			v.TempValue = valueSI
			e.Env[v.Name] = valueSI
		} else if !v.Solved {
			// Value from optimize iteration
			e.Env[v.Name] = v.TempValue
		} else {
			e.Env[v.Name] = v.Guess
		}
	}
}

func (e *Equation) RunProgram() (float64, error) {
	/*
		Run the program to evaluate the equation
	*/
	e.UpdateEnv(false)
	output, err := expr.Run(e.Exec, e.Env)
	if err != nil {
		log.Printf("%v", err)
		return math.NaN(), err
	} else {
		return output.(float64), nil
	}
}

type BlockEquations struct {
	/*
		Block from Tarjan algorithm
		Equations: Array of equation structs on block
		Variables: Array of variables ONLY of block
			equations
		Index: Block index of equations system
		Solved: When all block eqns are solved
	*/
	Equations []*Equation
	Variables []*Variable
	Index     int
	Solved    bool
	Result    optimize.Result
}

type Solution struct {
	/*
		Solution struct for imhotep

	*/
	Variables []Variable
	Equations []Equation
	Blocks    []BlockEquations
}

type LowestVar struct {
	/*
		Stores the information about a node relationship in the
		adjacency matrix to the representation of the equations system

		Val: 1 or 0 ; variable exist or not in the equations
		Row: index of the row (equation) in the matrix
		Col: index of the col (variable) in the matrix
	*/
	Val, Row, Col int
}
