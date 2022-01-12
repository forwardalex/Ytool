package db

import (
	"database/sql"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/layzeInit"
)

var conn *sql.DB

func GetConn() *sql.DB {
	if conn != nil {
		return conn
	}
	assembly := layzeInit.GetAssembly(enum.GetAssemblyEnum().MySQL)

	if assembly == nil {
		return nil
	}

	return assembly.(*sql.DB)
}
