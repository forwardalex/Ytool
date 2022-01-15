package struct2struct

import (
	"github.com/jinzhu/copier"
)

func StructCopy(to, from interface{}) error {
	return copier.Copy(to, from)
}
