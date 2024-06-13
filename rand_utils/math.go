package rand_utils

// 随机区域 [min,max)左闭右开
func RandScope(min, max int, random Random) int {
	if random == nil {
		random = globalRand
	}

	if min == max {
		return min
	}
	return random.Intn(max-min) + min
}

// 随机区域 [min,max]左闭右闭
func CLCRRandScope(min, max int, random Random) int {
	if random == nil {
		random = globalRand
	}

	if min >= max {
		return min
	}
	return random.Intn(max+1-min) + min
}

func WeightRandom(weights []int, random Random) int {
	if random == nil {
		random = globalRand
	}

	if len(weights) == 0 {
		panic("non element random")
	}
	if len(weights) == 1 {
		return 0
	}
	var sum int
	for _, w := range weights {
		sum += w
	}
	r := random.Intn(sum)

	var t int
	for i, w := range weights {
		t += w
		if t > r {
			return i
		}
	}
	return len(weights) - 1
}
