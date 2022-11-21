package case_label

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

type Entity struct {
	Id     int `orm:"id,primary,table_comment:'标签管理'" json:"label_id"`
	UserId int `orm:"user_id,comment:'用户id'" json:"user_id"`
	Type   int `orm:"type,comment:'标签类型 1：系统标签，2：自定义标签'" json:"label_type"`
	//Key      string      `orm:"key,size:50,comment:'标签key'" json:"label_key"`
	GroupName string      `orm:"group_name,size:50,comment:'标签分组'" json:"label_group_name"`
	Name      string      `orm:"name,size:50,comment:'标签名称'" json:"label_name"`
	Sort      int         `orm:"sort,comment:'排序'" json:"sort"`
	CreateAt  *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	Status    int         `orm:"status,size:2,default:1,comment:'状态，1：有效，0：删除'" json:"status" default:"1"`
	CaseId    int64       `json:"case_id"`
}

var (
	Table       = "case_label"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "case_label cl"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type GroupLabel struct {
	GroupName string       `json:"group_name"`
	LabelList []*LabelList `json:"label_list"`
}

type LabelList struct {
	Id     int    `json:"label_id"`
	Name   string `json:"label_name"`
	Type   int    `json:"label_type"`
	CaseId int64  `json:"case_id"`
}

//创建病例提交参数
type CreateLabelReq struct {
	LabelName string `v:"required|check-label-name#标签名称必填|标签已存在"`
}

func init() {
	//自定义验证规则，检查病例id值是否合法
	if err := gvalid.RegisterRule("check-label-name", CheckerLabelName); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查标签名称是否存在
func CheckerLabelName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("name", value).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}
