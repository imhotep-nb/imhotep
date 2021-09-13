package constructors

import (
	"errors"
	"imhotep/types"
	"math"

	"log"

	"github.com/imhotep-nb/units/quantity"
)

func NewVariable(Name string, Guess *float64, Upperlim *float64, Lowerlim *float64,
	Comment *string, Unit *string) (*types.Variable, error) {

	/*
		Parsear los limites para que cuando el usuario ingrese Inf or -Inf sean nil
	*/
	var err error

	if Name == "" {

		err = errors.New("variables should have name")
		log.Printf("%v", err)
		return nil, err
	}

	newVar := types.Variable{Name: Name}

	if Upperlim == nil {
		newVar.Upperlim = math.Inf(1)
	} else {
		newVar.Upperlim = *Upperlim
	}

	if Lowerlim == nil {
		newVar.Lowerlim = math.Inf(-1)
	} else {
		newVar.Lowerlim = *Lowerlim
	}

	if Guess == nil {
		newVar.Guess = 0
	} else {
		newVar.Guess = *Guess
	}

	if !(newVar.Guess > newVar.Lowerlim && newVar.Guess < newVar.Upperlim) {

		if Guess == nil {
			// TODO Promediar los limites y cuando el promedio sea NAN logica especial

		} else {
			err = errors.New("guess is outside of limits")
			log.Printf("Variable %v malformed guess %v: %v", newVar.Name, newVar.Guess, err)
			return nil, err
		}
	}

	if Comment == nil {
		newVar.Comment = ""
	} else {
		newVar.Comment = *Comment
	}

	var tempUnit quantity.Quantity

	if Unit == nil {
		tempUnit, err = quantity.ParseSymbol("m/m")

	} else {
		tempUnit, err = quantity.ParseSymbol(*Unit)
	}

	if err != nil {
		log.Printf("The variable %v can't be assign the unit %v: %v\n", newVar.Name, *Unit, err)
		return nil, err

	} else {
		newVar.Unit = *tempUnit.Unit
		newVar.Dimensionality = tempUnit.Dimensionality()
	}

	if newVar.Upperlim <= newVar.Lowerlim {
		err = errors.New("upperlim should be greater than lowerlim")
		log.Printf("Upperlim %v is lower than or equal to Lowerlim %v : %v", *Upperlim, *Lowerlim, err)
		return nil, err
	}

	return &newVar, nil

}
