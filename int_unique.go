package golib

import (
	"encoding/json"
	"slices"
	"sync"
)

type Int64Unique struct {
	Ids   []int64
	idMap map[int64]bool
	mtx   sync.Mutex
}

func NewInt64Unique() *Int64Unique {
	m := Int64Unique{}
	m.idMap = make(map[int64]bool, 0)
	return &m
}

func (m *Int64Unique) Add(ids ...int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for _, id := range ids {
		if !m.Contains(id) {
			m.Ids = append(m.Ids, id)
			m.idMap[id] = true
		}
	}
}

func (m *Int64Unique) Remove(ids ...int64) (n int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for _, id := range ids {
		if m.Contains(id) {
			m.Ids = slices.DeleteFunc(m.Ids, func(id2 int64) bool { return id == id2 })
			delete(m.idMap, id)
			n++
		}
	}
	return n
}

func (m *Int64Unique) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Ids)
}

func (m *Int64Unique) Contains(id int64) bool {
	_, exists := m.idMap[id]
	return exists
}

func (m *Int64Unique) Count() int {
	return len(m.Ids)
}
