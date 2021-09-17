package test

import (
	rec "Ytool/recover"
)

func TestBlame() {
	defer rec.RecoverFromPanic("TestBlame")
	panic("testing blame")
}
