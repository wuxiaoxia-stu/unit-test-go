package pre_pair

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id               int         `orm:"id,primary,table_comment:'已配对记录'" json:"id"`
	UkeyCode         string      `orm:"ukey_code,size:5,comment:'房间'" json:"ukey_code"`
	AuthorNumber     string      `orm:"author_number,size:20,comment:'设备编号'" json:"author_number"`
	DeviceNumber     string      `orm:"device_number,size:30,comment:'设备编号'" json:"device_number"`
	Uuid             string      `orm:"uuid,size:40,comment:'UUID'" json:"uuid"`
	PublicKey        string      `orm:"public_key,size:1024,comment:'公钥'" json:"public_key"` // 授权公钥
	Ip               string      `orm:"ip,size:20,comment:'医院名称'" json:"ip"`
	Floor            string      `orm:"floor,size:10,comment:'楼层'" json:"floor"`
	Room             string      `orm:"room,size:10,comment:'房间'" json:"room"`
	MachineType      string      `orm:"machine_type,size:10,comment:'超声机型号'" json:"machine_type"`
	ClientVersion    string      `orm:"client_version,size:30,comment:'客户端版本'" json:"client_version"`
	AlgorithmVersion string      `orm:"algorithm_version,size:30,comment:'算法版本'" json:"algorithm_version"`
	CreateAt         *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt         *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
	Status           int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'" json:"status" default:"1"`
}

var (
	Table       = "pre_pair"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "pre_pair pp"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//客户端配对提交数据
type PairClientReq struct {
	AuthorNumber string `p:"author_number" v:"required#参数错误,客户端授权编号必填"` // 客户端授权编号
	SerialNumber string `p:"serial_number" v:"required#参数错误,设备序列号必填"`   // 客户端设备系列号
	//Ciphertext   string `p:"ciphertext" v:"required#参数错误,加密信息必填"`			 // 密文
}

type BreakPairClientReq struct {
	AuthorNumber string `p:"author_number" v:"required#参数错误,客户端授权编号必填"` // 客户端授权编号
}

type PairReq struct {
	ClientAuthorNumber string `p:"client_author_number" v:"required#参数错误,客户端授权编号必填"`  // 客户端授权编号
	ClientSerialNumber string `p:"client_serial_number" v:"required#参数错误,客户端设备系列号必填"` // 客户端设备系列号
	ServerAuthorNumber string `p:"service_author_number" v:"required#参数错误,服务端授权编号必填"` // 服务端授权编号
	//ServerSerialNumber string `p:"service_serial_number" v:"required#参数错误,服务端设备系列号必填"` // 服务端授权编号
	Signature string `p:"signature" v:"required#参数错误,服务端签名信息必填"` // 服务端数据签名
}
