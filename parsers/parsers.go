package parsers

import (
	"encoding/json"
	"errors"
	"imhotep/constructors"
	"imhotep/types"
	"io/ioutil"
	"log"
	"strings"
)

func ParseText(File string, Vars *[]*types.Variable,
	Eqns *[]*types.Equation) (bool, error) {
	/*
	   This function parse a file text string to a
	*/
	var input types.APIInput
	buf, err := ioutil.ReadFile(File)

	if err != nil {
		log.Printf("%v", err)
		return false, err
	}

	json.Unmarshal(buf, &input)

	// Reset Vars and Equations
	*Vars = []*types.Variable{}
	*Eqns = []*types.Equation{}
	varsS := []string{}
	// Extract the variables, for each line
	tempLine := ""
	// varsD: Variables to replace to ""
	varsD := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	// varsR: Variables to replace to "-", for
	// variables extraction
	varsR := []string{"(", "+", "*", "-", ")"}
	lineCopy := make([]string, len(input.Equations))
	copy(lineCopy, input.Equations)
	for i, line := range input.Equations {
		// Clear spaces
		if line != "" {
			tempLine = lineCopy[i]
			for _, varD := range varsD {
				tempLine = strings.ReplaceAll(tempLine, varD, "")
				line = strings.ReplaceAll(line, varD, "")
			}
			// Put "-"
			for _, varR := range varsR {
				tempLine = strings.ReplaceAll(tempLine, varR, "&")
			}
			// Split "-"
			log.Printf("The string for splitting is: %v", tempLine)
			varsT := strings.Split(tempLine, "&")
			log.Printf("Load a variables list in string format: %v. Length of: %v", varsT, len(varsT))
			// Evaluate availability
			for _, varT := range varsT {
				if varT != " " && varT != "" {
					// Check if varT lives in varsS
					liveInside := false
					for _, varS := range varsS {
						if varS == varT {
							liveInside = true
						}
					}
					// If not, add it!
					if !liveInside {
						varsS = append(varsS, varT)
					}
				}
			}
		}
	}

	if len(varsS) != len(input.Equations) {
		err = errors.New("mismatch number equations with variables")
		// TODO give variables and equations numbers
		log.Printf("%v", err)
		return false, err
	}

	// For for generate variable structs:
	for _, varJSON := range input.Variables {

		// validate if the varJSON exits in the identified variables in the equations (varsS)
		var existVar = false
		for _, varS := range varsS {
			if varS == varJSON.Name {
				existVar = true
				break
			}
		}

		if !existVar {
			err = errors.New("variables in JSON mismatch with variables in equations")
			// TODO give variables names from JSON and from equations
			log.Printf("%v", err)
			return false, err
		}

		guess := varJSON.Guess
		lowerlim := varJSON.Lowerlim
		upperlim := varJSON.Upperlim
		comment := varJSON.Comment
		unit := varJSON.Unit
		newVar, err := constructors.NewVariable(varJSON.Name, &guess, &upperlim, &lowerlim, &comment, &unit)
		if err != nil {
			log.Printf("Error in variable creation %v: %v", varJSON.Name, err)
			// &Vars = *[]*types.Variable{}
			// &Eqns = *[]*types.Equation{}
			return false, nil
		}
		*Vars = append(*Vars, newVar)
	}
	// For every line,create a equation
	for i, line := range input.Equations {
		if line != "" {
			log.Printf("The equation is: %v", line)
			newEq, err2 := constructors.NewEquation(line, *Vars, uint16(i), uint16(i))
			if err2 != nil {
				log.Printf("%v", err2)
				return false, err2
			}
			*Eqns = append(*Eqns, newEq)
		}
	}
	return true, nil
}
