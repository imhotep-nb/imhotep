package parsers

import (
	"errors"
	"imhotep/constructors"
	"imhotep/types"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/imhotep-nb/units/quantity"
)

func ParseText(input types.APIInput, Vars *[]*types.Variable,
	Eqns *[]*types.Equation, Settings *types.SolverSettings,
	onlyVars bool) (bool, error) {
	/*
	   This function parse a file text string to a eqns and vars structs
	   If no error, the boolean param is the APIInput.Debug setting
	*/

	log.Printf("Debug: %v", input.Debug)

	*Settings = input.Settings

	// It need replace explicit units with conversion factors in SI
	// so concatenate eqn to replace the whole units at the same time.
	// Use a equations separator to concat it; the only requeriment is that
	// there is no risk that the separator could be in the eqn.
	concatEqns := ""
	eqnSeparator := "@@&##&@@##"
	nEqns := len(input.Equations)
	for i, line := range input.Equations {

		if i < (nEqns - 1) {
			concatEqns += line.Text + eqnSeparator
		} else if i == nEqns-1 {
			concatEqns += line.Text
		}

	}

	// Replace explicit units with conversion factors to SI
	concatEqnsParsed, err := ParseExplicitUnits(concatEqns)
	var unitsParsedTexts []string
	var funcStrings []string
	if err != nil {
		log.Printf("Can't parse units: %v\n", err)
		return false, err
	} else {
		// Find functions to avoid it when identify variables
		// Take avantage that already are concatenated all lines
		funcRegex := regexp.MustCompile(`[a-zA-Z0-9]+[\(]`)
		funcStrings = funcRegex.FindAllString(concatEqnsParsed, -1)
		unitsParsedTexts = strings.Split(concatEqnsParsed, eqnSeparator)
		log.Printf("Functions founded: %v\n", funcStrings)
	}

	for i, line := range unitsParsedTexts {
		input.Equations[i].UnitsParsedText = line
	}

	// Reset Vars and Equations
	*Vars = []*types.Variable{}
	*Eqns = []*types.Equation{}

	// Identify vars text in equations
	varsS, varsInEqns, errVars := IdentifyVars(input.Equations, funcStrings)
	if errVars != nil {
		log.Printf("Fail to identify vars: %v\n", errVars)
	}

	if len(varsS) != nEqns {
		err := errors.New("mismatch number equations with variables")
		// TODO give variables and equations numbers
		log.Printf("%v", err)
		log.Printf("Variables: %v", varsS)
		log.Printf("Equations Number: %v", nEqns)
		return false, err
	}

	varsIndexByName, errVars := MakeVars(varsS, Vars, input.Variables, onlyVars)
	if errVars != nil {
		log.Printf("Fail to created vars: %v\n", errVars)
	} else if onlyVars {
		return true, nil
	}

	// For every line,create a equation
	for i, line := range input.Equations {
		if line.UnitsParsedText != "" {
			log.Printf("The equation is: %v ", line.UnitsParsedText)

			// parse the variable Name with the variable index
			varsIndexInEqn_i := make([]int, len(varsInEqns[i])) // index of each variable
			for j, nameVar := range varsInEqns[i] {
				varsIndexInEqn_i[j] = varsIndexByName[nameVar]
			}

			log.Printf("and the indices are %v\n", varsIndexInEqn_i)
			newEq, err2 := constructors.NewEquation(line.UnitsParsedText, *Vars,
				uint16(i), uint16(line.Line), varsIndexInEqn_i)
			if err2 != nil {
				log.Printf("%v", err2)
				return false, err2
			}
			*Eqns = append(*Eqns, newEq)
		}
	}
	return input.Debug, nil
}

func ParseExplicitUnits(eqnText string) (string, error) {
	/*
	   This function parse a string (like a equation string)  using a regex to identify
	   all explicit units, after that use quantity symbol parse to convert the unit to SI
	   and finally replace each explicit unit in the text for the conversión factor to SI

	   Example:

	   input: "x + y - 10[ft] = 50[ft]"
	   output: "x + y - 10*0.304800 = 50*0.304800
	*/
	newEqnText := eqnText

	unitsRegex := regexp.MustCompile(`\[(.*?)\]`)
	unitsInEqn := unitsRegex.FindAllString(eqnText, -1)

	unitsUniques := make(map[string]bool)
	for _, unitText := range unitsInEqn {

		// Check that is a new unit
		if !unitsUniques[unitText] {

			// Remove square brackets []
			unitClear := strings.ReplaceAll(unitText, "[", "")
			unitClear = strings.ReplaceAll(unitClear, "]", "")

			// Parse unit and check that is a valit one.
			unitIdentify, err := quantity.ParseSymbol(unitClear)

			if err != nil {
				log.Printf("The unit %v can't be identify: %v\n", unitText, err)
				return "", err

			} else {

				// Get factor conversion in string format to replace it in the eqns.
				valueSI := unitIdentify.ToSI().Value()
				converFactor := strconv.FormatFloat(valueSI, 'e', -1, 64)
				newEqnText = strings.ReplaceAll(newEqnText, unitText, "*"+converFactor)

				// Save that unit to avoid repeat the process
				unitsUniques[unitText] = true
			}
		}
	}

	log.Printf("Equation with units parsed: %v", newEqnText)

	return newEqnText, nil

}

func IdentifyVars(inputEquations []types.EquationJSON, funcStrings []string) ([]string, map[int][]string, error) {
	/*
	   This function identify variables in the inputEquations and return it
	   like an array of strings, where each item is a variable found. Also
	   return the vars found in each equations.
	*/
	varsS := []string{}
	nEqns := len(inputEquations)

	// Extract the variables, for each line
	tempLine := ""
	// varsR: Variables to replace to "-", for
	// variables extraction
	varsR := []string{"(", "+", "*", "-", ")", ".", "/", "e", "phi", "pi", ","}
	lineCopy := make([]types.EquationJSON, nEqns)
	// Backup of equation json objects
	copy(lineCopy, inputEquations)
	// This is for replace string variables.
	// Example:
	// h = entropy('Water', P, T, 20)
	// where 'Water' is a string variable
	var reStringVars = regexp.MustCompile(`\'\w+'`)
	// This regular expression is for looking numbers wich no be
	// part of some equation name
	var reAloneNumbers = regexp.MustCompile(`(^[0-9][\*\/\+\-])|([\s\+\-\*\/\(][0-9]{1,})*`)
	// This slice will store the vars on string format
	varsInEqns := make(map[int][]string)
	for i, line := range inputEquations {
		varsInEqns[i] = []string{}
		// Clear spaces
		if line.UnitsParsedText != "" {
			// tempLine store the equation line, but will be modified
			tempLine = lineCopy[i].UnitsParsedText
			// Remove strings vars
			tempLine = reStringVars.ReplaceAllString(tempLine, "")
			// Remove functions
			for _, varFunc := range funcStrings {
				tempLine = strings.ReplaceAll(tempLine, varFunc, "")
			}
			// Remove numbers
			tempLine = reAloneNumbers.ReplaceAllString(tempLine, "")
			// TODO: Remove spaces
			tempLine = strings.ReplaceAll(tempLine, " ", "")
			// Replace operator symbols with  "&"
			for _, varR := range varsR {
				tempLine = strings.ReplaceAll(tempLine, varR, "&")
			}
			// Split "&", but use FieldsFunc instead, to avoid empty values in slice
			log.Printf("The string for splitting is: %v", tempLine)
			f := func(c rune) bool {
				return c == '&'
			}
			varsT := strings.FieldsFunc(tempLine, f)
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
					// add var in the row of the equation i
					varsInEqns[i] = append(varsInEqns[i], varT)
				}
			}
		}
	}

	return varsS, varsInEqns, nil

}

func MakeVars(varsS []string, Vars *[]*types.Variable,
	inputVariables []types.VariableJSON, onlyVars bool) (map[string]int, error) {
	/*
		This functions create the vars struct from the name vars string array
		identify on IdentifyVars().

		There are some cases:
		1. User input partially or totally the vars in JSON API input
		2. User input doesn't have any variable definition
		3. User input doesn't have any variable definition and only requeries
		   get variables
	*/
	varsIndexByName := make(map[string]int)

	if onlyVars {
		// Generate dummy Vars (only for names)
		for i, varsS_ := range varsS {
			newVar := types.Variable{Name: varsS_}
			*Vars = append(*Vars, &newVar)
			varsIndexByName[varsS_] = i
		}
	} else {
		// For generate variable structs:
		for i, varJSON := range inputVariables {

			// validate if the varJSON exits in the identified variables in the equations (varsS)
			var existVar = false
			for _, varS := range varsS {
				if varS == varJSON.Name {
					existVar = true
					break
				}
			}

			if !existVar {
				err := errors.New("variables in JSON mismatch with variables in equations")
				// TODO give variables names from JSON and from equations
				log.Printf("%v", err)
				log.Printf("Variables varsS: %v, Variables varJSON: %v\n", varsS, varJSON)
				return varsIndexByName, err
			}

			guess := varJSON.Guess
			lowerlim := varJSON.Lowerlim
			upperlim := varJSON.Upperlim
			comment := varJSON.Comment
			unit := varJSON.Unit
			log.Printf("The variable %v has index %v\n", varJSON.Name, i)
			newVar, err := constructors.NewVariable(varJSON.Name, uint16(i), &guess,
				&upperlim, &lowerlim, &comment, &unit)
			if err != nil {
				log.Printf("Error in variable creation %v: %v", varJSON.Name, err)
				// &Vars = *[]*types.Variable{}
				// &Eqns = *[]*types.Equation{}
				return varsIndexByName, err
			}
			*Vars = append(*Vars, newVar)
			// Assign the index to the variable name
			varsIndexByName[varJSON.Name] = i
		}
	}

	return varsIndexByName, nil
}
