package main

import (
	"imhotep/parsers"
	"imhotep/solver"
	"imhotep/types"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/solve", func(c *fiber.Ctx) error {
		input := new(types.APIInput)
		err := c.BodyParser(input)
		if err != nil {
			out := types.APIOutput{
				Info: types.Info{
					Errors: []string{err.Error()},
				},
			}
			return c.JSON(out)
		}
		solution, errSol := solveProblem(*input)
		if errSol != nil {
			return c.JSON(solution)
		}
		return c.JSON(solution)
	})
	log.Fatal(app.Listen(":3000"))

}

func solveProblem(input types.APIInput) (types.APIOutput, error) {
	Vars := []*types.Variable{}
	Eqns := []*types.Equation{}
	Settings := types.SolverSettings{}
	debug, err := parsers.ParseText(input, &Vars, &Eqns, &Settings)
	if err != nil {
		log.Printf("Something fails: %v", err)
		out := types.APIOutput{
			Info: types.Info{
				Errors: []string{err.Error()},
			},
		}
		return out, err
	}

	solution, errSol := solver.Solver(Vars, Eqns, Settings, debug)
	if errSol != nil {
		log.Print(errSol.Error())
		return solution, errSol
	}
	log.Print(solution)
	return solution, nil
}
