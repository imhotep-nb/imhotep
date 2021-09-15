package main

import (
	"fmt"
	"imhotep/parsers"
	"imhotep/types"
	"log"
)

func main() {

	Vars := []*types.Variable{}
	Eqns := []*types.Equation{}
	_, err := parsers.ParseText("/home/devgtc/texto", &Vars, &Eqns)
	if err != nil {
		log.Printf("Something fails: %v", err)
		return
	}
	for _, eqn := range Eqns {
		val, errExec := eqn.RunProgram()
		if errExec != nil {
			log.Printf("Falló la ejecución de la ecuación %v: %v", eqn.Text, errExec)
		}
		fmt.Println(val)
	}
}
