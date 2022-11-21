package case_patient

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id             int         `orm:"id,primary,table_comment:'病例表'" json:"id"`
	PatientId      string      `orm:"patient_id,size:64,not null,comment:'病人id'" json:"patient_id" form:"patient_id"`            // 病人id(病人表)
	Name           string      `orm:"name,size:50,comment:'病人姓名'" json:"patient_name" form:"patient_name"`                       // 病人姓名
	Age            int         `orm:"age,comment:'病人年龄'" json:"patient_age" form:"patient_age"`                                  // 病人年龄
	Sex            int         `orm:"sex,comment:'病人性别'" json:"patient_sex" form:"patient_sex"`                                  // 病人性别
	AccessionNo    string      `orm:"accession_no,size:150" json:"patient_accession_no" form:"patient_accession_no"`             // ???
	StudyUid       string      `orm:"study_uid,size:150" json:"patient_study_uid" form:"patient_study_uid"`                      // ???
	PatLocalId     string      `orm:"pat_local_id,size:150" json:"patient_patLocal_id" form:"patient_patLocal_id"`               // ???
	LastMensesDate string      `orm:"last_menses_date,size:150" json:"patient_last_menses_date" form:"patient_last_menses_date"` // ???
	CreateAt       *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt       *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
}

var (
	Table       = "case_patient"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_patient cp"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
