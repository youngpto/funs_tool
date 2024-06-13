package coll_utils

import (
	"github.com/youngpto/funs_tool/define"
	"github.com/youngpto/funs_tool/math_utils"
	"strconv"
)

type ConvertFunc[O, N any] func(o O) N

func Convert2Type[O, N any](data []O, convertFunc ConvertFunc[O, N]) []N {
	var result []N
	for _, datum := range data {
		result = append(result, convertFunc(datum))
	}
	return result
}

func Int2String[T define.Integer](data []T) []string {
	return Convert2Type(data, func(i T) string {
		return strconv.FormatInt(int64(i), 10)
	})
}

func String2Int[T define.Integer](data []string) []T {
	return Convert2Type(data, func(s string) T {
		i, _ := strconv.ParseInt(s, 10, 64)
		return T(i)
	})
}

type Tuple = []any

func Zip(data ...Tuple) []Tuple {
	var result []Tuple
	length := math_utils.Min(Convert2Type(data, func(tuple Tuple) int {
		return len(tuple)
	})...)

	for i := 0; i < length; i++ {
		var tuple Tuple
		for _, datum := range data {
			tuple = append(tuple, datum[i])
		}
		result = append(result, tuple)
	}
	return result
}

func Map[K comparable, V any](keys []K, values []V) map[K]V {
	result := make(map[K]V)
	length := math_utils.Min(len(keys), len(values))

	for i := 0; i < length; i++ {
		result[keys[i]] = values[i]
	}
	return result
}
