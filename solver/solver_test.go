package solver

import (
	"reflect"
	"testing"
)

func TestConvertFullPseudograph(t *testing.T) {

	// tarjan paper example
	adjacencyMatrix := [][]int{{1, 1, 1, 0, 0, 0, 0, 0}, {0, 1, 1, 0, 0, 0, 0, 1}, {0, 1, 0, 1, 0, 0, 0, 0}, {0, 0, 1, 0, 1, 1, 0, 0}, {1, 0, 0, 1, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 0, 1, 0}, {0, 1, 1, 0, 0, 0, 1, 0}, {0, 0, 1, 0, 0, 0, 0, 1}}
	expectedReOrderEqn := map[int]int{0: 0, 1: 1, 2: 3, 3: 5, 4: 4, 5: 6, 6: 2, 7: 7}

	_, reOrderEqn, err := ConvertFullPseudograph(adjacencyMatrix)

	if err != nil {
		t.Logf("Fails to reorder eqns: %v", err)
		t.Fail()
	}

	res := reflect.DeepEqual(reOrderEqn, expectedReOrderEqn)
	if !res {
		t.Logf("Bad reorder eqns, expected %v, got %v", expectedReOrderEqn, reOrderEqn)
		t.Fail()
	}

}
