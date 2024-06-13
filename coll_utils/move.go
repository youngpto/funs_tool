package coll_utils

func MoveRight[T any](slice []T, step int) []T {
	length := len(slice)
	step = step % length
	if length == 0 || step == 0 {
		return slice
	}

	if step < 0 {
		step = length + step
	}

	newS := make([]T, length)
	copy(newS[step:], slice[:length-step])
	for i := 0; i < step; i++ {
		newS[i] = slice[length-step+i]
	}

	return newS
}

func MoveLeft[T any](slice []T, step int) []T {
	return MoveRight(slice, -step)
}

func Reverse[T any](slice []T) []T {
	left := 0
	right := len(slice) - 1
	for left < right {
		slice[left], slice[right] = slice[right], slice[left]
		left++
		right--
	}
	return slice
}
