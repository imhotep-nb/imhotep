package main

import (
	"fmt"
	"imhotep/constructors"
)

func main() {

	name := "✖️"
	guess := 500000.0
	upperlim := 1000.0
	lowerlim := -1000.0
	comment := "This pleases imhotep 😏"
	unit := "kg/m"

	varPyENL, err := constructors.NewVariable(name, &guess, &upperlim, &lowerlim, &comment, &unit)

	if err == nil {
		fmt.Printf("Primera variable con constructor: %v \n", *varPyENL)

		fmt.Println(varPyENL.Unit)
	}

}
