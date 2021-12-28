package solver

import (
	"errors"
	"imhotep/types"
	"log"
	"math"
	"time"

	"github.com/looplab/tarjan"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

func Solver(Vars []*types.Variable, Eqns []*types.Equation,
	settingsSolver types.SolverSettings, debug bool) (types.APIOutput, error) {
	/*
		Split the equations into groups of small equations
		called blocks, that could be solved in order saving
		time of execution using Tarjan algorithm

		McNunn, G. S. (2013). Using Tarjan's algorithm to organize and schedule the
		computational workflow in a federated system of models and databases.
		[https://web.archive.org/web/20201021194957/https://lib.dr.iastate.edu/etd/13566/]
	*/

	adjacencyMatrix, errAdj := MakeAdjacencyMatrix(Eqns)
	if errAdj != nil {
		log.Print(errAdj.Error())
		return types.APIOutput{}, errAdj
	}
	blocksEqnIndex, blocksEqnIndexInv, _, errBlocks := MakeEquationBlocks(Eqns, adjacencyMatrix)
	if errBlocks != nil {
		log.Print(errBlocks.Error())
		return types.APIOutput{}, errBlocks
	}
	log.Printf("Los bloques: %v", blocksEqnIndex)

	// Create block equations from types
	blocks := make([]types.BlockEquations, len(blocksEqnIndex))
	for i, blockIndexes := range blocksEqnIndex {
		eqnList := make([]*types.Equation, len(blockIndexes))
		indexVars := []int{}
		varsList := []*types.Variable{}
		for j, EqnIndex := range blockIndexes {
			eqnList[j] = Eqns[blocksEqnIndexInv[EqnIndex.(int)]]
			for _, VarIndex := range eqnList[j].IndexVars {
				add := true
				for _, Var := range indexVars {
					if Var == VarIndex {
						add = false
					}
				}
				if add {
					indexVars = append(indexVars, VarIndex)
					varsList = append(varsList, Vars[VarIndex])
				}
			}
		}
		block := types.BlockEquations{
			Equations: eqnList,
			Variables: varsList,
			Index:     i,
			Solved:    false,
		}
		blocks[i] = block
	}

	// In this section, the blocks will be solved

	for i, block := range blocks {
		result, errS := SolverBlock(block, settingsSolver)
		log.Print("------------------------------------------------------------------")
		log.Printf("Este es el bloque número %v", i)
		if errS != nil {
			log.Printf("Block fails: %v", errS)
			log.Printf("The equations: %v", *block.Equations[0])
			return types.APIOutput{}, errS
		} else {
			blocks[i].Result = optimize.Result{
				Stats:  result.Stats,
				Status: result.Status,
			}
			log.Printf("result.Status: %v\n", result.Status)
			log.Printf("result.X: %0.4g\n", result.X)
			log.Printf("result.F: %0.4g\n", result.F)
			log.Printf("Time: %v microseconds\n", result.Runtime.Microseconds())
			log.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)

			for _, varS := range Vars {
				log.Printf("[%v] %v = %v\n", varS.Solved, varS.Name, varS.Guess)
			}
		}
	}

	// Output Structure
	output := types.APIOutput{}
	if debug {
		equations := make([]types.EquationJSON, len(Eqns))
		for i, eqn := range Eqns {
			equations[i] = types.EquationJSON{
				Text: eqn.Text,
				Line: int(eqn.Line),
			}
		}
		output.Eqns = equations
	}
	variables := make([]types.VariableJSON, len(Vars))
	for i, variable := range Vars {
		variables[i] = types.VariableJSON{
			Name:     variable.Name,
			Guess:    variable.Guess,
			Upperlim: variable.Upperlim,
			Lowerlim: variable.Lowerlim,
			Comment:  variable.Comment,
			Unit:     variable.Unit.String(),
		}
	}
	output.Vars = variables
	output.Settings = settingsSolver

	stats := types.Stats{
		Blocks: make([]types.BlockInfo, len(blocks)),
		Global: optimize.Stats{},
	}
	for i, block := range blocks {
		stats.Blocks[i] = types.BlockInfo{
			Stats:  block.Result.Stats,
			Status: block.Result.Status.String(),
		}
		stats.Global.FuncEvaluations += block.Result.FuncEvaluations
		stats.Global.GradEvaluations += block.Result.GradEvaluations
		stats.Global.MajorIterations += block.Result.MajorIterations
		stats.Global.Runtime += block.Result.Runtime
	}
	output.Stats = stats

	// output.Info.Graph = graph

	return output, nil
}

func MakeAdjacencyMatrix(Eqns []*types.Equation) ([][]int, error) {
	// Create the adjacency Matrix with rows and cols representing equations and variables respectively
	adjacencyMatrix := make([][]int, len(Eqns))
	for i, eqn := range Eqns {
		adjacencyMatrix[i] = make([]int, len(Eqns))

		for _, val := range eqn.IndexVars {
			adjacencyMatrix[i][val] = 1
		}
	}
	return adjacencyMatrix, nil
}

func MakeEquationBlocks(Eqns []*types.Equation, adjacencyMatrix [][]int) ([][]interface{}, map[int]int, map[interface{}][]interface{}, error) {
	// Convert adjacency Matrix into pseudographe
	_, reOrderEqn, errM := ConvertFullPseudograph(adjacencyMatrix)

	if errM != nil {
		log.Print(errM.Error())
		return nil, nil, nil, errM
	}
	log.Printf("Re order de equations: %v\n", reOrderEqn)
	// Create the graph struct to the input of tarjan's package
	graph := make(map[interface{}][]interface{})
	reOrderEqnInv := make(map[int]int)
	for i, eqn := range Eqns {
		elements := make([]interface{}, len(eqn.IndexVars))

		for j, val := range eqn.IndexVars {
			elements[j] = val
		}
		graph[reOrderEqn[i]] = elements
		reOrderEqnInv[reOrderEqn[i]] = i
	}

	log.Printf("Graph to tarjan: %v\n", graph)
	output := tarjan.Connections(graph)
	log.Printf("Blocks equations: %v\n", output)
	return output, reOrderEqnInv, graph, nil
}

func SolverBlock(blockEqn types.BlockEquations,
	settingsSolver types.SolverSettings) (*optimize.Result, error) {
	/*
		Solve one block from system equations using optimize.Minimize
		and the method specified for user in settings
	*/

	settings := optimize.Settings{
		GradientThreshold: settingsSolver.GradientThreshold,
		MajorIterations:   settingsSolver.MajorIterations,
		Runtime:           time.Duration(settingsSolver.Runtime),
		FuncEvaluations:   settingsSolver.FuncEvaluations,
		GradEvaluations:   settingsSolver.GradEvaluations,
		Concurrent:        settingsSolver.Concurrent,
	}

	objective := func(input []float64) float64 {
		output := 0.0
		i := 0
		for _, varBlock := range blockEqn.Variables {
			if !varBlock.Solved {
				varBlock.TempValue = input[i]
				i += 1
			}
		}
		for _, eqnBlock := range blockEqn.Equations {
			out, err := eqnBlock.RunProgram()
			if err != nil {
				// TODO: Handle errors here
				log.Printf("Error on equation %v : %v\n",
					eqnBlock.Text, err)
				return math.NaN()
			}
			output += out * out
		}
		return output
	}

	gradObjective := func(grad, x []float64) {
		fd.Gradient(grad, objective, x, nil)
	}

	problem := optimize.Problem{
		Func: objective,
		Grad: gradObjective,
	}

	// Initialize X with guesses for optimize iteration
	X := []float64{}
	for _, varB := range blockEqn.Variables {
		if !varB.Solved {
			X = append(X, varB.Guess)
		}
	}
	result, err := optimize.Minimize(problem, X, &settings, nil)
	if err != nil {
		log.Printf("Minimize fails: %v \n", err)
		return nil, err

	} else if err = result.Status.Err(); err != nil {
		log.Printf("Status fails: %v \n", err)
		return nil, err

	} else {
		blockEqn.Solved = true

		i := 0
		for _, varBlock := range blockEqn.Variables {
			if !varBlock.Solved {
				varBlock.Solved = true
				varBlock.Guess = result.X[i]
				i += 1
			}
		}

		return result, nil
	}
}

func ConvertFullPseudograph(adjacencyMatrix [][]int) ([][]int, map[int]int, error) {
	/*
		Re order the adjacency matrix (the graphe representation of the equations system)
		in a full pseudograph if it is possible (here we call A FULL pseudograp, that one
		where all nodes must have one loop)

		In this way, we relation a variable with a equations, nevertless to say, the
		variable must exist in the equation, thats why we say that all nodes must have one
		loop. The new adjacency Matrix will have a diagonal full with 1.

		Not all adjacency matrix can convert into a full pseudograph. If can't achieve it,
		we can assure that the equations system is inconsistent, and it will fire an error.

		The ROWS of the adjacency matrix represent the VARIABLES
		The COLS of the adjancecy matrix represent the EQUATIONS

		Example:

			1) x + y - z = 0
			2) y - z = 10
			3) z = 5

			Adjacency matrix to the 3x3 equations system:
				x	y	z
			1)	1	1	1
			2)	0	1	1
			3)  0	0	1
	*/

	n := len(adjacencyMatrix)
	reOrderEqn := make(map[int]int)
	pseudograph := make([][]int, n)

	// To ignore the eqns that already are ordered and the vars assigned
	varsAlreadyAssign := make(map[int]bool)
	eqnsAlreadyAssign := make(map[int]bool)
	currentLastRow := n - 1

	var contadorFors int // only to know the number of loops

	// In each iteration we order a row (equation)
	// that means iter n times
	for iter := 0; iter < n; iter++ {
		contadorFors += 1
		sumCols := make([]int, n)
		sumRows := make([]int, n)
		minSumRows := types.LowestVar{}
		minSumCols := types.LowestVar{}
		minToAdd := types.LowestVar{}
		var lastJcol int
		lastiRow := make(map[int]int)

		// Go through rows
		for i, eqn := range adjacencyMatrix {

			if eqnsAlreadyAssign[i] {
				continue
			}

			contadorFors += 1
			for j, existVar := range eqn {

				if varsAlreadyAssign[j] {
					continue
				}

				contadorFors += 1
				if existVar == 1 {
					sumCols[j] += 1
					sumRows[i] += 1
					lastJcol = j
					lastiRow[j] = i
				}

				// check that any acumulative sum is zero at the end
				if i == currentLastRow && sumCols[j] == 0 {
					err := errors.New("fails to reorder the adjacency matrix to assign a eqn to a variable")
					log.Printf(" index variable %v can't be assign to any eqn: %v\n", i, err)
					return [][]int{}, map[int]int{}, err
				}

				// update only in the last row, when ensures that the acumulative sum in each
				// column is the total variables sum
				if i == currentLastRow && (minSumCols.Val == 0 || sumCols[j] < minSumCols.Val) {
					minSumCols.Val = sumCols[j]
					minSumCols.Row = lastiRow[j]
					minSumCols.Col = j
				}
			}

			// check that any acumulative sum is zero at the end
			if i == currentLastRow && sumRows[i] == 0 {
				err := errors.New("fails to reorder the adjacency matrix to assign a eqn to a variable")
				log.Printf(" index eqn %v can't be assign to any variable: %v\n", i, err)
				return [][]int{}, map[int]int{}, err
			}

			// update when end the row (minSumRows only initializate) or when
			// identify a row with the lowest acummulative sum
			if minSumRows.Val == 0 || sumRows[i] < minSumRows.Val {
				minSumRows.Val = sumRows[i]
				minSumRows.Row = i
				minSumRows.Col = lastJcol
			}

			// The pair eqn with variable (row with col) is determined when finish
			// to sum all rows an cols, and it will select the lowest acumulative sum
			// either from rows or columns
			minToAdd = types.LowestVar{}
			if minSumRows.Val == 1 {
				// If only there is a variable, that row (eqn) can be assign to
				// that variable without wait until go through all rows
				minToAdd.Val = 1
				minToAdd.Row = i
				minToAdd.Col = minSumRows.Col

			} else if i == currentLastRow && minSumRows.Val <= minSumCols.Val {
				minToAdd.Val = 1
				minToAdd.Row = minSumRows.Row
				minToAdd.Col = minSumRows.Col
			} else if i == currentLastRow && minSumRows.Val > minSumCols.Val {
				minToAdd.Val = 1
				minToAdd.Row = minSumCols.Row
				minToAdd.Col = minSumCols.Col
			}

			// Save the pair eqn and var identify above
			if minToAdd.Val == 1 {

				reOrderEqn[minToAdd.Row] = minToAdd.Col
				// create the connection matrix of the full pseudograph
				pseudograph[minToAdd.Col] = adjacencyMatrix[minToAdd.Row]

				varsAlreadyAssign[minToAdd.Col] = true
				eqnsAlreadyAssign[minToAdd.Row] = true

				// Determinate the current last eqn
				for k := 0; k < n; k++ {

					if !eqnsAlreadyAssign[k] {
						contadorFors += 1
						currentLastRow = k
					}
				}
				break
			}
		}
		log.Printf("iter %v : %v (%v, %v) minRow and %v (%v, %v) minCol, pairing %v  -> %v",
			iter, minSumRows.Val, minSumRows.Row, minSumRows.Col, minSumCols.Val, minSumCols.Row,
			minSumCols.Col, minToAdd.Row, minToAdd.Col)
	}
	log.Printf("Convert full pseudograph take %v for loops\n", contadorFors)

	return pseudograph, reOrderEqn, nil
}
