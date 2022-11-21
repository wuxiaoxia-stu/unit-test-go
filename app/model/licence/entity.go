package licence

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id           int         `orm:"id,primary,table_comment:'授权记录表'" json:"id"`
	Licence      string      `orm:"licence,size:text,comment:'证书'" json:"licence"`
	LicenceKey   string      `orm:"licence_key,size:text,comment:'公钥'" json:"licence_key"`
	UkeyCode     string      `orm:"ukey_code,size:5,comment:'Ukey代码'" json:"ukey_code"`
	UkeyCrypt    string      `orm:"ukey_crypt,size:500,comment:'ukey硬件ID'" json:"ukey_crypt"`
	AuthorNumber string      `orm:"author_number,size:20,comment:'设备编号'" json:"author_number"`
	DeviceNumber string      `orm:"device_number,size:30,comment:'设备编号'" json:"device_number"`
	CreateAt     *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
	Status       int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'" json:"status" default:"1"`
}

var (
	Table       = "licence"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "licence l"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type QueryHospital struct {
	RegionId []string `p:"region_id" v:"required#参数错误：区域信息必填"`
}

type QueryUkey struct {
	Uuid         string `p:"uuid" v:"required#参数错误：uuid不存在"`
	SerialNumber string `p:"serial_number" v:"required#参数错误：ukey系列号不存在"`
	HospitalId   int    `p:"hospital_id"`
	Signature    string `p:"signature" v:"required#参数错误：签名信息必填"`
}

//云端ukey数据结构
type UkeyData struct {
	Id           int         `orm:"id,primary,size:4,table_comment:'授权Ukey记录'" json:"id"`
	SerialNumber string      `orm:"serial_number,size:20,comment:'硬件ID'" json:"serial_number"`
	Code         string      `orm:"code,size:50,comment:'UKey代码'" json:"code"`
	Type         int         `orm:"type,size:2,comment:'UKey类型'" json:"type"`
	Publickey    string      `orm:"public_key,size:500,comment:'公钥'" json:"public_key"`
	AadminId     int         `orm:"admin_id,size:4,comment:'公钥'" json:"admin_id"`
	AuthTimes    int         `orm:"auth_times,comment:'授权次数'" json:"auth_times"`
	UsedTimes    int         `orm:"used_times,comment:'已使用授权次数'" json:"used_times"`
	CreateAt     *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId   int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status       int         `orm:"status,size:2,not null,default:1" json:"status"`
}

type LicenceCreate struct {
	AddressNumber    string `json:"address_number"`     // 地址编号
	HospitalId       int    `json:"hospital_id"`        // 医院
	UkeySerialNumber string `json:"ukey_serial_number"` // ukey唯一编号
	Signature        string `json:"signature"`          // 数据签名
	State            int    `json:"state"`              // 是否需要新建编号 1-需要， 2-不需要
	Role             int    `json:"role"`               // 客户端角色 （1-windows 客户端，2-web 端）
	Uuid             string `json:"uuid"`
	GpuName          string `json:"gpu_name"`
	CpuName          string `json:"cpu_name"`
	CpuId            string `json:"cpu_id"`
	CpuDeviceId      string `json:"cpu_device_id"`
	BaseboardID      string `json:"baseboard_id"`   // baseboard_id
	DiskID           string `json:"disk_id"`        // 磁盘id
	ClientVersion    string `json:"client_version"` // 版本号
	AiVersion        string `json:"ai_version"`     // ai版本
}

type LicenceCreateRsp struct {
	AuthorNumber string `json:"author_number"`
	Licence      string `json:"licence"`
	LicenceKey   string `json:"licence_key"`
	UkeyCode     string `json:"ukey_code"`
	UkeyCrypt    string `json:"ukey_crypt"`
}

type ActivationReq struct {
	AuthorNumber  string `json:"author_number" p:"author_number" v:"required#参数错误：授权编号必填"`    // 授权编号
	SerialNumber  string `json:"serial_number" p:"serial_number" v:"required#参数错误：设备系列号必填"`   // 设备序列号
	MachineBrand  string `json:"machine_brand" p:"machine_brand" v:"required#参数错误：超声机品牌必填"`   // 超声机品牌
	MachineNumber string `json:"machine_number" p:"machine_number" v:"required#参数错误：超声机信息必填"` // 超声机
	Floor         string `json:"floor" form:"floor" p:"floor" v:"required#参数错误：楼层信息必填"`       // 楼层
	Room          string `json:"room" form:"room" p:"room" v:"required#参数错误：房号信息必填"`          // 房号
}
