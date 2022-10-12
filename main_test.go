package main

import (
	"encoding/json"
	"imhotep/types"
	"io/ioutil"
	"math"
	"testing"
)

func TestSolveProblem(t *testing.T) {

	file := "./testing/eqn-system-to-test.json"
	buf, err := ioutil.ReadFile(file)

	if err != nil {
		t.Logf("Fails to read JSON testing eqns: %v", err)
		t.Fail()
	}

	var input types.TestInputEquations
	json.Unmarshal(buf, &input)

	// Necesary to avoid false positive if the input is empty
	if input == nil {
		t.Log("Can't parse test file to JSON")
		t.Fail()
	}

	for i, example := range input {

		apiInput, errGen := APITestingInputGenerator(&example)

		if errGen != nil {
			t.Logf("Cant generate APIInput to example N° %v: %v", i, errGen)
			t.Fail()
		}

		out, errSol := solveProblem(apiInput, false)

		if errSol != nil {
			t.Logf("Fail to solve problem with example N° %v: %v", i, errSol)
			t.Fail()
		}

		// file, errJSON := json.Marshal(out)

		// if errJSON != nil {
		// 	t.Logf("Fail to create JSON %v: %v", i, errJSON)
		// 	t.Fail()
		// }

		tolerance := out.Settings.GradientThreshold
		for _, varOut := range out.Vars {

			for _, varRef := range example.VariablesSolved {
				if varOut.Name == varRef.Name {

					diff := math.Abs(varRef.Guess - varOut.Guess)
					if diff > tolerance {
						t.Logf("Fail to solve var %v ,Reference: %v , Calculated: %v, Diff: %v > GradientThreshold = %v", varRef.Name, varRef.Guess, varOut.Guess, diff, tolerance)
						t.Fail()
					}
				}
			}

		}
		// 		errWrite := ioutil.WriteFile("./testing/res.json", file, 0644)

		// 		if errWrite != nil {
		// 			t.Logf("Fail to write solve: %v", errWrite)
		// 			t.Fail()
		// 		}
	}

}

func APITestingInputGenerator(example *struct {
	Name            string               "json:\"name\""
	Equations       []string             "json:\"eqns\""
	Variables       []types.VariableJSON "json:\"vars\""
	VariablesSolved []types.VariableJSON "json:\"varsSolved\""
}) (types.APIInput, error) {

	// Identify the variables equations and create equations array objetcs
	var eqns []types.EquationJSON
	for i, textEqn := range example.Equations {

		newEqn := types.EquationJSON{
			Text: textEqn,
			Line: i,
		}
		eqns = append(eqns, newEqn)
	}

	input := types.APIInput{
		Equations: eqns,
		Variables: example.Variables,
	}

	// create variables array objects

	return input, nil
}
