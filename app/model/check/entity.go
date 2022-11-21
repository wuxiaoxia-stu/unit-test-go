package check

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

type Entity struct {
	Id               string      `orm:"id,size:50,unique,table_comment:'检查项配置表'" json:"check_id"`
	Type             int         `orm:"type,size:2,comment:'使用场景：1：实时分析，2：质控考核'" json:"type"`
	Name             string      `orm:"name,size:100,comment:'检查名称'" json:"check_name"`
	CheckType        int         `orm:"check_type,size:2,comment:'孕期：0 早 1中 2 晚'" json:"check_type"`
	Level            int         `orm:"level,size:2,comment:'检查级别'" json:"check_level"`
	CheckItem        string      `orm:"check_item,size:500,comment:'当前检查切面项'" json:"-"`
	Sort             int         `orm:"sort,size:4,comment:'排序'" json:"sort"`
	UpdateAt         *gtime.Time `orm:"update_at,comment:'修改时间'" json:"-"`
	DeleteAt         *gtime.Time `orm:"delete_at,comment:'删除时间'" json:"-"`
	Status           int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'" json:"status" default:"1"`
	CurrCheckItem    []string    `json:"curr_check_item"`
	DefaultCheckItem []string    `json:"default_check_item"`
}

var (
	Table       = "check"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "check c"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//map[string][]*PartsItem
type PartsItem struct {
	PlaneId     string `json:"plane_id"`
	PlaneNameCh string `json:"plane_name_ch"`
	PlaneNameEn string `json:"plane_name_en"`
	IsCheck     bool   `json:"is_check"`
}

type PartsDefaultCheck struct {
	Name      string   `json:"name"`
	Type      int      `json:"type"`
	Level     int      `json:"level"`
	CheckType int      `json:"check_type"`
	Sort      int      `json:"sort"`
	Default   []string `json:"default"`
}

type SaveCheckOptions struct {
	CheckId       string   `p:"check_id" v:"required#参数错误,检查ID必填"`
	Type          int      `p:"type" v:"required|in:1,2#参数错误,类型值必填|参数错误，类型值取值异常"`
	CheckName     string   `p:"check_name" v:"required#参数错误,检查名称必填|检查名称重复"`
	CheckType     int      `p:"check_type" v:"required|in:0,1,2#参数错误,检查类型必填|参数错误：检查类型取值异常"`
	CheckLevel    int      `p:"check_level" v:"required|in:1,2,3,100,200#参数错误,检查级别必填|参数错误：检查级别取值异常"`
	CurrCheckItem []string `p:"curr_check_item"`
	Sort          int      `p:"sort"`
	Status        int      `p:"status" v:"required|in:0,1#参数错误,状态值必填|参数错误，状态值取值异常"`
}

func init() {
	//自定义验证规则，检查type值是否合法
	if err := gvalid.RegisterRule("check-check-id", CheckerCheckId); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-check-id-type-1", CheckerCheckIdType1); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-check-id-type-2", CheckerCheckIdType2); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-check-name", CheckerCheckName); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查role_id值是否合法
func CheckerCheckId(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
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

//自定义验证规则，检查role_id值是否合法
func CheckerCheckIdType1(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("id", value).Where("type", 1).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查role_id值是否合法
func CheckerCheckIdType2(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("id", value).Where("type", 2).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查role_id值是否合法
func CheckerCheckName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var u *SaveCheckOptions
	if err := gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		return gerror.New("数据解析失败")
	}

	var info *Entity
	where := g.Map{"name": value, "type": u.Type, "status": 1, "id !=": u.CheckId}
	err := M.Where(where).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}
