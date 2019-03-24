package statsrunner

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/stat"
)

// calcRelativeFrequency calculates the relative frequency for a given set of observations obs
func calcRelativeFrequency(obs []string) map[string]float64 {
	result := make(map[string]float64)
	totalObservations := len(obs)
	if totalObservations == 0 {
		return result
	}

	for _, o := range obs {
		result[o]++
	}

	for val, absCount := range result {
		result[val] = absCount / float64(totalObservations)
	}

	return result
}

func calcMeanStdDev(x, weights []float64) (mean, std float64) {
	return stat.MeanStdDev(x, weights)
}

func calcMeanStdDevUint16(x, weights []uint16) (mean, std float64) {
	xFloat64 := make([]float64, 0, len(x))
	for _, val := range x {
		xFloat64 = append(xFloat64, float64(val))
	}
	var weightsFloat64 []float64
	if weights != nil {
		weightsFloat64 := make([]float64, 0, len(weights))
		for _, val := range weights {
			weightsFloat64 = append(weightsFloat64, float64(val))
		}
	}

	return calcMeanStdDev(xFloat64, weightsFloat64)
}

func calcMeanStdDevUint32(x, weights []uint32) (mean, std float64) {
	xFloat64 := make([]float64, 0, len(x))
	for _, val := range x {
		xFloat64 = append(xFloat64, float64(val))
	}
	var weightsFloat64 []float64
	if weights != nil {
		weightsFloat64 := make([]float64, 0, len(weights))
		for _, val := range weights {
			weightsFloat64 = append(weightsFloat64, float64(val))
		}
	}

	return calcMeanStdDev(xFloat64, weightsFloat64)
}

func calcMedian(x, weights []float64) (float64, error) {
	if len(x) == 0 {
		return 0.0, fmt.Errorf("Cannot calculate Median: No elements in slice")
	}
	sort.Float64s(x)
	return stat.Quantile(0.5, stat.Empirical, x, weights), nil
}

func calcMedianUint16(x, weights []uint16) (float64, error) {
	if len(x) == 0 {
		return 0.0, fmt.Errorf("Cannot calculate Median: No elements in slice")
	}

	xFloat64 := make([]float64, 0, len(x))
	for _, val := range x {
		xFloat64 = append(xFloat64, float64(val))
	}
	var weightsFloat64 []float64
	if weights != nil {
		weightsFloat64 := make([]float64, 0, len(weights))
		for _, val := range weights {
			weightsFloat64 = append(weightsFloat64, float64(val))
		}
	}

	return calcMedian(xFloat64, weightsFloat64)
}

func calcMedianUint32(x, weights []uint32) (float64, error) {
	if len(x) == 0 {
		return 0.0, fmt.Errorf("Cannot calculate Median: No elements in slice")
	}

	xFloat64 := make([]float64, 0, len(x))
	for _, val := range x {
		xFloat64 = append(xFloat64, float64(val))
	}
	var weightsFloat64 []float64
	if weights != nil {
		weightsFloat64 := make([]float64, 0, len(weights))
		for _, val := range weights {
			weightsFloat64 = append(weightsFloat64, float64(val))
		}
	}

	return calcMedian(xFloat64, weightsFloat64)
}
