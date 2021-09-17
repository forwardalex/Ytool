package test

import (
	rec "Ytool/recover"
)

func TestBlame() {
	defer rec.RecoverFromPanic("TestBlame1")
	panic("testing blame1")
}
