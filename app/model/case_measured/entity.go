package case_measured

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id           int         `orm:"id,primary,table_comment:'病例-测量值'" json:"measured_id"`
	SerialNumber string      `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	CaseId       int64       `orm:"case_id,size:8,not null,comment:'病例号'" json:"case_id"`
	Type         int         `orm:"type,size:2,comment:'类型：0：默认，1：心脏'" json:"measured_type"`
	Name         string      `orm:"name,size:50,comment:'测量名称'" json:"measured_name"`
	Value        int         `orm:"value,comment:'测量值'" json:"measured_value"`
	ValuePts     int         `orm:"value_pts,comment:'测量值pts'" json:"measured_value_pts"`
	Min          int         `orm:"min,comment:'最小值'" json:"measured_min"`
	Max          int         `orm:"max,comment:'最大值'" json:"measured_max"`
	Sort         int         `orm:"sort,comment:'排序'" json:"sort"`
	CreateAt     *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
}

var (
	Table       = "case_measured"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_measured cm"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

// 设置测量值
type CaseSetMeasuredReq struct {
	CaseId           int64  `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	MeasuredType     string `v:"required|in:0,1#参数错误：测量类型必填|测量值类型错误"`
	MeasuredName     string `v:"required#参数错误：测量名称必填"`
	MeasuredValue    int    `v:"required#参数错误：测量值必填"`
	MeasuredValuePts int    `default:"-1"`
	MeasuredMin      int    `default:"0"`
	MeasuredMax      int    `default:"0"`
	Sort             int    `default:"0"`
}
