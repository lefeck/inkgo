package database

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

var modelPath = "config/rbac_model.conf"

func NewCasbin(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("failed to create Casbin adapter: %s", err)
	}
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.Fatalf("failed to create Casbin enforcer: %s", err)
	}
	// 注册自定义函数 isAdmin
	enforcer.AddFunction("isAdmin", func(args ...interface{}) (interface{}, error) {
		role := args[0].(string)
		return role == "admin", nil
	})

	// 注册 keyMatch2 用于路径匹配
	enforcer.AddFunction("keyMatch2", util.KeyMatch2Func)

	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy: %s", err)
	}
	return enforcer
}
