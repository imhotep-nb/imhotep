package main

import (
	"imhotep/constructors"
	"imhotep/parsers"
	"imhotep/solver"
	"imhotep/types"
	"log"

	"github.com/looplab/tarjan"
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
	blocksEquation := []*types.BlockEquations{}
	newBlock, errB := constructors.NewBlockEquation(Eqns, Vars, 0)
	if errB != nil {
		log.Printf("Block determination fails: %v\n", errB)
	}
	blocksEquation = append(blocksEquation, newBlock)
	log.Printf("This block list is: %v\n", &blocksEquation)

	result, errS := solver.SolverBlock(*(blocksEquation)[0], Settings)

	if errS != nil {
		log.Printf("Block fails: %v", err)
	} else {
		log.Printf("result.Status: %v\n", result.Status)
		log.Printf("result.X: %0.4g\n", result.X)
		log.Printf("result.F: %0.4g\n", result.F)
		log.Printf("Time: %v microseconds\n", result.Runtime.Microseconds())
		log.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)

		for _, varS := range Vars {
			log.Printf("[%v] %v = %v\n", varS.Solved, varS.Name, varS.Guess)
		}
	}

	// Create the adjacency Matrix with rows and cols representing equations and variables respectively
	adjacencyMatrix := make([][]int, len(Eqns))
	for i, eqn := range Eqns {
		adjacencyMatrix[i] = make([]int, len(Eqns))

		for _, val := range eqn.IndexVars {
			adjacencyMatrix[i][val] = 1
		}
	}

	// Convert adjacency Matrix into pseudographe
	_, reOrderEqn, _ := solver.ConvertFullPseudograph(adjacencyMatrix)
	log.Printf("Re order de equations: %v\n", reOrderEqn)

	// Create the graph struct to the input of tarjan's package
	graph := make(map[interface{}][]interface{})
	for i, eqn := range Eqns {
		elements := make([]interface{}, len(eqn.IndexVars))

		for j, val := range eqn.IndexVars {
			elements[j] = val
		}
		graph[reOrderEqn[i]] = elements
	}

	log.Printf("Graph to tarjan: %v\n", graph)
	output := tarjan.Connections(graph)
	log.Printf("Blocks equations: %v\n", output)

}
