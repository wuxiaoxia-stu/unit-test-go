package case_shot

import (
	"aiyun_local_srv/app/model/case_plane"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	ShotId       string               `orm:"shot_id,size:50,unique,table_comment:'病例-截图'" json:"shot_id"`
	SerialNumber string               `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	CaseId       int64                `orm:"case_id,size:8,not null,comment:'病例号'" json:"case_id"`
	Type         int                  `orm:"type,size:2,comment:'类型，1：ai截图 2：用户截图 3：用户导入'" json:"shot_type"`
	Name         string               `orm:"name,size:100,comment:'截图名称'" json:"shot_name"`
	Path         string               `orm:"path,size:250,comment:'路径'" json:"shot_path"`
	Pts          int64                `orm:"pts,size:8,comment:'帧数值'" json:"shot_pts"`
	PlaneNum     int                  `orm:"plane_num,size:4,comment:'包含切面个数'" json:"plane_num"`
	UsedTimes    int                  `orm:"used_times,size:8,comment:'被使用次数'" json:"used_times"`
	Sort         int                  `orm:"sort,comment:'排序'" json:"shot_sort"`
	CreateAt     *gtime.Time          `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt     *gtime.Time          `orm:"update_at,comment:'修改时间'" json:"update_at"`
	PlaneList    []*case_plane.Entity `json:"relate_plane_list"` // 截图关联的切面列表
}

var (
	Table       = "case_shot"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_shot cs"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

// 保存截图
type CaseSaveShotReq struct {
	CaseId    int64  `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	ShotId    string `v:"required#截图ID必填"`
	ShotType  int    `v:"required|in:1,2,3#参数错误：截图类型必填|截图类型取值异常"`
	ShotName  string `v:"required#截图名称必填"`
	ShotPath  string `v:"required#截图路径必填"`
	ShotPts   int64
	ShotSort  int `default:"0"`
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

type CaseDeleteShotReq struct {
	ShotIds []string `v:"required#请选择至少选择一张图片"`
}
