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

func Solver(blockEqn types.BlockEquations, settingsSolver types.SolverSettings) {
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
		log.Printf("%v\n", output)
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
		log.Printf("The solver fails: %v", err)
		return
	}
	log.Printf("Finally!: %v", result)
}
