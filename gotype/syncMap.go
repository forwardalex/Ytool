package gotype

import (
	"encoding/json"
	"sync"
)

// 读写锁 map
type RWMap struct {
	sync.RWMutex
	data map[interface{}]interface{}
}

// NewRWMap 新建RWMap
func NewRWMap() *RWMap {
	return &RWMap{
		data: make(map[interface{}]interface{}),
	}
}

// Get 获取某个key的值
func (m *RWMap) Get(key interface{}) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Set 设置某个key的值
func (m *RWMap) Set(key interface{}, val interface{}) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = val
}

// Delete 删除某个key
func (m *RWMap) Delete(key interface{}) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
}

// JsonMarshal json序列化
func (m *RWMap) JsonMarshal() ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	return json.Marshal(m.data)
}
