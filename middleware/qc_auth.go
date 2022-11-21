package middleware

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/model/qc_user"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils"
	"aiyun_local_srv/library/utils/cache"
	"aiyun_local_srv/library/utils/jwt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

//远程质控白名单接口
var qc_white_list = []string{
	"/api/qc/user/login",
	"/api/qc/user/sms-login",
	"/api/qc/user/register",
	"/api/qc/user/get-sms",
	"/api/qc/user/find-pwd-verify",
	"/api/qc/user/find-pwd",
	"/api/qc/device-reg",
}

// 远程质控权限控制
func QcApiAuth(r *ghttp.Request) {
	if utils.StrInArray(r.RequestURI, qc_white_list) {
		r.Middleware.Next()
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		response.Json(r, 9, "登录已过期，请重新登录")
	}

	token_state, err := cache.Get(_const.TOKEN_CACHE_KEY("qc_user:" + tokenStr))
	if err != nil {
		response.ErrorSys(r, err)
	}

	if token_state == "" {
		response.Json(r, 9, "登录已过期，请重新登录")
	} else if token_state == "0" {
		response.Json(r, 9, "该用户在其他设备登录，您被迫下线")
	}

	data, err := jwt.ParseToken(tokenStr, []byte(g.Cfg().GetString("qc.sign", "qc_sign")))
	if err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	var u *qc_user.Entity
	if err = gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	r.SetCtxVar("uid", u.Id)
	r.SetCtxVar("role_type", u.RoleType)
	r.SetCtxVar("user_info", u)
	//判断权限  目前只有登录权限
	// 执行下一步请求逻辑

	r.Middleware.Next()
}
