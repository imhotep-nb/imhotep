package main

import (
	"imhotep/constructors"
	"imhotep/parsers"
	"imhotep/solver"
	"imhotep/types"
	"log"
)

func main() {

	Vars := []*types.Variable{}
	Eqns := []*types.Equation{}
	Settings := types.SolverSettings{}
	_, err := parsers.ParseText("./testing/texto.json", &Vars, &Eqns, &Settings)
	if err != nil {
		log.Printf("Something fails: %v", err)
		return
	}
	blocksEquation := []*types.BlockEquations{}
	newBlock, errB := constructors.NewBlockEquation(Eqns, Vars, 0)
	if errB != nil {
		log.Printf("Block determination fails: %v\n", errB)
	}
	blocksEquation = append(blocksEquation, newBlock)
	log.Printf("This block list is: %v\n", blocksEquation)
	solver.Solver(*(blocksEquation)[0], Settings)
	// for _, eqn := range Eqns {
	// 	val, errExec := eqn.RunProgram()
	// 	if errExec != nil {
	// 		log.Printf("Falló la ejecución de la ecuación %v: %v", eqn.Text, errExec)
	// 	}
	// 	fmt.Println(val)
	// }
}
