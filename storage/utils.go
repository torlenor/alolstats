package storage

import (
	"fmt"
	"net/url"
)

func extractURLStringParameter(parameters url.Values, parameterName string) (string, error) {
	var parameterValue string
	if val, ok := parameters[parameterName]; ok {
		if len(val) == 0 {
			return "", fmt.Errorf("%s parameter was empty in request", parameterName)
		} else if len(val) > 1 {
			return "", fmt.Errorf("Got more than one %s parameters in request", parameterName)
		}
		parameterValue = val[0]
	} else {
		return "", fmt.Errorf("%s parameter was missing in request", parameterName)
	}
	return parameterValue, nil
}
