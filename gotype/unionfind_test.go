package gotype

import (
	"fmt"
	"testing"
)

func TestUnionSet(t *testing.T) {
	set := NewUnionSet(10)
	set.Union(1, 2)
	set.Union(1, 3)
	set.Union(2, 3)
	set.Union(2, 4)
	is, err := set.IsConnected(1, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(set.set)
	fmt.Println(is)

}
