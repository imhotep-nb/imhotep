package solver

import (
	"fmt"
	"imhotep/types"
	"log"
	"math"
	"time"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

func Hola() {
	fmt.Printf("Hola Imhoteeeep")
}

func SolverBlock(blockEqn types.BlockEquations,
	settingsSolver types.SolverSettings) (*optimize.Result, error) {

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
		for i, varBlock := range blockEqn.Variables {
			varBlock.TempValue = input[i]
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
		X = append(X, varB.Guess)
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

		for i, varBlock := range blockEqn.Variables {
			varBlock.Solved = true
			varBlock.Guess = result.X[i]
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
					// update when start the row (minsumCols only initializate) or when
					// identify a col with the lowest acummulative sum
					if minSumCols.Val == 0 || sumCols[j] < minSumCols.Val {
						minSumCols.Val = sumCols[j]
						minSumCols.Row = i
						minSumCols.Col = j
					}
				}
			}

			// update when end the row (minSumRows only initializate) or when
			// identify a row with the lowest acummulative sum
			if minSumRows.Val == 0 || sumRows[i] < minSumRows.Val {
				minSumRows.Val = sumRows[i]
				minSumRows.Row = i
				minSumRows.Col = lastJcol
			}

			// Conditions to determinate which pairing of eqn and variable assign
			minToAdd = types.LowestVar{}
			// If only there is a variable, that row (eqn) must be assign to that variable
			if minSumRows.Val == 1 {
				minToAdd.Val = 1
				minToAdd.Row = i
				minToAdd.Col = minSumRows.Col

				// At the end of current row,if a variable exist only in a equation,
				//it must be assign to that eqn
			} else if i == currentLastRow && minSumCols.Val == 1 {
				minToAdd.Val = 1
				minToAdd.Row = minSumCols.Row
				minToAdd.Col = minSumCols.Col

			} else if i == currentLastRow && minSumCols.Val < minSumRows.Val {
				minToAdd.Val = 1
				minToAdd.Row = minSumCols.Row
				minToAdd.Col = minSumCols.Col
			} else if i == currentLastRow && minSumRows.Val < minSumCols.Val {
				minToAdd.Val = 1
				minToAdd.Row = minSumRows.Row
				minToAdd.Col = minSumRows.Col
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
	}
	fmt.Printf("Convert full pseudograph take %v for loops\n", contadorFors)

	return pseudograph, reOrderEqn, nil
}
