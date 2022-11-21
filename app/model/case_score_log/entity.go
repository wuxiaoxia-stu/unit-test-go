package qc_score_log

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var (
	Table       = "qc_score_log"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "qc_score_log qsl"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type Entity struct {
	Id              int64       `orm:"id,primary,size:8,start:10000000,table_comment:'病例分數修改日志'" json:"case_id"`
	Type            int         `orm:"type,size:2,comment:'数据类型：1：病例数据，2：病例切面数据，3：病例结构数据'" json:"serial_number"`
	CaseId          int         `orm:"case_id,comment:'病例id'" json:"case_id"`
	GroupId         string      `orm:"group_id,size:20,comment:'切面组ID'" json:"plane_group_id"`
	PlaneId         string      `orm:"plane_id,size:10,comment:'切面ID'" json:"plane_id"`
	PlaneHash       int64       `orm:"plane_hash,size:8,comment:'切面hash'" json:"plane_hash"`
	PlaneIndex      int         `orm:"plane_index,size:4,comment:'切面索引'" json:"plane_index"`
	PlaneNameCh     string      `orm:"plane_name_ch,size:150,comment:'切面名称'" json:"plane_name_ch"`
	PlaneNameEn     string      `orm:"plane_name_en,size:150,comment:'切面名称'" json:"plane_name_en"`
	StructureId     string      `orm:"structure_id,size:20,comment:'部位hash'" json:"structure_id"`
	StructureHash   int64       `orm:"structure_hash,size:8,comment:'部位hash'" json:"structure_hash"`
	StructureNameCh string      `orm:"structure_name_ch,size:150,comment:'部位名称'" json:"structure_name_ch"`
	StructureNameEn string      `orm:"structure_name_en,size:150,comment:'部位名称'" json:"structure_name_en"`
	Score           int         `orm:"score,size:4,comment:'分数'" json:"score"`
	TotalScore      int         `orm:"total_score,size:4,comment:'总数'" json:"total_score"`
	UserId          int         `orm:"user_id,comment:'用户id（修改分数时必填）'" json:"user_id"`
	QcUserId        int         `orm:"qc_user_id,comment:'QC用户id（远程质控修改分数时必填）'" json:"qc_user_id"`
	Remark          string      `orm:"remark,size:text,comment:'备注'" json:"remark"`
	CreateAt        *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
}
