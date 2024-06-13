package math_utils

import "github.com/youngpto/funs_tool/define"

func Max[T define.Number](data ...T) T {
	var max T
	for i, v := range data {
		if i == 0 {
			max = v
			continue
		}
		if v > max {
			max = v
		}
	}
	return max
}

func Min[T define.Number](data ...T) T {
	var min T
	for i, v := range data {
		if i == 0 {
			min = v
			continue
		}
		if v < min {
			min = v
		}
	}
	return min
}

func Sum[T define.Number](data ...T) T {
	var sum T
	for _, v := range data {
		sum += v
	}
	return sum
}

func Abs[T define.Number](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func In[T define.Number](val T, start, end T) bool {
	return val >= start && val <= end
}

func Ternary[T any](ifVal T, cond bool, elVal T) (t T) {
	if cond {
		t = ifVal
	} else {
		t = elVal
	}
	return
}
