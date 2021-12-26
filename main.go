package main

import (
	"imhotep/parsers"
	"imhotep/solver"
	"imhotep/types"
	"log"
)

func main() {

	Vars := []*types.Variable{}
	Eqns := []*types.Equation{}
	Settings := types.SolverSettings{}
	_, err := parsers.ParseText("./testing/tarjan.json", &Vars, &Eqns, &Settings)
	if err != nil {
		log.Printf("Something fails: %v", err)
		return
	}

	log.Print("Desde acá perros")

	solution, errSol := solver.Solver(Vars, Eqns, Settings)
	if errSol != nil {
		log.Print(errSol.Error())
	}
	log.Print(solution)
}
