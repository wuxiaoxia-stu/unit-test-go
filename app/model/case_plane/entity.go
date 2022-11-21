package case_plane

import (
	"aiyun_local_srv/app/model/case_plane_structure"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	SerialNumber    string                         `orm:"serial_number,size:20,comment:'设备系列号',table_comment:'病例-切面'" json:"serial_number"`
	CaseId          int64                          `orm:"case_id,size:8,not null,comment:'病例号'" json:"case_id"`
	ShotId          string                         `orm:"shot_id,size:50,comment:'截图id'" json:"shot_id"`
	GroupId         string                         `orm:"group_id,size:20,comment:'切面组ID'" json:"plane_group_id"`
	PlaneId         string                         `orm:"plane_id,size:10,comment:'切面ID'" json:"plane_id"`
	PlaneHash       int64                          `orm:"plane_hash,size:8,comment:'切面hash'" json:"plane_hash"`
	PlaneIndex      int                            `orm:"plane_index,size:4,comment:'切面索引'" json:"plane_index"`
	PlaneNameCh     string                         `orm:"plane_name_ch,size:150,comment:'切面名称'" json:"plane_name_ch"`
	PlaneNameEn     string                         `orm:"plane_name_en,size:150,comment:'切面名称'" json:"plane_name_en"`
	PlaneScore      int                            `orm:"plane_score,size:4,comment:'分数'" json:"plane_score"`
	PlaneTotalScore int                            `orm:"plane_total_score,size:4,comment:'总数'" json:"plane_total_score"`
	Sort            int                            `orm:"sort,comment:'排序'" json:"shot_sort"`
	CreateAt        *gtime.Time                    `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt        *gtime.Time                    `orm:"update_at,comment:'修改时间'" json:"update_at"`
	StructureList   []*case_plane_structure.Entity `json:"structure_list"`
	ShotPath        string                         `json:"shot_path"`
	ShotName        string                         `json:"shot_name"`
}

var (
	Table       = "case_plane"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_plane cp"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type PlaneList struct {
	GroupId         string
	PlaneId         string `v:"required#参数错误：切面ID必填"`
	PlaneNameCh     string `v:"required#参数错误：切面名称必填"`
	PlaneNameEn     string `v:"required#参数错误：切面英文名称必填"`
	PlaneHash       int64  `v:"required#参数错误：切面HASH必填"`
	PlaneScore      int    `v:"required#参数错误：切面得分必填"`
	PlaneScoreTotal int    `v:"required#参数错误：切面总分必填"`
	PlaneIndex      int    `v:"required#参数错误：切面索引"`
	StructureList   []*StructureList
}

type StructureList struct {
	StructureId         string `v:"required#参数错误：部位ID必填"`
	StructureHash       int64  `v:"required#参数错误：部位HASH必填"`
	StructureNameCh     string `v:"required#参数错误：部位名称必填"`
	StructureNameEn     string `v:"required#参数错误：部位英文名称必填"`
	StructureScore      int    `v:"required#参数错误：部位分数必填"`
	StructureScoreTotal int    `v:"required#参数错误：部位总分必填"`
}

type ShotUsedReq struct {
	CaseId        int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	IsReset       bool
	PlaneShotList []*PlaneShotList
}

type PlaneShotList struct {
	ShotId     string `v:"required|check-shot-id-exist#参数错误：截图ID必填|参数错误：截图不存在"`
	PlaneId    string `v:"required#参数错误：切面ID必填"`
	PlaneIndex int
}

type PlaneScoreDetailReq struct {
	CaseId int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
}

type ModifyScoreReq struct {
	CaseId    int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	GroupId   string
	PlaneId   string
	ScoreItem []*ScoreItem
	Remark    string
}

type ScoreItem struct {
	StructureId string
	Score       int
}
