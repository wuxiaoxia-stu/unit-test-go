package case_label_relate

import (
	"aiyun_local_srv/app/model/case_label"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

type Entity struct {
	Id        int         `orm:"id,primary,table_comment:'病例标签'" json:"id"`
	Userid    int         `orm:"user_id,comment:'用户id'" json:"user_id"`
	CaseId    int64       `orm:"case_id,size:8,comment:'病例ID'" json:"case_id"`
	LabelId   int         `orm:"label_id,comment:'标签ID'" json:"label_id"`
	Sort      int         `orm:"sort,comment:'排序'" json:"sort"`
	CreateAt  *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	DeteteAt  *gtime.Time `orm:"delete_at,comment:'删除时间'" json:"delete_at"`
	LabelName string      `json:"label_name"`
	LabelType string      `json:"label_type"`
}

var (
	Table       = "case_label_relate"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_label_relate clr"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//创建病例提交参数
type CaseAddLabelReq struct {
	CaseId  int64 `v:"required|check-case-id#参数错误：病例ID必填|病例不存在"`
	LabelId int   `v:"required|check-label-id#参数错误：标签ID必填|标签不存在"`
}

func init() {
	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-label-id", CheckerLabelId); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查标签是否存在
func CheckerLabelId(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *case_label.Entity
	err := case_label.M.Where("id", value).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}
