package common

import "sort"

type SI64 struct {
	key   string
	value int64
}

type SliceSI64 []*SI64

func (ss SliceSI64) Len() int {
	return len(ss)
}

func (ss SliceSI64) Less(i, j int) bool {
	return ss[i].value < ss[i].value
}

func (ss SliceSI64) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func SortMapSI64(sm map[string]int64) []string {
	ss := make([]*SI64, 0, len(sm))
	result := make([]string, 0, len(sm))
	for k, v := range sm {
		ss = append(ss, &SI64{
			key:   k,
			value: v,
		})
	}
	sort.Sort(SliceSI64(ss))
	for _, item := range ss {
		result = append(result, item.key)
	}
	return result
}
