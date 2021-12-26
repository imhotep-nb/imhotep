package parsers

import "testing"

func TestParseExplicitUnits(t *testing.T) {

	eqnText := "x + y - 10[ft] = 50[ft]"
	expectedOut := "x + y - 10*3.048e-01 = 50*3.048e-01"

	out, err := ParseExplicitUnits(eqnText)

	if err != nil {
		t.Log("Algo no cuadra")
		t.Fail()
	}

	if out != expectedOut {
		t.Log("Shit! :o")
		t.Fail()
	}
}
