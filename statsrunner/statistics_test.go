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

func TestCalcMeanStdDev(t *testing.T) {
	values := []float64{2, 2, 6, 6, 6, 6, 8, 10}
	actualMean, actualStdDev := calcMeanStdDev(values, nil)
	if !almostEqual(actualMean, 5.75) {
		t.Errorf("Result for mean is wrong, was = %f, should be = %f", actualMean, 5.75)
	}
	if !almostEqual(actualStdDev, 2.71240536372107) {
		t.Errorf("Result for standard deviation is wrong")
	}
}

func TestCalcMedian(t *testing.T) {
	values := []float64{2, 2, 6, 6, 6, 6, 8, 10}
	actualMedian := calcMedian(values, nil)
	if !almostEqual(actualMedian, 6) {
		t.Errorf("Result for median is wrong, was = %f, should be = %f", actualMedian, 6.0)
	}

	values = []float64{3, 6, 7, 8, 8, 9, 10, 13, 15, 16, 20}
	actualMedian = calcMedian(values, nil)
	if !almostEqual(actualMedian, 9) {
		t.Errorf("Result for median is wrong, was = %f, should be = %f", actualMedian, 9.0)
	}
}
