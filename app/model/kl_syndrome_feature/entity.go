package kl_syndrome_feature

import (
	"github.com/gogf/gf/frame/g"
)

// Entity is the golang structure for table kl_syndrome_feature.
type Entity struct {
	Id              int    `orm:"id,primary,size:4,table_comment:'知识图谱-综合征-特征-关系表'" json:"id"`
	SyndromeId      int    `orm:"syndrome_id,comment:'综合征ID'" json:"syndrome_id"`
	FeatureId       int    `orm:"feature_id,comment:'特征ID'" json:"feature_id"`
	Type            string `orm:"type,size:5,comment:'类型 1：特征病变，2：其它可见病变'" json:"type"`
	FeatureRootId   int    `orm:"feature_root_id,comment:'特征顶级节点ID'" json:"feature_root_id"`
	FeatureIdPath   string `orm:"feature_id_path,size:30,comment:'特征节点路径'" json:"feature_id_path"`
	FeatureName     string `json:"feature_name"`
	FeatureNameEn   string `json:"feature_name_en"`
	FeatureRootName string `json:"feature_root_name"`
}

var (
	// Table is the table name of kl_syndrome_feature.
	Table       = "kl_syndrome_feature"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_syndrome_feature ksf"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type SyndromeCount struct {
	SyndromeId int `json:"syndrome_id"`
	Count      int `json:"count"`
}

// feature_syndrome 特征关联的综合征
type FeatureSyndromeReq struct {
	SyndromeUuid    string     `json:"syndrome_uuid"`      // 综合征uuid
	SyndromeSerial  string     `json:"syndrome_serial"`    // 综合征编号名称
	SyndromeName    string     `json:"syndrome_name"`      // 综合征名称
	SyndromeNameEn  string     `json:"syndrome_name_en"`   // 综合征名称(英文)
	SyndromeQuality int        `json:"syndrome_quality"`   // 综合征质量
	HitGroupCount   int        `json:"hit_group_count"`    // qt特征检索综合征，综合征命中多少个作为检索条件的超声特征
	Features        []*Feature `json:"features,omitempty"` // 特征病变(形态学部位特征信息-数组)
	Others          []*Feature `json:"others,omitempty"`   // 其他特征病变(形态学部位特征信息-数组)
	Type1Count      int        `json:"-"`                  // 汇总特征病变的总数量，用以排序
}

// kl_syndrome_morphology_feature 形态学部位特征信息
type Feature struct {
	PartUuid           string `json:"part_uuid"`            // 特征部位系统uuid
	PartSerial         string `json:"part_serial"`          // 特征部位系统编号名称
	PartName           string `json:"part_name"`            // 部位名称（仅赋值：有序得形态学列表）
	FeatureUuid        string `json:"feature_uuid"`         // 特征uuid
	FeatureSerial      string `json:"feature_serial"`       // 特征编号
	FeatureName        string `json:"feature_name"`         // 特征名称
	SyndromeFeatureOpt string `json:"syndrome_feature_opt"` // 形态学特征常见状态；SY01-特征病变，SY02-其他可见病变
	PartNameEn         string `json:"part_name_en"`         // 部位名称(英文)（仅赋值：有序得形态学列表）
	FeatureNameEn      string `json:"feature_name_en"`      // 特征名称(英文)
}
