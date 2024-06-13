package rand_utils

func RandGet[T any](slice []T, count int, random Random) []T {
	if random == nil {
		random = globalRand
	}
	if count <= 0 {
		panic("rand count not 0")
	}
	if count == 1 {
		return []T{slice[random.Intn(len(slice))]}
	}
	var result []T
	perm := random.Perm(len(slice))
	for i := 0; i < count; i++ {
		result = append(result, slice[perm[i]])
	}
	return result
}

func RandOne[T any](slice []T, random Random) T {
	if random == nil {
		random = globalRand
	}

	return slice[random.Intn(len(slice))]
}
