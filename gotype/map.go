package gotype

import (
	"reflect"
	"sync"
)

type Map struct {
	m map[interface{}]interface{}
	sync.Mutex
	deleteTimes int64
}

func (m *Map) Store(key, value interface{}) {
	m.Lock()
	defer m.Unlock()
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
	if m.deleteTimes > 1000 {
		m.copy()
	}
	m.m[key] = value
}

func (m *Map) Load(key interface{}) (interface{}, reflect.Type) {
	m.Lock()
	defer m.Unlock()
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
	_, ok := m.m[key]
	if ok {
		return m.m[key], reflect.TypeOf(m.m[key])
	}
	return nil, nil
}

func (m *Map) Delete(key interface{}) {
	m.Lock()
	defer m.Unlock()
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
	delete(m.m, key)
	m.deleteTimes++
}

//deltetime over limit copy a new map instead of old
func (m *Map) copy() {
	newmap := make(map[interface{}]interface{}, len(m.m))
	for k, v := range m.m {
		newmap[k] = v
	}
	m.deleteTimes = 0
	m.m = newmap
}
