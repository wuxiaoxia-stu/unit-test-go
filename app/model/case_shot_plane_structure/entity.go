package case_shot_plane_structure

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var (
	Table       = "case_shot_plane_structure"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_shot_plane_structure csps"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type Entity struct {
	SerialNumber        string      `orm:"serial_number,size:20,comment:'设备系列号',table_comment:'病例-截图-切面-结构-分数'" json:"serial_number"`
	CaseId              int64       `orm:"case_id,size:8,comment:'病例ID'" json:"case_id"`
	ShotId              string      `orm:"shot_id,size:50,comment:'截图ID'" json:"shot_id"`
	PlaneId             string      `orm:"plane_id,size:10,comment:'切面ID'" json:"plane_id"`
	PlaneHash           int64       `orm:"plane_hash,size:8,comment:'切面hash'" json:"plane_hash"`
	StructureId         string      `orm:"structure_id,size:20,comment:'部位hash'" json:"structure_id"`
	StructureHash       int64       `orm:"structure_hash,size:8,comment:'部位hash'" json:"structure_hash"`
	StructureNameCh     string      `orm:"structure_name_ch,size:150,comment:'部位名称'" json:"structure_name_ch"`
	StructureNameEn     string      `orm:"structure_name_en,size:150,comment:'部位名称'" json:"structure_name_en"`
	StructureScore      int         `orm:"structure_score,size:4,comment:'部位得分'" json:"structure_score"`
	StructureScoreTotal int         `orm:"structure_score_total,size:4,comment:'部位分数限制'" json:"structure_score_total"`
	CreateAt            *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt            *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
}
