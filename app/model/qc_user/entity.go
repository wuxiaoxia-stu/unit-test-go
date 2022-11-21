package qc_user

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/library/utils/cache"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

var QC_USER_ROLE_TYPE = map[int]string{
	1: "普通用户",
	3: "考官",
	5: "专家",
	7: "审核员",
	9: "管理员",
}

type Entity struct {
	Id              int         `orm:"id,primary,table_comment:'远程质控用户表'" json:"user_id"`
	Username        string      `orm:"username,size: 30,comment:'用户名'" json:"username"`
	Password        string      `orm:"password,size:50,comment:'密码'" json:"-"`
	RoleType        int         `orm:"role_type,default:1,comment:'角色类型'" json:"role_type"`
	IdNo            string      `orm:"id_no,size:20,comment:'身份证号'" json:"id_no"`
	RealName        string      `orm:"real_name,size:50,comment:'真实姓名'" json:"real_name"`
	Avatar          string      `orm:"avatar,size:150,comment:'头像'" json:"avatar"`
	RegionId        string      `orm:"region_id,size:10,comment:'医院所在区域'" json:"region_id"`
	RegionName      string      `orm:"region_name,size:50,comment:'区域名称'" json:"region_name"`
	HospitalId      int         `orm:"hospital_id,size:4,comment:'医院'" json:"hospital_id"`
	HospitalName    string      `orm:"hospital_name,size:50,comment:'医院名称'" json:"hospital_name"`
	PositionalTitle string      `orm:"positional_title,size:10,comment:'职称'" json:"positional_title"`
	CreateAt        *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	Status          int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'"`
	BaseUrl         string      `json:"base_url"`
}

var (
	Table       = "qc_user"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "qc_user u"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

func init() {
	//自定义验证规则，检查用户名是否合法
	if err := gvalid.RegisterRule("check-qc-username-exist", CheckerQcUsernameExist); err != nil {
		panic(err)
	}

	//自定义验证规则，检查用户名是否存在
	if err := gvalid.RegisterRule("check-qc-username-not-exist", CheckerQcUsernameNotExist); err != nil {
		panic(err)
	}

	// 检查手机号是否已经被使用
	if err := gvalid.RegisterRule("check-qc-phone-exist", CheckerQcPhoneExist); err != nil {
		panic(err)
	}

	//自定义验证规则，检查验证码是够正确
	if err := gvalid.RegisterRule("check-captcha-code-correct", CheckerCaptchaCodeCorrect); err != nil {
		panic(err)
	}

	//自定义验证规则，检查验证码是够正确
	if err := gvalid.RegisterRule("check-phone-code-correct", CheckerPhoneCodeCorrect); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查用户名是否存在
func CheckerQcUsernameExist(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("username", value).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查用户名是否合法
func CheckerQcUsernameNotExist(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("username", value).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查手机号是否 被注册
func CheckerQcPhoneExist(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var d *GetSms
	if err := gconv.Struct(data, &d); err != nil {
		return err
	}

	var info *Entity
	err := M.Where("username", value).Scan(&info)
	if err != nil {
		return err
	}

	if (d.Type == "register" || d.Type == "change_username") && info != nil {
		return gerror.New("该手机号已被注册")
	}
	if (d.Type == "login" || d.Type == "find_pwd") && info == nil {
		return gerror.New("该手机号未注册")
	}
	return nil
}

//自定义验证规则，检查验证码是够正确
func CheckerCaptchaCodeCorrect(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var d *GetSms
	if err := gconv.Struct(data, &d); err != nil {
		return err
	}

	//code, err := cache.Get(_const.CAPTCHA_CODE_CACHE_KEY(d.CaptchaId))
	//if err != nil {
	//	return err
	//}
	//
	//if code != d.CaptchaCode {
	//	return gerror.New(message)
	//}
	//
	//if err := cache.Del(_const.CAPTCHA_CODE_CACHE_KEY(d.CaptchaId)); err != nil {
	//	return err
	//}

	return nil
}

//自定义验证规则，检查验证码是够正确
func CheckerPhoneCodeCorrect(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var d *Register
	if err := gconv.Struct(data, &d); err != nil {
		return err
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("register", d.Username))
	if err != nil {
		return err
	}

	if code != d.Code {
		return gerror.New(message)
	}

	return nil
}

type Login struct {
	Username string `v:"required#用户名必填"`
	Password string `v:"required#密码必填"`
}

type SmsLogin struct {
	Username string `v:"required|check-qc-username-not-exist#用户名必填|用户不存在"`
	Code     string `v:"required#短信验证码必填"`
}

type Register struct {
	Avatar          string `default:"avatar/sys_default/default.jpeg"`
	Username        string `v:"required|phone|check-qc-username-exist#手机号必填|手机号格式异常|该手机号已经被注册"`
	Code            string `v:"required|check-phone-code-correct#验证码必填|验证码错误"`
	Password        string `v:"required#用户密码必填"`
	IdNo            string `v:"required|resident-id#身份证号码必填|身份证号码异常"`
	RealName        string `v:"required#用户姓名必填"`
	RegionId        string `v:"required#地区必填"`
	HospitalId      int    `v:"required#医院必填"`
	PositionalTitle string
}

type GetSms struct {
	Type     string `v:"required|in:register,login,find_pwd,change_username,modify_pwd#参数错误:使用场景必填|参数错误：使用场景取值异常"`
	Username string `v:"required|check-qc-phone-exist#手机号必填|信息异常"`
}

type FindPwdReq struct {
	Username   string `v:"required|check-qc-username-not-exist#用户名必填|用户不存在"`
	Code       string `v:"required#短信验证码必填"`
	Password   string `v:"required#用户密码必填"`
	RePassword string `v:"required|same:password#用户密码必填|两次密码输入不一致"`
}

// 修改用户基本信息, 如果修改了手机号，那么验证码必填
type ModifyReq struct {
	Avatar          string `default:"avatar/qc/default.jpeg"`
	IdNo            string `v:"required|resident-id#身份证号码必填|身份证号码异常"`
	RealName        string `v:"required#用户姓名必填"`
	RegionId        string `v:"required#地区必填"`
	HospitalId      int    `v:"required#医院必填"`
	PositionalTitle string
	Password        string `v:"required#密码必填"`
}

// 修改手机号
type ModifyUsernameReq struct {
	Username string `v:"required|phone#手机号必填|手机号格式异常"`
	Code     string `v:"required#短信验证码必填"`
	Password string `v:"required#密码必填"`
}

type UserBindReq struct {
	SerialNumber string `v:"required#设备系列号必填"`
}

type PageParams struct {
	TimeBegin string
	TimeEnd   string
	Keyword   string `p:"keyword"`
	Status    int    `p:"status"`
	Page      int    `p:"page"`
	PageSize  int    `p:"page_size"`
	Order     string `p:"order" default:"id"`
	Sort      string `p:"sort" default:"DESC"`
}
