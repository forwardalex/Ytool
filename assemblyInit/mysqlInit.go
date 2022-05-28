package assemblyInit

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/forwardalex/Ytool/debug"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/model"
	_ "github.com/go-sql-driver/mysql" //mysql 驱动
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Dsn string

type MySqlInit struct {
	obj model.AssemblyObj
}

func (impl *MySqlInit) InitAssembly(ctx context.Context) interface{} {
	db, err := initDB()
	if err != nil {
		log.Fatal(context.Background(), "", err.Error())
	}
	return db
}

func (*MySqlInit) GetAssemblyType() enum.Enum {
	return enum.GetAssemblyEnum().MySQL
}

func (impl *MySqlInit) GetAssemblyObj() *model.AssemblyObj {
	return &impl.obj
}

//InitDB 初始化数据库
func initDB() (*sql.DB, error) {
	var (
		dbConf model.DbConf
		err    error
	)

	err = getDBConfig(&dbConf)

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	fmt.Println(path[:index])

	// 如果是开发环境，手动指定数据库地址
	debug.ConfigDebugDB(&dbConf)

	if err != nil {
		log.Error(context.Background(), "err ", err)
		return nil, err
	}
	Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.User, dbConf.Password, dbConf.Host,
		dbConf.Port, dbConf.Database)
	log.Info(context.Background(), "db=", Dsn)
	conn, err := sql.Open("mysql", Dsn)
	if err != nil {
		log.Error(context.Background(), "open mysql err:", err.Error())
		return nil, err
	}
	return conn, nil
}

func getDBConfig(conf *model.DbConf) error {
	return nil
}
