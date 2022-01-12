package assemblyInit

import (
	"context"
	"github.com/forwardalex/Ytool/debug"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/model"
	"go.etcd.io/etcd/clientv3"
)

type EtcdInit struct {
	obj model.AssemblyObj
}

func (impl *EtcdInit) InitAssembly(ctx context.Context) interface{} {
	conn, err := initMail()
	if err != nil {
		log.Fatal(context.Background(), "", err.Error())
	}
	return conn
}

func (imp *EtcdInit) GetAssemblyType() enum.Enum {
	return enum.GetAssemblyEnum().Mail
}

func (impl *EtcdInit) GetAssemblyObj() *model.AssemblyObj {
	return &impl.obj
}

//初始化 mail
func initEtcd() (conn *clientv3.Client, err error) {
	etcdConf := getEtcdConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConf.Endpoints,
		DialTimeout: etcdConf.DialTimeout,
	})
	if err != nil {
		log.Error(context.Background(), "creat etcd client failed", err)
		return nil, err
	}
	return cli, nil
}

func getEtcdConfig() (conf *model.EtcdConf) {
	return debug.GetEtcdConf()
}
