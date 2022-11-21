package user

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

//用户角色 1:考生  2:考官 3：审核员  7：普通管理员 9：超级管理员
var RoleTypeOptions = map[int]string{
	0: "管理员",
	1: "普通用户",
	2: "考官",
	3: "审核员",
}

type Entity struct {
	Id           int         `orm:"id,table_comment:'用户表'" json:"user_id"`
	Username     string      `orm:"username,size: 10,comment:'姓名'" json:"username"`
	RoleType     int         `orm:"role_type,default:2,comment:'角色 0：管理员 1：普通用户组'" json:"role_type"`
	Password     string      `orm:"password,size:32,comment:'角色时间'" json:"password"`
	CreateAt     *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
	DeleteAt     *gtime.Time `orm:"delete_at,comment:'删除时间'" json:"delete_at"`
	Status       int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'"`
	RoleName     string      `json:"role_name"`
	SerialNumber string      `json:"serial_number"`
	AuthorNumber string      `json:"author_number"`
}

var (
	Table       = "user"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "user c"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type UserList struct {
	Id       int    `json:"user_id"`
	Username string `json:"username"`
	RoleType int    `json:"role_type"`
	RoleName string `json:"role_name"`
	Password string `json:"password"`
	UpdateAt int    `json:"update_at"`
	DeleteAt int    `json:"delete_at"`
}

type UserData struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleType int    `json:"role_type"`
	UpdateAt int    `json:"update_at"`
	DeleteAt int    `json:"delete_at"`
}

type LoginReq struct {
	UserId       int    `p:"user_id" v:"required#参数错误：用户id必填"`
	Password     string `p:"password" v:"required#参数错误：登录密码必填"`
	AuthorNumber string `p:"author_number"  v:"required#参数错误：客户端授权码"`
}

type SaveUserDataReq struct {
	UserData string `json:"user_data" p:"user_data" v:"required#参数错误,用户数据不存在"` // 主任密钥中的用户数据字串
}
