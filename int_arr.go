package golib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type IntArr []int64

func (ia IntArr) StringArr() []string {
	s := []string{}
	for _, i := range ia {
		s = append(s, strconv.FormatInt(i, 10))
	}
	return s
}

// Join returns int elements joined by "el". As a special case
// this returns "<empty>" if the Array is empy.
func (ia IntArr) Join(el string) string {
	s := ia.StringArr()
	if len(s) == 0 {
		return "<empty>"
	} else {
		return strings.Join(s, el)
	}
}

func (ia IntArr) AvgInt64() int64 {
	s := ia.Sum()
	return int64(math.Round(float64(s) / float64(len(ia))))
}

func (ia IntArr) Sum() int64 {
	var s int64
	for _, i := range ia {
		s += i
	}
	return s
}

// Min returns the smallest number in the array. 0 if the array is empty.
func (ia IntArr) Min() int64 {
	var m int64

	for idx, i := range ia {
		if idx == 0 {
			// seed
			m = i
			continue
		}
		if i < m {
			m = i
		}
	}
	return m
}

// Max returns the biggest number in the array. 0 if the array is empty.
func (ia IntArr) Max() int64 {
	var m int64
	for idx, i := range ia {
		if idx == 0 {
			// seed
			m = i
			continue
		}
		if i > m {
			m = i
		}
	}
	return m
}

func (ic IntArr) Contains(i2 int64) bool {
	return ic.IndexOf(i2) > -1
}

func (ic IntArr) IndexOf(i2 int64) int {
	if ic == nil {
		return -1
	}
	for idx, i := range ic {
		if i == i2 {
			return idx
		}
	}
	return -1
}

// Add appends all given i to the Array.
func (ic *IntArr) Add(i_s ...int64) {
	for _, i := range i_s {
		*ic = append(*ic, i)
	}
}

// AddUnique only adds if i is not in the Array yet
func (ic *IntArr) AddUnique(i_s ...int64) {
	for _, i := range i_s {
		if !ic.Contains(i) {
			ic.Add(i)
		}
	}
}

// Remove removes i_s from the array (if exist)
func (ic *IntArr) Remove(i_s ...int64) {
	iMap := map[int64]bool{}
	for _, i := range i_s {
		iMap[i] = true
	}

	ic2 := IntArr{}
	for _, i := range *ic {
		if !iMap[i] {
			ic2 = append(ic2, i)
		}
	}
	*ic = ic2[:]
}

func (ic1 IntArr) Diff(ic2 IntArr) IntDiffMap {
	m := IntDiffMap{
		New:     IntArr{},
		Same:    IntArr{},
		Removed: IntArr{},
	}

	for _, i1 := range ic1 {
		if ic2.Contains(i1) {
			m.Same.Add(i1)
		} else {
			m.Removed.Add(i1)
		}
	}

	for _, i2 := range ic2 {
		if !ic1.Contains(i2) {
			m.New.Add(i2)
		}
	}

	return m
}

func Int64InArray(id int64, arr []int64) bool {
	for _, item := range arr {
		if item == id {
			return true
		}
	}
	return false
}

type IntDiffMap struct {
	New     IntArr
	Same    IntArr
	Removed IntArr
}

func (idm IntDiffMap) String() string {
	return fmt.Sprintf(`New: %v Same: %v Removed: %v`, idm.New, idm.Same, idm.Removed)
}
