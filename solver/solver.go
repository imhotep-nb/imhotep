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
	fmt.Printf("settingsSolver: %v\n", settingsSolver)
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
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result.Status: %v\n", result.Status)
	fmt.Printf("result.X: %0.4g\n", result.X)
	fmt.Printf("result.F: %0.4g\n", result.F)
	fmt.Printf("Time: %v microseconds\n", result.Runtime.Microseconds())
	fmt.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)
}
