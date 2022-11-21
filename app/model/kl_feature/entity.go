package kl_feature

import (
	"aiyun_local_srv/app/model/kl_feature_atlas"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table kl_feature.
type Entity struct {
	Id         int                        `orm:"id,primary,size:4,table_comment:'知识图谱-特征'" json:"id"`
	Pid        int                        `orm:"pid,size:4,comment:'上级ID'" json:"pid"`
	Uuid       string                     `orm:"uuid,size:50,comment:'UUID'" json:"uuid"`
	Level      int                        `orm:"level,size:2,comment:'层级'" json:"level"`
	Name       string                     `orm:"name,size:100,comment:'特征名称'" json:"name"`
	NameEn     string                     `orm:"name_en,size:100,comment:'特征英文名称'" json:"name_en"`
	Invisible  int                        `orm:"invisible,size:2,default:0,comment:'超声不可见病理特征'" json:"invisible"`
	Define     string                     `orm:"define,size:text,comment:'定义'" json:"define"`
	DefineEn   string                     `orm:"define_en,size:text,comment:'定义'" json:"define_en"`
	Diagnose   string                     `orm:"diagnose,size:text,comment:'超声诊断要点'" json:"diagnose"`
	DiagnoseEn string                     `orm:"diagnose_en,size:text,comment:'超声诊断要点'" json:"diagnose_en"`
	Consult    string                     `orm:"consult,size:text,comment:'预后咨询'" json:"consult"`
	ConsultEn  string                     `orm:"consult_en,size:text,comment:'超声诊断要点'" json:"consult_en"`
	Other      string                     `orm:"other,size:text,comment:'其他'" json:"other"`
	OtherEn    string                     `orm:"other_en,size:text,comment:'其他'" json:"other_en"`
	Sort       int                        `orm:"sort,size:4,,default:0,comment:'排序'" json:"sort"`
	CreateAt   *gtime.Time                `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time                `orm:"update_at" json:"update_at"`
	OperatorId int                        `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status     int                        `orm:"status,size:2,not null,default:1" json:"status"`
	Children   []*Entity                  `json:"children"`
	Atlas      []*kl_feature_atlas.Entity `json:"atlas"`
	Type       string                     `json:"type"` // 1:特征病变 2:其它可见病变
	RootId     int                        `json:"root_id"`
	RootPath   string                     `json:"root_path"`
	IsCheck    bool                       `json:"is_check"` // 选中
}

var (
	// Table is the table name of kl_feature.
	Table       = "kl_feature"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_feature kf"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type FeatureListRsp struct {
	Id          int               `json:"id"`
	Pid         int               `json:"pid"`
	Uuid        string            `json:"uuid"`         // 系统-部位-特征分组-特征 uuid
	Serial      string            `json:"serial"`       // 系统-部位-特征分组-特征 列号
	Name        string            `json:"name"`         // 系统-部位-特征分组-特征 名称
	NameEn      string            `json:"name_en"`      // 系统-部位-特征分组-特征 名称（英文）
	ContentType int32             `json:"content_type"` // 内容类型；1-系统，2-部位，3-特征分组，4-特征
	Children    []*FeatureListRsp `json:"children"`     // 部位或特征详情-数组
	NotClassify []*FeatureListRsp `json:"not_classify"` // 未分类的特征(数组)
}
