package gotype

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map{}
	for i := 0; i < 1000; i++ {
		m.Store(i, i)
	}
	assert.Equal(t, 1000, len(m.m))
	for i := 1001; i > 0; i-- {
		m.Delete(i)
	}
	m.Store(1, 10)
	fmt.Println(m.deleteTimes)
	//assert.Equal(t, int64(0),m.deleteTimes)
	fmt.Println(m.Load(1))
	//for i:=0;i<999;i++{
	//	m.Store(i,i)
	//}
	//assert.Equal(t, int64(0),m.deleteTimes)
}
