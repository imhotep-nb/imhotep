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
	debug, err := parsers.ParseText("./testing/tarjan.json", &Vars, &Eqns, &Settings)
	if err != nil {
		log.Printf("Something fails: %v", err)
		return
	}

	solution, errSol := solver.Solver(Vars, Eqns, Settings, debug)
	if errSol != nil {
		log.Print(errSol.Error())
	}
	log.Print(solution)
}
