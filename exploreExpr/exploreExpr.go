package main

import (
	"fmt"
	"math"

	"github.com/antonmedv/expr"
)

type Env struct {
	x float64
}

func main() {
	code := `5*x-x + sin(x)`
	env := map[string]interface{}{
		"x":   3.0,
		"cos": math.Cos,
		"tan": math.Tan,
		"sin": math.Sin,
	}

	program, err := expr.Compile(code, expr.Env(env), expr.AsFloat64())
	if err != nil {
		panic(err)
	}
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output.(float64) + 4)
}
