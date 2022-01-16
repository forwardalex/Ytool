package gotype

import "errors"

type unionSet struct {
	set []int
}

func NewUnionSet(size int) *unionSet {
	buf := make([]int, size)
	for i := 0; i < size; i++ {
		buf[i] = i // 初始时，所有节点均指向自己
	}
	return &unionSet{set: buf}
}

func (set *unionSet) GetSize() int {
	return len(set.set)
}

func (set *unionSet) GetID(p int) (int, error) {
	if p < 0 || p > len(set.set) {
		return 0, errors.New(
			"failed to get ID,index is illegal.")
	}
	return set.set[p], nil
}

func (set *unionSet) IsConnected(p, q int) (bool, error) {
	if p < 0 || p > len(set.set) || q < 0 || q > len(set.set) {
		return false, errors.New(
			"failed to get ID,index is illegal.")
	}
	return set.set[p] == set.set[q], nil
}

func (set *unionSet) Union(p, q int) error {
	b, err := set.IsConnected(p, q)
	if err != nil {
		return err
	}

	if !b {
		pID := set.set[p]
		qID := set.set[q]
		for k, v := range set.set {
			if v == pID {
				set.set[k] = qID
			}
		}
	}
	return nil
}
