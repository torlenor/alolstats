package utils

import (
	"sort"
	"strconv"
	"strings"
)

// HashInt takes a int slice and returns a string value
// which uniquely identifies this set of ints.
func HashInt(i []int) string {
	var s string
	for _, i := range i {
		s = s + strconv.Itoa(i) + "_"
	}

	return strings.TrimSuffix(s, "_")
}

// HashSortedInt takes a int slice, sorts it and returns a string value
// which uniquely identifies this set of ints.
// Beware: This sorts the slice passed to the function!
func HashSortedInt(i []int) string {
	sort.Ints(i)
	return HashInt(i)
}
