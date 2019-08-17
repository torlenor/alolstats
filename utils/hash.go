package utils

import (
	"sort"
	"strconv"
	"strings"
)

func HashSortedInt(i []int) string {
	sort.Ints(i)
	var s string
	for _, i := range i {
		s = s + strconv.Itoa(i) + "_"
	}

	return strings.TrimSuffix(s, "_")
}
