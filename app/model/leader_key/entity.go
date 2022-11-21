package leader_key

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id            int         `orm:"id,primary,table_comment:'主任秘钥绑定记录'" json:"id"`
	LeaderKeyCode string      `orm:"leader_key_code,size:5,comment:'秘钥代码'" json:"leader_key_code"`
	SerialNumber  string      `orm:"serial_number,size:16,comment:'秘钥系列号'" json:"leader_key_serial_number"`
	CreateAt      *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt      *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
	Status        int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'" json:"status" default:"1"`
}

var (
	Table       = "leader_key"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "leader_key lk"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type LeaderBindReq struct {
	//LeaderKeyCode      string `json:"leader_key_code" p:"leader_key_code" v:"required#参数错误,主任秘钥编码必填"`            // 主任密钥唯一编号
	LeaderSerialNumber string `json:"leader_serial_number" p:"leader_serial_number" v:"required#参数错误,主任秘钥系列号必填"` // 主任密钥唯一编号
	AuthorNumber       string `json:"author_number" p:"author_number" v:"required#参数错误,客户端授权编号必填"`               // 客户端授权id
	UserData           string `json:"user_data" p:"user_data" v:"required#参数错误,用户数据不存在"`                         // 主任密钥中的用户数据字串
}

type CloudLeaderBindReq struct {
	//LeaderKeyCode      string `json:"leader_key_code" p:"leader_key_code" v:"required#参数错误：主任秘钥代码必填"`           // 主任密钥代号
	LeaderSerialNumber string `json:"leader_serial_number" p:"leader_serial_number" v:"required#参数错误,主任秘钥系列号必填"` // 主任密钥唯一编号
	AuthorNumber       string `json:"author_number" p:"author_number" v:"required#参数错误：客户端授权码必填"`                // 客户端授权id
	ServerAuthorNumber string `json:"server_author_number" p:"server_author_number" v:"required#参数错误：服务端授权码必填"`  // 服务端授权id
	//Signature          string `json:"signature" p:"signature" v:"required#参数错误：数据签名必填"`                         // 数据签名
}
