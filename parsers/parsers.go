package parsers

import (
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
	buf, err := ioutil.ReadFile(File)
	if err != nil {
		log.Printf("%v", err)
		return false, err
	}
	// s it's the string with the contents of file
	s := string(buf)
	// Separete in lines
	s = strings.ReplaceAll(s, "\r", "")
	v := strings.Split(s, "\n")
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
	lineCopy := make([]string, len(v))
	copy(lineCopy, v)
	for i, line := range v {
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
	// For for generate variable structs:
	for _, varS := range varsS {
		guess := 0.0
		lowerlim := -100.0
		upperlim := 100.0
		comment := ""
		newVar, err := constructors.NewVariable(varS, &guess, &upperlim, &lowerlim, &comment, nil)
		if err != nil {
			log.Printf("Error in variable creation %v: %v", varS, err)
			// &Vars = *[]*types.Variable{}
			// &Eqns = *[]*types.Equation{}
			return false, nil
		}
		*Vars = append(*Vars, newVar)
	}
	// For every line,create a equation
	for i, line := range v {
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
