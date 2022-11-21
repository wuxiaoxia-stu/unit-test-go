package case_info

import (
	"aiyun_local_srv/app/model/case_measured"
	"aiyun_local_srv/app/model/case_plane"
	"aiyun_local_srv/app/model/case_shot"
	"context"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

type Entity struct {
	Id             int64                   `orm:"id,primary,size:8,start:10000000,table_comment:'病例表'" json:"case_id"`
	SerialNumber   string                  `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	UserId         int                     `orm:"user_id,comment:'用户id'" json:"user_id"`
	Type           int                     `orm:"type,comment:'病例类型 1：实时分析病例 2：单人质控  3：多人质控 4：远程质控病例'" json:"case_type"`
	MultiGroupId   string                  `orm:"multi_group_id,size:50,comment:'分组标识'" json:"multi_group_id"`
	MultiIndex     int                     `orm:"multi_index,size:4,default:0,comment:'分组索引'" json:"multi_index"`
	CheckId        string                  `orm:"check_id,size:50,comment:'检查ID'" json:"check_id"`
	CheckType      int                     `orm:"check_type,comment:'检查类型 0 早 1中 2 晚'" json:"check_type"`
	CheckLevel     int                     `orm:"check_level,comment:'检查级别'" json:"check_level"`
	DoctorName     string                  `orm:"doctor_name,size:100,comment:'医生名称'" json:"doctor_name"`
	DepartmentName string                  `orm:"department_name,size:100,comment:'科室名称'" json:"department_name"`
	QcCaseDate     string                  `orm:"qc_case_date,size:20,comment:'质控病例时间'" json:"qc_case_date"`
	Path           string                  `orm:"path,size:150,comment:'病例路径'" json:"case_path"`
	PathCloud      string                  `orm:"path_cloud,size:150,comment:'病例云路径'" json:"case_path_cloud"`
	Length         int                     `orm:"length,comment:'病例时长',default:0" json:"case_length"`
	FetusName      string                  `orm:"fetus_name,size:20,comment:'胎儿名称 F1 F2 F3'" json:"fetus_name"`
	FetusNum       int                     `orm:"fetus_num,size:2,comment:'胎儿数量',default:1" json:"fetus_num"`
	JoinFlag       int                     `orm:"join_flag,size:2,comment:'续接标志',default:0" json:"join_flag"`
	PatientId      string                  `orm:"patient_id,size:50,comment:'患者ID'" json:"patient_id"`
	ResultStatus   int                     `orm:"result_status,default:1,comment:'病例结果切面状态 1 正常 2 疑似 3 未检， 4 异常'" json:"case_result"`
	Score          int                     `orm:"score,comment:'病例得分',default:0" json:"case_score"`
	TotalScore     int                     `orm:"total_score,comment:'病例总分',default:0" json:"case_total_score"`
	FlagCount      int                     `orm:"flag_count,comment:'标签数量',default:0" json:"flag_count"`
	AbnormalCount  int                     `orm:"abnormal_count,comment:'非正常标签数量',default:0" json:"abnormal_count"`
	GaAua          int32                   `orm:"ga_aua,comment:'AUA'" json:"ga_aua"`
	GaLmp          int32                   `orm:"ga_lmp,comment:'LMP'" json:"ga_lmp"`
	LmpDate        int64                   `orm:"lmp_date,size:8,comment:'孕期'" json:"lmp_date"`
	StopAt         *gtime.Time             `orm:"stop_at,comment:'结束时间'" json:"stop_at"`
	CreateAt       *gtime.Time             `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt       *gtime.Time             `orm:"update_at,comment:'修改时间'" json:"update_at"`
	HashId         string                  `orm:"hash_id,size:50,comment:'hash_id'" json:"-"`
	UploadId       string                  `orm:"upload_id,size:50,comment:'上传ID'" json:"upload_id"`
	IsHidden       int                     `orm:"is_hidden,comment:'是否隐藏，0：不隐藏，1：隐藏，2，永久隐藏',default:0" json:"-"`
	StopTime       int64                   `json:"stop_time"`
	CreateTime     int64                   `json:"create_time"`
	UpdateTime     int64                   `json:"update_time"`
	CheckName      string                  `json:"check_name"`          // 检查名称
	Username       string                  `json:"username"`            // 用户名
	PatientName    string                  `json:"patient_name"`        // 病人年龄
	PatientAge     int                     `json:"patient_age"`         // 病人年龄
	PatientSex     int                     `json:"patient_sex"`         // 病人性别
	AccessionNo    string                  `json:"patient_accessionNo"` // ???
	StudyUid       string                  `json:"patient_studyUID"`    // ???
	PatLoaclId     string                  `json:"patient_patLocalId"`  // ???
	LastMensesDate string                  `json:"patient_lastMensesDate"`
	LabelList      []*CaseLabel            `json:"label_list"`
	MeasuredList   []*case_measured.Entity `json:"measured_list"`
	ShotList       []*case_shot.Entity     `json:"shot_list"`
	AutoShotCount  int                     `json:"auto_shot_count"`
}

type CaseLabel struct {
	LabelId   int    `json:"label_id"`
	LabelName string `json:"label_name"`
	LabelType string `json:"label_type"`
}

var (
	Table       = "case_info"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_info ci"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//创建病例提交参数
type CreateCaseReq struct {
	CaseType     int    `v:"required|in:1,2,3,4#参数错误：病例类型必填|參數错误：类型值异常"`
	CheckId      string `v:"required#参数错误：检查ID必填"`
	MultiGroupId string
	MultiIndex   int
}

type SaveScoreStandardReq struct {
	CaseId    int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	PlaneList []*case_plane.PlaneList
}

//病例更新接口
type UpdateCaseReq struct {
	CaseId                int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	DoctorName            string
	QcCaseDate            string
	DepartmentName        string
	CasePath              string
	CasePathCloud         string
	FetusName             string
	JoinFlag              int    `default:"-1"`
	CaseLength            int    `default:"-1"`
	PatientId             string `p:"patient_id"`
	PatientName           string `p:"patient_name"` // 病人年龄
	PatientAge            int    `p:"patient_age"`  // 病人年龄
	PatientSex            int    `p:"patient_sex"`  // 病人性别
	PatientAccessionNo    string `p:"patient_accession_no"`
	PatientStudyUid       string `p:"patient_study_uid"`
	PatientPatLoaclId     string `p:"patient_pat_loacl_id"`
	PatientLastMensesDate string `p:"patient_last_menses_date"`
	ResultStatus          int    `p:"case_result_status" default:"-1"`
	GaAua                 int32  `p:"ga_aua" default:"-1"`
	GaLmp                 int32  `p:"ga_lmp" default:"-1"`
	LmpDate               int64  `p:"lmp_date" default:"-1"`
}

//病例结束提交参数
type FinishCaseReq struct {
	CaseId     int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	CaseLength int   `v:"required#参数错误：病例时长必填"`
	CaseResult int   `p:"case_result" v:"in:1,2,3,4#参数错误：病例結果取值异常" default:"-1"`
}

type TimeLine struct {
	Date string `json:"date"`
}

type Year struct {
	Year     string   `json:"year"`
	Children []*Month `json:"children_item"`
}

type Month struct {
	Month    string `json:"month"`
	Children []*Day `json:"children_item"`
}

type Day struct {
	Day string `json:"day"`
}

//type Day struct {
//	Day string `json:"day"`
//}

type CaseTimelineReq struct {
	KeyWords      string
	CaseType      string
	CheckType     string
	CheckLevel    string
	DoctorName    string
	PatientName   string
	TimeBegin     string
	TimeEnd       string
	JoinFlag      int `default:"-1"`
	CaseResult    int `default:"-1"`
	CaseLengthMin int `default:"-1"`
	CaseLengthMax int `default:"-1"`
	AuaWeekBegin  int `default:"-1"`
	AuaWeekEnd    int `default:"-1"`
	LmpWeekBegin  int `default:"-1"`
	LmpWeekEnd    int `default:"-1"`
	CaseLabels    []string
	UserIds       []int
}

type CaseListReq struct {
	Dates []string
}

type CaseListJoinReq struct {
	Date       string
	JoinFlag   bool
	ExistLabel bool
}

func init() {
	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-case-id", CheckerCaseId); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查case_id值是否合法
func CheckerCaseId(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("id", value).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

type PageReqParams struct {
	KeyWord       string `p:"keyword"`
	Status        int    `p:"status" default:"-1"`
	Page          int    `p:"page" default:"1"`
	PageSize      int    `p:"page_size" default:"10"`
	Order         string `p:"order" default:"id"`
	Sort          string `p:"sort" default:"DESC"`
	CaseId        int64  `default:"-1"`
	StartTime     string `p:"start_time"`
	EndTime       string `p:"end_time"`
	JoinFlag      int    `default:"-1"`
	CaseType      string
	CheckType     string
	CheckLevel    string
	DoctorName    string
	PatientName   string
	TimeBegin     string
	TimeEnd       string
	CaseResult    int `default:"-1"`
	CaseLengthMin int `default:"-1"`
	CaseLengthMax int `default:"-1"`
	AuaWeekBegin  int `default:"-1"`
	AuaWeekEnd    int `default:"-1"`
	LmpWeekBegin  int `default:"-1"`
	LmpWeekEnd    int `default:"-1"`
	CaseLabels    []string
	UserIds       []int
}

type DelReq struct {
	Ids      []int64 `v:"required#请选择至少一个病例"`
	Password string  `v:"required#密码必填"`
}

type CaseHiddenStatusReq struct {
	CaseIds []int64 `v:"required#请选择至少一个病例"`
}

type HiddenListReq struct {
	KeyWord   string `p:"keyword"`
	Status    int    `p:"status" default:"-1"`
	Page      int    `p:"page" default:"1"`
	PageSize  int    `p:"page_size" default:"10"`
	Order     string `p:"order" default:"id"`
	Sort      string `p:"sort" default:"DESC"`
	CaseId    int64  `default:"-1"`
	StartTime string `p:"start_time"`
	EndTime   string `p:"end_time"`
	JoinFlag  int    `default:"-1"`
}
