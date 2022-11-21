package kl_syndrome

import (
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_syndrome_feature"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table kl_syndrome.
type Entity struct {
	Id               int                           `orm:"id,primary,size:4,table_comment:'知识图谱-综合征'" json:"id"`
	Uuid             string                        `orm:"uuid,size:50,comment:'UUID'" json:"uuid"`
	Type             int                           `orm:"type,size:2,comment:'类型:1：遗传综合征，2：宫内感染，3：致畸剂'" json:"type"`
	SubType          int                           `orm:"sub_type,size:2,comment:'子类型'" json:"sub_type"`
	Name             string                        `orm:"name,size:100,comment:'特征名称'" json:"name"`
	NameEn           string                        `orm:"name_en,size:100,comment:'特征英文名称'" json:"name_en"`
	GeneLocation     string                        `orm:"gene_location,size:50,comment:'基因点位'" json:"gene_location"`
	GeneLocationEn   string                        `orm:"gene_location_en,size:50,comment:'基因点位'" json:"gene_location_en"`
	GeneticsDesc     string                        `orm:"genetics_desc,size:50,comment:'遗传类型'" json:"genetics_desc"`
	GeneticsDescEn   string                        `orm:"genetics_desc_en,size:50,comment:'遗传类型'" json:"genetics_desc_en"`
	FeatureIds       string                        `orm:"feature_ids,size:text,comment:'特征id'" json:"feature_ids"`
	Diagnose         string                        `orm:"diagnose,size:text,comment:'超声诊断要点'" json:"diagnose"`
	DiagnoseEn       string                        `orm:"diagnose_en,size:text,comment:'超声诊断要点'" json:"diagnose_en"`
	Consult          string                        `orm:"consult,size:text,comment:'预后咨询'" json:"consult"`
	ConsultEn        string                        `orm:"consult_en,size:text,comment:'预后咨询'" json:"consult_en"`
	Other            string                        `orm:"other,size:text,comment:'其他'" json:"other"`
	OtherEn          string                        `orm:"other_en,size:text,comment:'其他'" json:"other_en"`
	Sort             int                           `orm:"sort,size:4,,default:0,comment:'排序'" json:"sort"`
	CreateAt         *gtime.Time                   `orm:"create_at" json:"create_at"`
	UpdateAt         *gtime.Time                   `orm:"update_at" json:"update_at"`
	OperatorId       int                           `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status           int                           `orm:"status,size:2,not null,default:1" json:"status"`
	Children         []*Entity                     `json:"children"`
	Features         []*kl_feature.Entity          `json:"features"`
	SyndromeFeatures []*kl_syndrome_feature.Entity `json:"syndrome_features"`
}

var (
	// Table is the table name of kl_syndrome.
	Table       = "kl_syndrome"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_syndrome ks"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

var TypeTree = []*SyndromeListRsp{
	&SyndromeListRsp{Type: 1, SyndromeName: "遗传综合征", Children: []*SyndromeListRsp{
		&SyndromeListRsp{Type: 1, SubType: 1, SyndromeName: "生长发育受限为特征"},
		&SyndromeListRsp{Type: 1, SubType: 2, SyndromeName: "生长过度为特征"},
		&SyndromeListRsp{Type: 1, SubType: 3, SyndromeName: "颜面部异常为特征"},
		&SyndromeListRsp{Type: 1, SubType: 4, SyndromeName: "大脑异常为特征"},
		&SyndromeListRsp{Type: 1, SubType: 5, SyndromeName: "肢体异常为特征"},
		&SyndromeListRsp{Type: 1, SubType: 6, SyndromeName: "骨骼发育不良为特征"},
		&SyndromeListRsp{Type: 1, SubType: 7, SyndromeName: "颅缝早闭为特征"},
		&SyndromeListRsp{Type: 1, SubType: 8, SyndromeName: "多发异常为特征"},
		&SyndromeListRsp{Type: 1, SubType: 9, SyndromeName: "软组织异常为特征"},
		&SyndromeListRsp{Type: 1, SubType: 10, SyndromeName: "序列征和联合症"},
		&SyndromeListRsp{Type: 1, SubType: 13, SyndromeName: "染色体异常综合征"},
	}},
	&SyndromeListRsp{Type: 2, SyndromeName: "宫内感染"},
	&SyndromeListRsp{Type: 3, SyndromeName: "致畸剂"},
}

// syndrome_list 综合征检索列表
type SyndromeListRsp struct {
	Id               int
	Type             int
	SubType          int
	SyndromeSequence int                `json:"syndrome_sequence"` // 综合征类型key/综合征分组key
	SyndromeUuid     string             `json:"syndrome_uuid"`     // 综合征uuid
	SyndromeName     string             `json:"syndrome_name"`     // 综合征类型名称/综合征分组名称/综合征名称
	SyndromeNameEn   string             `json:"syndrome_name_en"`  // 综合征类型名称/综合征分组名称/综合征名称
	Children         []*SyndromeListRsp `json:"children"`          // 综合征检索列表-数组
}

type FeatureSyndromeReq struct {
	FeatureUuidArray []string `p:"feature_uuid_array"`
}

type SyndromeUuidArrayReq struct {
	SyndromeUuidArray []string `p:"syndrome_uuid_array"`
}
