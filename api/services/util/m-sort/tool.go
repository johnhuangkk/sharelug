package m_sort

import (
	"errors"
	"sort"
)

func SortStringsAscAndGetFirst(strs []string) (string, error) {
	if len(strs) == 0 {
		return "",errors.New("Slice is empty")
	}
	sort.Slice(strs, func(i, j int) bool {
		return strs[i] < strs[j]
	})
	return strs[0],nil
}
