package main

import (
	"fmt"
	"imhotep/constructors"
	"imhotep/types"
	"log"
)

func main() {

	name := "x"
	guess := 5.0
	upperlim := 1000.0
	lowerlim := -1000.0
	comment := "This pleases imhotep 😏 y qué?"
	unit := "kg/m"

	varPyENL, _ := constructors.NewVariable(name, &guess, &upperlim, &lowerlim, &comment, &unit)

	name2 := "y"
	guess2 := 3.0
	upperlim2 := 1000.0
	lowerlim2 := -1000.0
	comment2 := "This pleases imhotep 😏 2?"
	unit2 := "kg/m"

	varPyENL2, _ := constructors.NewVariable(name2, &guess2, &upperlim2, &lowerlim2, &comment2, &unit2)

	// First equation of Imhotep
	vars := []*types.Variable{varPyENL, varPyENL2}
	firstEquation, err3 := constructors.NewEquation("5*x + y", vars, 0, 0)
	if err3 != nil {
		log.Printf("%v", err3)
		return
	}
	out, err := firstEquation.RunProgram()
	if err != nil {
		log.Printf("%v", err)
		return
	}
	fmt.Printf("%v", out)

	// Vamos a cambiar los valores!!
	varPyENL.Guess = 10
	varPyENL2.Guess = 2
	out, err = firstEquation.RunProgram()
	if err != nil {
		log.Printf("%v", err)
		return
	}
	fmt.Printf("%v", out)

}
