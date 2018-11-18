package statsrunner

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
