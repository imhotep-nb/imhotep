package parsers

import (
	"encoding/json"
	"errors"
	"imhotep/constructors"
	"imhotep/types"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/imhotep-nb/units/quantity"
)

func ParseText(File string, Vars *[]*types.Variable,
	Eqns *[]*types.Equation, Settings *types.SolverSettings) (bool, error) {
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
	concatEqnsParsed, err := parseExplicitUnits(concatEqns)
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
	varsS := []string{}
	// Extract the variables, for each line
	tempLine := ""
	// varsD: Variables to replace to ""
	varsD := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	// varsR: Variables to replace to "-", for
	// variables extraction
	varsR := []string{"(", "+", "*", "-", ")", ".", "/", "e", "phi", "pi"}
	lineCopy := make([]types.EquationJSON, nEqns)
	copy(lineCopy, input.Equations)
	varsDirectory := make(map[int][]int)
	for i, line := range input.Equations {
		varsDirectory[i] = []int{}
		// Clear spaces
		if line.UnitsParsedText != "" {

			tempLine = lineCopy[i].UnitsParsedText

			// Remove functions
			for _, varFunc := range funcStrings {
				tempLine = strings.ReplaceAll(tempLine, varFunc, "")
			}
			// Remove numbers and spaces
			for _, varD := range varsD {
				tempLine = strings.ReplaceAll(tempLine, varD, "")
				line.UnitsParsedText = strings.ReplaceAll(line.UnitsParsedText, varD, "")
			}
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
					indexT := -1
					for index, varS := range varsS {
						if varS == varT {
							liveInside = true
							indexT = index
						}
					}
					// If not, add it!
					if !liveInside {
						varsS = append(varsS, varT)
						varsDirectory[i] = append(varsDirectory[i], len(varsS)-1)
					} else {
						varsDirectory[i] = append(varsDirectory[i], indexT)
					}
				}
			}
		}
	}

	if len(varsS) != nEqns {
		err = errors.New("mismatch number equations with variables")
		// TODO give variables and equations numbers
		log.Printf("%v", err)
		log.Printf("Variables: %v", varsS)
		log.Printf("Equations Number: %v", nEqns)
		return false, err
	}

	// For for generate variable structs:
	for i, varJSON := range input.Variables {

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
		log.Printf("The variable %v has index %v\n", varJSON.Name, i)
		newVar, err := constructors.NewVariable(varJSON.Name, uint16(i), &guess,
			&upperlim, &lowerlim, &comment, &unit)
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
		if line.UnitsParsedText != "" {
			log.Printf("The equation is: %v ", line.UnitsParsedText)
			log.Printf("and the indices are %v\n", varsDirectory[i])
			newEq, err2 := constructors.NewEquation(line.UnitsParsedText, *Vars,
				uint16(i), uint16(line.Line), varsDirectory[i])
			if err2 != nil {
				log.Printf("%v", err2)
				return false, err2
			}
			*Eqns = append(*Eqns, newEq)
		}
	}
	return true, nil
}

func parseExplicitUnits(eqnText string) (string, error) {
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
