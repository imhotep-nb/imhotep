package main

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

func Objetivo(input []float64) float64 {
	// 5y - 4x^2 + 3x - 7 = 0
	// 3y + 15x^2 + 42x - 5777 = 0
	var x = input[0]
	var y = input[1]
	out := [2]float64{0, 0}
	out[0] = 5*y - 4*x*x + 3*x - 7
	out[1] = 3*y + 15*x*x + 42*x - 5777
	return out[0]*out[0] + out[1]*out[1]
}

func main() {
	// p := optimize.Problem{
	// 	Func: functions.ExtendedRosenbrock{}.Func,
	// 	Grad: functions.ExtendedRosenbrock{}.Grad,
	// }

	// x := []float64{1.3, 0.7, 0.8, 1.9, 1.2}
	// result, err := optimize.Minimize(p, x, nil, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = result.Status.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("result.Status: %v\n", result.Status)
	// fmt.Printf("result.X: %0.4g\n", result.X)
	// fmt.Printf("result.F: %0.4g\n", result.F)
	// fmt.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)

	gradObjetivo := func(grad, x []float64) {
		fd.Gradient(grad, Objetivo, x, nil)
	}

	p2 := optimize.Problem{
		Func: Objetivo,
		Grad: gradObjetivo,
	}

	x2 := []float64{10, 100}
	result2, err2 := optimize.Minimize(p2, x2, nil, nil)
	if err2 != nil {
		log.Fatal(err2)
	}
	if err2 = result2.Status.Err(); err2 != nil {
		log.Fatal(err2)
	}
	fmt.Printf("result.Status: %v\n", result2.Status)
	fmt.Printf("result.X: %0.4g\n", result2.X)
	fmt.Printf("result.F: %0.4g\n", result2.F)
	fmt.Printf("result.Stats.FuncEvaluations: %d\n", result2.Stats.FuncEvaluations)
}
