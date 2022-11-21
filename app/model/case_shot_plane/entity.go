package case_shot_plane

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

var (
	Table       = "case_shot_plane"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_shot_plane csp"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type Entity struct {
	SerialNumber    string      `orm:"serial_number,size:20,comment:'设备系列号',table_comment:'病例-截图-切面-分数'" json:"serial_number"`
	CaseId          int64       `orm:"case_id,size:8,comment:'病例ID'" json:"case_id"`
	ShotId          string      `orm:"shot_id,size:50,comment:'截图ID'" json:"shot_id"`
	GroupId         string      `orm:"group_id,size:20,comment:'切面组ID'" json:"plane_group_id"`
	PlaneId         string      `orm:"plane_id,size:10,comment:'切面ID'" json:"plane_id"`
	PlaneHash       int64       `orm:"plane_hash,size:8,comment:'切面hash'" json:"plane_hash"`
	PlaneNameCh     string      `orm:"plane_name_ch,size:150,comment:'切面名称'" json:"plane_name_ch"`
	PlaneNameEn     string      `orm:"plane_name_en,size:150,comment:'切面名称'" json:"plane_name_en"`
	PlaneScore      int         `orm:"plane_score,size:4,comment:'切面得分'" json:"plane_score"`
	PlaneScoreTotal int         `orm:"plane_score_total,size:4,comment:'切面分数限制'" json:"plane_score_total"`
	CreateAt        *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt        *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
}

type ShotScoreSaveReq struct {
	CaseId    int64  `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	ShotId    string `v:"required|check-shot-id-exist#参数错误：截图ID必填|参数错误：截图不存在"`
	PlaneList []*PlaneList
}

type PlaneList struct {
	GroupId         string
	PlaneId         string `v:"required#参数错误：切面ID必填"`
	PlaneNameCh     string `v:"required#参数错误：切面名称必填"`
	PlaneNameEn     string `v:"required#参数错误：切面英文名称必填"`
	PlaneHash       int64  `v:"required#参数错误：切面HASH必填"`
	PlaneScore      int    `v:"required#参数错误：切面得分必填"`
	PlaneScoreTotal int    `v:"required#参数错误：切面总分必填"`
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
	CaseId     int64  `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	ShotId     string `v:"required|check-shot-id-exist#参数错误：截图ID必填|参数错误：截图不存在"`
	OldPlaneId string
	NewPlaneId string `v:"required|check-plane-is-used#参数错误：切面ID必填|参数错误：此切面已经被使用"`
}

func init() {
	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-shot-id-exist", CheckerShotIdExist); err != nil {
		panic(err)
	}

	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-plane-is-used", CheckerPlaneIsUsed); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查shot_id值是否合法
func CheckerShotIdExist(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("shot_id", value).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查切面是够被图片使用
func CheckerPlaneIsUsed(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var d *ShotUsedReq
	if err := gconv.Struct(data, &d); err != nil {
		g.Log().Error(err.Error())
		return err
	}

	var info *Entity
	err := M.Where(g.Map{"case_id": d.CaseId, "plane_id": d.NewPlaneId, "is_used": 1}).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}
