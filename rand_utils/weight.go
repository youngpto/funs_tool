package rand_utils

import (
	"errors"
	"github.com/youngpto/funs_tool/coll_utils"
)

func WeightGet(weights []int, random Random) (idx int) {
	if random == nil {
		random = globalRand
	}
	sum := 0
	for _, i := range weights {
		sum += i
	}
	if sum <= 0 {
		panic("sum <= 0")
	}
	if len(weights) == 0 {
		return 0
	}
	target := random.Intn(sum)
	var t int
	for i, weight := range weights {
		t += weight
		if t >= target {
			return i
		}
	}
	return len(weights) - 1
}

type RandWeight[T any] struct {
	elems  []WeightElement[T]
	random Random
}

type WeightElement[T any] struct {
	Weight  int
	Element T
}

func NewRandWeight[T any](wes []WeightElement[T], rds ...Random) *RandWeight[T] {
	es := make([]WeightElement[T], len(wes))
	copy(es, wes)

	rad := globalRand
	if len(rds) > 0 {
		rad = rds[0]
	}

	return &RandWeight[T]{
		elems:  es,
		random: rad,
	}
}

func (rnd *RandWeight[T]) Get() (t T, err error) {
	if len(rnd.elems) == 0 {
		err = errors.New("len is 0")
		return
	}
	weights := coll_utils.Convert2Type(rnd.elems, func(o WeightElement[T]) int {
		return o.Weight
	})
	idx := WeightGet(weights, rnd.random)
	return rnd.elems[idx].Element, nil
}

func (rnd *RandWeight[T]) GetAndDrop() (t T, err error) {
	if len(rnd.elems) == 0 {
		err = errors.New("len is 0")
		return
	}
	weights := coll_utils.Convert2Type(rnd.elems, func(o WeightElement[T]) int {
		return o.Weight
	})
	idx := WeightGet(weights, rnd.random)
	t = rnd.elems[idx].Element
	rnd.elems = append(rnd.elems[:idx], rnd.elems[idx+1:]...)
	return
}

func (rnd *RandWeight[T]) IsEmpty() bool {
	return len(rnd.elems) == 0
}

func (rnd *RandWeight[T]) SetElems(elems []WeightElement[T]) {
	es := make([]WeightElement[T], len(elems))
	copy(es, elems)
	rnd.elems = elems
}

func OnceWeightElem[T any](elems []WeightElement[T], random Random) (t T, err error) {
	if len(elems) == 0 {
		err = errors.New("len is 0")
		return
	}
	if random == nil {
		random = globalRand
	}
	weights := coll_utils.Convert2Type(elems, func(o WeightElement[T]) int {
		return o.Weight
	})
	idx := WeightGet(weights, random)
	return elems[idx].Element, nil
}
