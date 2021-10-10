package test

import (
	rec "github.com/forwardalex/Ytool/recover"
)

func TestBlame() {
	defer rec.RecoverFromPanic("TestBlame1")
	panic("testing blame1")
}
