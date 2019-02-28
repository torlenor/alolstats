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

func calcMedian(x, weights []float64) (float64, error) {
	if len(x) == 0 {
		return 0.0, fmt.Errorf("Cannot calculate Median: No elements in slice")
	}
	sort.Float64s(x)
	return stat.Quantile(0.5, stat.Empirical, x, weights), nil
}
