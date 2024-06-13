package coll_utils

// 差集
func Difference[T comparable](slice1, slice2 []T) []T {
	countMap := make(map[T]int)
	result := make([]T, 0)

	intersect := Intersect(slice1, slice2)

	for _, v := range intersect {
		countMap[v]++
	}

	for _, value := range slice1 {
		times, _ := countMap[value]
		if times == 0 {
			result = append(result, value)
		} else {
			countMap[value]--
		}
	}
	return result
}

// 并集
func Union[T comparable](slice1, slice2 []T) []T {
	countMap := make(map[T]int)
	result := make([]T, len(slice1))

	for index, v := range slice1 {
		countMap[v]++
		result[index] = v
	}

	for _, v := range slice2 {
		times, _ := countMap[v]
		if times == 0 {
			slice1 = append(slice1, v)
		} else {
			countMap[v]--
		}
	}
	return slice1
}

// 交集
func Intersect[T comparable](slice1, slice2 []T) []T {
	result := make([]T, 0)

	if len(slice1) == 0 || len(slice2) == 0 {
		return result
	}

	countMap := make(map[T]int)
	for _, v := range slice1 {
		countMap[v]++
	}

	for _, v := range slice2 {
		times, _ := countMap[v]
		if times > 0 {
			countMap[v]--
			result = append(result, v)
		}
	}
	return result
}

// 去重
func Duplicate[T comparable](slice []T) []T {
	result := make([]T, 0)
	set := make(map[T]struct{})
	for _, value := range slice {
		set[value] = struct{}{}
	}

	for val, _ := range set {
		result = append(result, val)
	}

	return result
}
