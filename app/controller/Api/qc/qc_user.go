// +----------------------------------------------------------------------
// 远程质控用户相关接口
// @Author sinbook <778774780@qq.com>
// @Copyright 版权所有 广州爱孕记信息科技有限公司 [ https://www.aiyunji.cn/ ]
// @Date 2022-10-26 10:01:30
// +----------------------------------------------------------------------

package qc

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/model/qc_user"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils"
	"aiyun_local_srv/library/utils/cache"
	"aiyun_local_srv/library/utils/jwt"
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/gogf/gf/util/gvalid"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//远程质控

var QcUser = qcUserApi{}

type qcUserApi struct{}

//登录
func (*qcUserApi) Login(r *ghttp.Request) {
	var req *qc_user.Login

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	user, err := service.QcService.GetUserInfo(g.Map{"username": req.Username})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if user == nil {
		response.Error(r, "用户名或密码错误")
	}

	if service.QcService.QcUserPasswordRule(req.Password) != user.Password {
		response.Error(r, "用户名或密码错误")
	}

	if user.Status != 1 {
		response.Error(r, "用户被禁止，请联系管理员")
	}

	//生成token
	token, err := jwt.GenerateQcUserLoginToken(user, r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("qc.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	//清除旧的token
	old_token, err := cache.Get(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if old_token != "" {
		if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+old_token), "0", time.Second*time.Duration(token_expires)); err != nil {
			response.ErrorSys(r, err)
		}

	}

	//记录用户id和token的映射关系
	if err := cache.Set(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	user.Password = ""
	response.Success(r, g.Map{
		"user":   user,
		"token":  token,
		"expire": time.Now().Unix() + token_expires,
	})
}

// 短信验证码登录
func (*qcUserApi) SmsLogin(r *ghttp.Request) {
	var req *qc_user.SmsLogin

	if err := r.Parse(&req); err != nil {
		if err.(gvalid.Error).FirstString() == "验证码错误" {
			response.Json(r, 2, err.(gvalid.Error).FirstString())
		}
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("login", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if code != req.Code {
		response.Error(r, "短信验证码错误")
	}

	user, err := service.QcService.GetUserInfo(g.Map{"username": req.Username})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if user == nil {
		response.Error(r, "手机号未注册")
	}

	if user.Status != 1 {
		response.Error(r, "用户被禁止，请联系管理员")
	}

	//生成token
	token, err := jwt.GenerateQcUserLoginToken(user, r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("qc.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	//清除旧的token
	old_token, err := cache.Get(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if old_token != "" {
		if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+old_token), "0", time.Second*time.Duration(token_expires)); err != nil {
			response.ErrorSys(r, err)
		}

	}

	//记录用户id和token的映射关系
	if err := cache.Set(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	// 清除code
	err = cache.Del(_const.PHONE_CODE_CACHE_KEY("login", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}

	user.Password = ""
	response.Success(r, g.Map{
		"user":   user,
		"token":  token,
		"expire": time.Now().Unix() + token_expires,
	})
}

//退出登录
func (*qcUserApi) Logout(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if err := cache.Del(_const.TOKEN_CACHE_KEY("qc_user:" + token)); err != nil {
		response.ErrorSys(r, err)
	}
	response.Success(r)
}

//刷新token
func (*qcUserApi) RefreshToken(r *ghttp.Request) {
	//验证登录
	info, err := service.QcService.GetUserInfo(g.Map{"id": r.GetCtxVar("uid").Int()})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "用户不存在")
	}

	if info.Status != 1 {
		response.Error(r, "用户被禁止，请联系管理员")
	}

	//生成token
	token, err := jwt.GenerateQcUserLoginToken(info, r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("qc.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	//清除旧的token
	old_token, err := cache.Get(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", info.Id)))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if old_token != "" {
		if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+old_token), "0", time.Second*time.Duration(token_expires)); err != nil {
			response.ErrorSys(r, err)
		}

	}

	//记录用户id和token的映射关系
	if err := cache.Set(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", info.Id)), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	info.Password = ""
	response.Success(r, g.Map{
		"user":   info,
		"token":  token,
		"expire": time.Now().Unix() + token_expires,
	})
}

//注册
func (*qcUserApi) Register(r *ghttp.Request) {
	var req *qc_user.Register

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//验证地区和医院数据是否合法
	region_name, err := service.RegionService.GetNameById(req.RegionId)
	if err != nil || region_name == "" {
		response.Error(r, "地区数据异常")
	}

	hospital_list, err := service.PublicService.CloudHospitalList(req.RegionId)
	if err != nil {
		response.ErrorSys(r, err)
	}

	hospital_name := ""
	for _, v := range hospital_list {
		if v.Id == req.HospitalId {
			hospital_name = v.Name
		}
	}
	if hospital_name == "" {
		response.Error(r, "医院不存在")
	}

	if err := service.QcService.Register(req, region_name, hospital_name); err != nil {
		response.ErrorDb(r, err)
	}

	user, err := service.QcService.GetUserInfo(g.Map{"username": req.Username})
	if err != nil {
		response.ErrorDb(r, err)
	}
	//生成token
	token, err := jwt.GenerateQcUserLoginToken(user, r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("qc.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	//清除旧的token
	old_token, err := cache.Get(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if old_token != "" {
		if err := cache.Set(_const.TOKEN_CACHE_KEY("qc_user:"+old_token), "0", time.Second*time.Duration(token_expires)); err != nil {
			response.ErrorSys(r, err)
		}

	}

	//记录用户id和token的映射关系
	if err := cache.Set(_const.TOKEN_CACHE_KEY(fmt.Sprintf("qc_user:%d", user.Id)), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	user.Password = ""
	response.Success(r, g.Map{
		"user":   user,
		"token":  token,
		"expire": time.Now().Unix() + token_expires,
	})

	response.Success(r)
}

//获取短信验证码
func (*qcUserApi) GetSms(r *ghttp.Request) {
	var req *qc_user.GetSms

	if err := r.Parse(&req); err != nil {
		if err.(gvalid.Error).FirstString() == "验证码错误" {
			response.Json(r, 2, err.(gvalid.Error).FirstString())
		}
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code := grand.N(100000, 999999)
	code = 123456
	if req.Type == "register" {
		//if err := aliyun_sms.SMS(req.Username, "SMS_251016135", g.MapStrAny{"code": code}); err != nil {
		//	response.Error(r, "短信发送失败")
		//}
	} else if req.Type == "login" {
		//if err := aliyun_sms.SMS(req.Username, "SMS_251016135", g.MapStrAny{"code": code}); err != nil {
		//	response.Error(r, "短信发送失败")
		//}
	} else if req.Type == "find_pwd" {
		//if err := aliyun_sms.SMS(req.Username, "SMS_251016135", g.MapStrAny{"code": code}); err != nil {
		//	response.Error(r, "短信发送失败")
		//}
	} else if req.Type == "change_username" {
		//if err := aliyun_sms.SMS(req.Username, "SMS_251016135", g.MapStrAny{"code": code}); err != nil {
		//	response.Error(r, "短信发送失败")
		//}
	} else if req.Type == "modify_pwd" {
		//if err := aliyun_sms.SMS(req.Username, "SMS_251016135", g.MapStrAny{"code": code}); err != nil {
		//	response.Error(r, "短信发送失败")
		//}
	}

	if err := cache.Set(_const.PHONE_CODE_CACHE_KEY(req.Type, req.Username), fmt.Sprintf("%d", code), time.Minute*5); err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r)
}

// 获取用户信息
func (*qcUserApi) Info(r *ghttp.Request) {
	uid := r.GetCtxVar("uid").Int()

	info, err := service.QcService.GetUserInfo(g.Map{"id": uid})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, info)
}

// 获取用户信息
func (*qcUserApi) List(r *ghttp.Request) {
	var req *qc_user.PageParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.QcService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

// 绑定用户，用于客户端静默登录
func (*qcUserApi) Bind(r *ghttp.Request) {
	var req *qc_user.UserBindReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}
	response.Success(r)
}

// 修改用户基本信息
func (*qcUserApi) Modify(r *ghttp.Request) {
	var req *qc_user.ModifyReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//验证地区和医院数据是否合法
	region_name, err := service.RegionService.GetNameById(req.RegionId)
	if err != nil || region_name == "" {
		response.Error(r, "地区数据异常")
	}

	hospital_list, err := service.PublicService.CloudHospitalList(req.RegionId)
	if err != nil {
		response.ErrorSys(r, err)
	}

	hospital_name := ""
	for _, v := range hospital_list {
		if v.Id == req.HospitalId {
			hospital_name = v.Name
		}
	}
	if hospital_name == "" {
		response.Error(r, "医院不存在")
	}

	uid := r.GetCtxVar("uid").Int()

	// 如果密码不为空验证短信验证码是否正确
	user, err := service.QcService.GetUserInfo(g.Map{"id": uid})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if user == nil {
		response.Error(r, "用户名信息异常")
	}

	if service.QcService.QcUserPasswordRule(req.Password) != user.Password {
		response.Error(r, "密码错误")
	}

	if err := service.QcService.UserModify(req, uid, region_name, hospital_name); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 修改用户手机号码
func (*qcUserApi) ModifyUsername(r *ghttp.Request) {
	var req *qc_user.ModifyUsernameReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	uid := r.GetCtxVar("uid").Int()

	// 如果密码不为空验证短信验证码是否正确
	user, err := service.QcService.GetUserInfo(g.Map{"id": uid})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if user == nil {
		response.Error(r, "用户名信息异常")
	}

	if req.Username != user.Username {
		code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("change_username", req.Username))
		if err != nil {
			response.ErrorSys(r, err)
		}
		if code != req.Code {
			response.Json(r, 4, "短信验证码错误")
		}
	}

	if service.QcService.QcUserPasswordRule(req.Password) != user.Password {
		response.Error(r, "密码错误")
	}

	if err := service.QcService.UserModifyUsername(req.Username, uid); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 用户头像上传
func (*qcUserApi) UploadAvatar(r *ghttp.Request) {
	tempFile := r.GetUploadFile("file")
	if tempFile == nil {
		response.Error(r, "未获取到图片")
	}

	//保存文件
	date := gtime.Date()
	path := fmt.Sprintf("%s/%s/", "avatar/qc", date)
	ext := utils.Ext(tempFile.Filename)
	fileName := tempFile.Filename
	tempFile.Filename = gmd5.MustEncrypt(grand.Letters(10)) + "." + ext

	_, err := tempFile.Save("public/" + path)
	if err != nil {
		response.Error(r, "文件上传失败")
	}

	base_url := g.Cfg().GetString("server.qc_server.Domain")
	response.Success(r, g.Map{
		"base_url": base_url,
		"path":     path + tempFile.Filename,
		"name":     fileName,
		"size":     tempFile.Size,
		"ext":      ext,
	})
}

// 找回密码短信验证
func (*qcUserApi) FindPwdVerify(r *ghttp.Request) {
	var req *qc_user.SmsLogin

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("find_pwd", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if code != req.Code {
		response.Json(r, 4, "短信验证码错误")
	}

	response.Success(r)
}

// 找回密码
func (*qcUserApi) FindPwd(r *ghttp.Request) {
	var req *qc_user.FindPwdReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("find_pwd", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if code != req.Code {
		response.Json(r, 4, "短信验证码错误")
	}

	if err := service.QcService.FindPwd(req); err != nil {
		response.ErrorDb(r, err)
	}

	cache.Del(_const.PHONE_CODE_CACHE_KEY("find_pwd", req.Username))

	response.Success(r)
}

// 系统默认头像
func (*qcUserApi) AvatarList(r *ghttp.Request) {
	base_url := g.Cfg().GetString("server.qc_server.Domain")
	avatar_list := []string{}
	filepath.Walk("public/avatar/sys_default", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			avatar_list = append(avatar_list, strings.TrimLeft(strings.ReplaceAll(path, "\\", "/"), "public/"))
		}

		return nil
	})

	response.Success(r, g.Map{
		"base_url":    base_url,
		"avatar_list": avatar_list,
	})
}

// 权限内修改密码短信验证（步骤一）
func (*qcUserApi) ModifyPwdVerify(r *ghttp.Request) {
	var req *qc_user.SmsLogin

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("modify_pwd", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if code != req.Code {
		response.Json(r, 4, "短信验证码错误")
	}

	response.Success(r)
}

// 权限内修改密码（步骤二）
func (*qcUserApi) ModifyPwd(r *ghttp.Request) {
	var req *qc_user.FindPwdReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	code, err := cache.Get(_const.PHONE_CODE_CACHE_KEY("modify_pwd", req.Username))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if code != req.Code {
		response.Json(r, 4, "短信验证码错误")
	}

	if err := service.QcService.FindPwd(req); err != nil {
		response.ErrorDb(r, err)
	}

	cache.Del(_const.PHONE_CODE_CACHE_KEY("modify_pwd", req.Username))

	response.Success(r)
}
