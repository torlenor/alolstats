package statsrunner

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestCalcRelativeFrequency(t *testing.T) {
	var obs []string
	result := calcRelativeFrequency(obs)
	if len(result) != 0 {
		t.Errorf("Calculation on an empty set shall give an empty result")
	}

	obs = append(obs, "A")
	obs = append(obs, "A")
	obs = append(obs, "A")
	obs = append(obs, "B")
	obs = append(obs, "A")
	obs = append(obs, "B")
	obs = append(obs, "C")
	obs = append(obs, "C")
	obs = append(obs, "c")
	obs = append(obs, "AAAAA")
	// A = 4, B = 2, C = 2, c = 1, AAAAA = 1 -> 10 in total
	// -> A = 0.4, B = 0.2, C = 0.2, c = 0.1, AAAAA = 0.1

	result = calcRelativeFrequency(obs)
	if len(result) != 5 {
		t.Errorf("There should have been 5 disctinct value types")
	}
	if val, ok := result["A"]; ok {
		if !almostEqual(val, 0.4) {
			t.Errorf("Result for 'A' is wrong")
		}
	} else {
		t.Errorf("There should be a result for 'A'")
	}

	if val, ok := result["B"]; ok {
		if !almostEqual(val, 0.2) {
			t.Errorf("Result for 'B' is wrong")
		}
	} else {
		t.Errorf("There should be a result for 'B'")
	}

	if val, ok := result["C"]; ok {
		if !almostEqual(val, 0.2) {
			t.Errorf("Result for 'C' is wrong")
		}
	} else {
		t.Errorf("There should be a result for 'C'")
	}

	if val, ok := result["c"]; ok {
		if !almostEqual(val, 0.1) {
			t.Errorf("Result for 'c' is wrong")
		}
	} else {
		t.Errorf("There should be a result for 'c'")
	}

	if val, ok := result["AAAAA"]; ok {
		if !almostEqual(val, 0.1) {
			t.Errorf("Result for 'AAAAA' is wrong")
		}
	} else {
		t.Errorf("There should be a result for 'AAAAA'")
	}
}
