package meeting_device

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var (
	Table       = "meeting_device"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "meeting_device md"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type Entity struct {
	Id            int         `orm:"id,primary,table_comment:'远程会诊-设备列表'" json:"id"`
	AuthorNumber  string      `orm:"author_number,size:20,comment:'设备授权码'" json:"author_number"`
	SerialNumber  string      `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	DeviceType    int         `orm:"device_type,size:4,comment:'客户端类型，1：客户端，2：专家端'" json:"device_type"`
	UserId        int         `orm:"user_id,size:4,comment:'当前登录用户'" json:"user_id"`
	Online        int         `orm:"online,size:2,comment:'在线状态 1：在线，0：离线'" json:"online"`
	CaseStatus    int         `orm:"case_status,size:2,comment:'病例状态 0：正常，1：异常'" json:"case_status"`
	AssistStatus  int         `orm:"assist_status,size:2,comment:'协助状态 0：不需要协助，1：等待协助', 2: 正在协助" json:"assist_status"`
	CoverImage    string      `orm:"cover,size:150,comment:'封面图'" json:"cover_image"`
	DoctorName    string      `orm:"doctor_name,size:50,comment:'医生名称'" json:"doctor_name"`
	RegionId      string      `orm:"region_id,size:20,comment:'区域ID'" json:"region_id"`
	HospitalId    string      `orm:"hospital_id,size:4,comment:'医院ID'" json:"hospital_id"`
	HospitalName  string      `orm:"hospital_name,size:50,comment:'医院'" json:"hospital_name"`
	Ip            string      `orm:"ip,size:30,comment:'IP'" json:"ip"`
	MachineBrand  string      `orm:"machine_brand,size:50,comment:'超声机型号'" json:"machine_brand"`
	Floor         string      `orm:"floor,size:50,comment:'楼层'" json:"floor"`
	Room          string      `orm:"room,size:50,comment:'房间、科室名称'" json:"room"`
	MachineNumber string      `orm:"machine_number,size:30,comment:'超声机编号'" json:"machine_number"`
	CreateAt      *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt      *gtime.Time `orm:"update_at,comment:'更新时间'" json:"update_at"`
	Avatar        string      `json:"avatar"`
	BaseUrl       string      `json:"base_url"`
}

type DeviceReg struct {
	AuthorNumber  string `v:"required#参数错误,客户端授权号必填"`
	SerialNumber  string `v:"required#参数错误,设备系列号必填"`
	DeviceType    int    `v:"required#参数错误,客户端设备类必填"`
	DoctorName    string `v:"required#参数错误,医生名称必填"`
	RegionId      int    `v:"required#参数错误,所在区域必填"`
	HospitalId    string `v:"required#参数错误,所属医院必填"`
	HospitalName  string `v:"required#参数错误,医院名称必填"`
	CoverImage    string
	MachineBrand  string
	Floor         string
	Room          string
	MachineNumber string
}

type DeviceListReq struct {
	Online       int
	CaseStatus   int
	RegionId     string
	HospitalId   int
	HospitalName string
}
