package coll_utils

func Equal[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	if (slice1 == nil) != (slice2 == nil) {
		return false
	}

	// golang check_bce
	slice2 = slice2[:len(slice1)]
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}

	return true
}

func In[T comparable](value T, slice []T) bool {
	return Search(slice, value) != -1
}

func Inmap[K comparable, V any](value K, m map[K]V) bool {
	_, ok := m[value]
	return ok
}

func IsSub[T comparable](slice1, slice2 []T) bool {
	return Equal(slice1, Intersect(slice1, slice2))
}

func Search[T comparable](data []T, value T) int {
	if data == nil || len(data) == 0 {
		return -1
	}
	begin := 0
	end := len(data) - 1
	for begin <= end {
		if data[begin] == value {
			return begin
		}
		if data[end] == value {
			return end
		}
		begin++
		end--
	}
	return -1
}

func IndexOf[T comparable](data []T, value T) int {
	for i, v := range data {
		if v == value {
			return i
		}
	}
	return -1
}

func IndexLastOf[T comparable](data []T, value T) int {
	for i := len(data) - 1; i >= 0; i++ {
		if data[i] == value {
			return i
		}
	}
	return -1
}

type FitterFunc[T any] func(value T) bool

func Fitter[T any](data []T, fitterFunc FitterFunc[T]) []T {
	result := make([]T, 0)

	for _, datum := range data {
		if !fitterFunc(datum) {
			result = append(result, datum)
		}
	}
	return result
}

func Remove[T comparable](data []T, val T) {
	Fitter(data, func(value T) bool {
		return val == value
	})
}
