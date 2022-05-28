// Package casbin TODO
package casbin

import (
	"fmt"
	"strings"

	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/store/db"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"golang.org/x/net/context"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"strconv"
	"time"

	// MySQL初始化
	_ "github.com/go-sql-driver/mysql"
)

var (
	// CasbinEnforcer TODO
	gormDB         *gorm.DB
	CasbinEnforcer *casbin.SyncedEnforcer
)

var chCasbinPolicy chan *chCasbinPolicyItem

type chCasbinPolicyItem struct {
	ctx context.Context
	e   *casbin.SyncedEnforcer
}

// CasbinRule TODO
// Increase the column size to 512.
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:512;uniqueIndex:unique_index"`
	V0    string `gorm:"size:512;uniqueIndex:unique_index"`
	V1    string `gorm:"size:512;uniqueIndex:unique_index"`
	V2    string `gorm:"size:512;uniqueIndex:unique_index"`
	V3    string `gorm:"size:512;uniqueIndex:unique_index"`
	V4    string `gorm:"size:512;uniqueIndex:unique_index"`
	V5    string `gorm:"size:512;uniqueIndex:unique_index"`
}

// RedisConf TODO
type RedisConf struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func initGorm() {
	dsn := db.GetDBConfig()
	var err error
	gormDB, err = gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Infof(context.Background(), "db conn::%s\n", dsn)
		return
	}
}

// InitCasbin 初始化Cabin权限服务
func InitCasbin() error {
	text :=
		`
	[request_definition]
	r = sub, dom, obj, act
	
	[policy_definition]
	p = sub, dom, obj, act
	
	[role_definition]
	g = _, _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
	`
	initGorm()
	a, _ := gormadapter.NewAdapterByDBWithCustomTable(gormDB, &CasbinRule{}, "auth.casbin_rule")
	m, _ := model.NewModelFromString(text)

	CasbinEnforcer, _ = casbin.NewSyncedEnforcer(m, a)
	if os.Getenv("ENV_NAME") != "Prod" && os.Getenv("ENV_NAME") != "Pre" {
		CasbinEnforcer.EnableLog(true)
	}
	CasbinEnforcer.AddPolicy()
	CasbinEnforcer.StartAutoLoadPolicy(10 * time.Minute)
	CasbinEnforcer.LoadPolicy()
	initLoadCasbin()

	err := initWatcher()
	if err != nil {
		log.Error(context.TODO(), "casbin watcher start failed:", err.Error())
		return err
	}
	return nil
}

func initWatcher() error {
	var (
		err         error
		redisConfig RedisConf
	)
	config := db.GetRedisConfig()
	redisConfig.Addr = config.Host + ":" + strconv.Itoa(config.Port)
	// 启动监听器
	env := os.Getenv("ENV_NAME")
	w, err := rediswatcher.NewWatcher(redisConfig.Addr, rediswatcher.WatcherOptions{
		Options:    *db.GetRedisConn().Options(),
		Channel:    fmt.Sprintf("/casbin_%s", env),
		IgnoreSelf: true,
	})
	if err != nil {
		return err
	}

	err = CasbinEnforcer.SetWatcher(w)
	if err != nil {
		return err
	}
	err = w.SetUpdateCallback(updateCallback)
	if err != nil {
		return err
	}
	return nil
}

func updateCallback(msg string) {
	log.Info(context.Background(), "update permission:", msg)
	LoadCasbinPolicy(context.Background(), CasbinEnforcer)
}

// initLoadCasbin 初始化casbin异步加载
func initLoadCasbin() {
	chCasbinPolicy = make(chan *chCasbinPolicyItem, 1)
	go func() {
		for item := range chCasbinPolicy {
			err := item.e.LoadPolicy()
			if err != nil {
				log.Errorf(item.ctx, "The load casbin policy error: %s", err.Error())
			} else {
				log.Infof(item.ctx, "casbin policy flush success")
			}
		}
	}()
}

// LoadCasbinPolicy 异步加载casbin权限策略
func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {

	if len(chCasbinPolicy) > 0 {
		log.Infof(ctx, "The load casbin policy is already in the wait queue")
		return
	}
	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}

// GetUserPermissione TODO
func GetUserPermission(user string, domain string, obj string) [][]string {
	var list [][]string
	results := CasbinEnforcer.GetPermissionsForUser(user, domain)
	for i := 0; i < len(results); i++ {
		name := results[i][2]
		action := results[i][3]
		if name == obj {
			list = append(list, []string{name, action})
		}
	}
	return list
}

// GetUserPermissionByTypeAndAction TODO
func GetUserPermissionByTypeAndAction(user string, domain string, ptype string, paction string) []string {
	var list []string
	results, err := CasbinEnforcer.GetImplicitPermissionsForUser(user, domain)
	if err != nil {
		log.Error(context.TODO(), "获取权限失败:", err)
		return list
	}
	for i := 0; i < len(results); i++ {
		name := results[i][2]
		action := results[i][3]
		if strings.Contains(name, ptype+"::") && action == paction {
			arr := strings.Split(name, "::")
			if len(arr) >= 2 {
				list = append(list, arr[len(arr)-1])
			}
		}

	}
	return list
}
