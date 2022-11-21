package qc_upload

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

var UploadStatus = map[int]string{
	1: "已压缩，未上传",
	2: "传输失败, 请重试",
	3: "上传完成",
	4: "下载失败,请重试",
	5: "下载完成",
	6: "病例创建完成",
}

type Entity struct {
	UploadId     string      `orm:"upload_id,size:50,unique,comment:'上传标识',table_comment:'上传记录'" json:"upload_id"`
	Password     string      `orm:"password,size:50,comment:'压缩、解压密码'" json:"password"`
	SerialNumber string      `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	ClientType   int         `orm:"client_type,size:2,comment:'来源，0：客户端上传，1：web端上传'" json:"client_type"`
	SourceType   int         `orm:"source_type,size:2,comment:'资源类型，1：U盘传输，2：病例抽样，3：工作站传输'" json:"source_type"`
	UserId       int         `orm:"user_id,comment:'用户id'" json:"user_id"`
	CheckId      string      `orm:"check_id,size:50,comment:'产筛ID'" json:"check_id"`
	CaseNum      int         `orm:"case_num,size:4,comment:'包含病例数量'" json:"case_num"`
	SuccCaseNum  int         `orm:"succ_case_num,size:4,comment:'创建成功的病例总数'" json:"succ_case_num"`
	TotalScore   int         `orm:"total_score,size:4,comment:'病例总得分'" json:"total_score"`
	PerCaseScore string      `orm:"pre_case_score,size:text,comment:'记录每个病例总分数，由于中位分计算'" json:"pre_case_score"`
	ImgNum       int         `orm:"img_num,size:4,comment:'包含图片数量'" json:"img_num"`
	FiLeSize     int64       `orm:"file_size,size:8,comment:'压缩包大小'" json:"file_size"`
	CreateAt     *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	ZipAt        *gtime.Time `orm:"zip_at,comment:'压缩完成时间'" json:"zip_at"`
	UploadAt     *gtime.Time `orm:"upload_at,comment:'上传完成时间'" json:"upload_at"`
	DownloadAt   *gtime.Time `orm:"download_at,comment:'下载完成时间'" json:"download_at"`
	CaseCreateAt *gtime.Time `orm:"case_create_at,comment:'病例创建时间'" json:"case_create_at"`
	Status       int         `orm:"status,size:2,default:1,comment:'状态，0：未压缩，1：已压缩，未上传 2：传输失败，请重试 3：上传完成，4：下载失败，请重试 5：下载完成 6：病例创建完成'" json:"status"`
	StatusText   string      `json:"status_text"`
	CheckName    string      `json:"check_name"`
	CheckType    int         `json:"check_type"`
	CheckLevel   int         `json:"check_level"`
}

var (
	Table       = "qc_upload"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "qc_upload qu"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

func init() {
	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-qc-upload-id", CheckerQcUploadId); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查case_id值是否合法
func CheckerQcUploadId(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("upload_id", value).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

type CreateUploadReq struct {
	SerialNumber string
	ClientType   int    `v:"required|in:0,1#客户端类型必填|客户端类型异常"`
	SourceType   int    `v:"required|in:1,2,3#资源来源必填|资源类型异常"`
	CheckId      string `v:"required|check-check-id-type-2#产筛ID必填|产筛不存在"`
}

//上传完成
type SetUploadStatusReq struct {
	UploadId string `v:"required|check-qc-upload-id#文件标识必填|上传ID异常"`
	Status   int    `v:"required|in:1,2,3,4,5,6#状态值必填|参数错误：状态值异常"`
	CaseNum  int    `v:"integer|min:0|max:200#病例数量为正整数|单个压缩包允许病例数量为1-200|单个压缩包允许病例数量为1-200"`
	ImgNum   int    `v:"integer|min:0#文件数量为正整数|文件数量为正整数"`
	FileSize int64  `v:"integer|min:0#文件大小必须为正整数|文件大小必须为正整数"`
}

type PageReqParams struct {
	KeyWord      string `p:"key_word"`
	SerialNumber string `p:"serial_number"`
	UserId       int    `p:"user_id"`
	Status       int    `p:"status" default:"-1"`
	Page         int    `p:"page" default:"1"`
	PageSize     int    `p:"page_size" default:"10"`
	Order        string `p:"order" default:"create_at"`
	Sort         string `p:"sort" default:"DESC"`
	TimeBegin    string
	TimeEnd      string
}

type UploadDeleteReq struct {
	UploadIds []string `json:"upload_ids"`
}

type QcTimelineReq struct {
	TimeBegin  string
	TimeEnd    string
	CheckId    string
	SourceType int
	Dates      []string
}

type TimeLine struct {
	Date string `json:"date"`
}

type Year struct {
	Year     string   `json:"year"`
	Children []*Month `json:"children"`
}

type Month struct {
	Month    string   `json:"month"`
	Children []string `json:"children"`
}

type QcCaseList struct {
	CaseId     int64
	DoctorName string
	Score      int
	CaseTime   string
}

type QcGroupListReq struct {
	RegionId   string
	HospitalId int
	TimeBegin  string
	TimeEnd    string
}

type QcGroupList struct {
	CaseCreateAt     *gtime.Time               `json:"case_create_at"`
	UserId           int                       `json:"user_id"`
	Realname         string                    `json:"realname"`
	SourceType       int                       `json:"source_type"`
	CheckId          string                    `json:"check_id"`
	CheckName        string                    `json:"check_name"`
	RegionId         string                    `json:"region_id"`
	RegionName       string                    `json:"region_name"`
	HospitalId       int                       `json:"hospital_id"`
	HospitalName     string                    `json:"hospital_name"`
	TotalNum         int                       `json:"total_num"`
	TotalCaseNum     int                       `json:"total_case_num"`
	TotalAveragScore int                       `json:"total_averag_score"`
	TotalMedianScore int                       `json:"total_median_score"`
	DateGroupList    map[string][]*QcGroupList `json:"date_group_list"`
}
