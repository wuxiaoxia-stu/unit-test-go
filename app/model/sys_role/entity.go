package sys_role

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table sys_admin.
type Entity struct {
	Id       int         `orm:"id,primary,table_comment:'角色管理表'" json:"id"`
	Name     string      `orm:"name,size:50,not null,default:"",comment:'角色名称'" json:"name"`
	CreateAt *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt *gtime.Time `orm:"update_at" json:"update_at"`
	Status   int         `orm:"status,size:2,not null,default:1" json:"status"`
	S        int         `json:"s"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "sys_admin"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "sys_admin sa"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
