package main

import (
	"imhotep/parsers"
	"imhotep/solver"
	"imhotep/types"
	"imhotep/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/vars", func(c *fiber.Ctx) error {
		input := new(types.APIInput)
		var logger []string
		err := c.BodyParser(input)
		if err != nil {
			out := types.APIOutput{
				Info: types.Info{
					Errors: []string{err.Error()},
				},
			}
			if input.Debug {
				out.Info.Logs = logger
			}
			return c.JSON(out)
		}
		out, errSol := solveProblem(*input, true, &logger)
		if input.Debug {
			out.Info.Logs = logger
		}
		if errSol != nil {
			out.Info.Errors = []string{errSol.Error()}
			return c.JSON(out)
		}
		return c.JSON(out)

	})

	app.Post("/solve", func(c *fiber.Ctx) error {
		input := new(types.APIInput)
		var logger []string
		err := c.BodyParser(input)
		if err != nil {
			out := types.APIOutput{
				Info: types.Info{
					Errors: []string{err.Error()},
				},
			}
			if input.Debug {
				out.Info.Logs = logger
			}
			return c.JSON(out)
		}
		solution, errSol := solveProblem(*input, false, &logger)
		if input.Debug {
			solution.Info.Logs = logger
		}
		if errSol != nil {
			solution.Info.Errors = []string{errSol.Error()}
			return c.JSON(solution)
		}
		return c.JSON(solution)
	})
	log.Fatal(app.Listen(":3000"))

}

func solveProblem(input types.APIInput, onlyVars bool, logger *[]string) (types.APIOutput, error) {
	Vars := []*types.Variable{}
	Eqns := []*types.Equation{}
	Settings := types.SolverSettings{}
	debug, err := parsers.ParseText(input, &Vars, &Eqns, &Settings, onlyVars, logger)
	if err != nil {
		utils.HandleLog(logger, "Something fails: %v", err)
		out := types.APIOutput{
			Info: types.Info{
				Errors: []string{err.Error()},
			},
		}
		return out, err
	}
	if onlyVars {
		out := types.APIOutput{}
		out.Vars = []types.VariableJSON{}
		for _, var_ := range Vars {
			varjson := types.VariableJSON{}
			varjson.Name = var_.Name
			out.Vars = append(out.Vars, varjson)
		}
		return out, nil
	}
	solution, errSol := solver.Solver(Vars, Eqns, Settings, debug, logger)
	if errSol != nil {
		utils.HandleLog(logger, errSol.Error())
		return solution, errSol
	}
	log.Print(solution)
	return solution, nil
}
