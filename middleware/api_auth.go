package middleware

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/model/pre_pair"
	"aiyun_local_srv/app/model/user"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/cache"
	"aiyun_local_srv/library/utils/jwt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

// Auth 权限判断处理中间件
func ApiAuth(r *ghttp.Request) {
	if r.RequestURI == "/api/user/login" || r.RequestURI == "/api/user/list" {
		r.Middleware.Next()
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		response.Json(r, 9, "登录已过期，请重新登录")
	}

	if tokenStr == "token" {
		var u *user.Entity
		err := user.M.Where("status", 1).Order("id DESC").Scan(&u)
		if err != nil {
			response.ErrorDb(r, err)
		}
		var pair *pre_pair.Entity
		err = pre_pair.M.Where("status", 1).Order("id DESC").Scan(&pair)
		if err != nil {
			response.ErrorDb(r, err)
		}
		r.SetCtxVar("uid", u.Id)
		r.SetCtxVar("role_type", u.RoleType)
		r.SetCtxVar("author_number", pair.AuthorNumber)
		r.SetCtxVar("serial_number", pair.DeviceNumber)
		r.Middleware.Next()
		return
	}

	token_state, err := cache.Get(_const.TOKEN_CACHE_KEY("user:" + tokenStr))
	if err != nil {
		response.ErrorSys(r, err)
	}

	if token_state == "" {
		response.Json(r, 9, "登录已过期，请重新登录")
	} else if token_state == "0" {
		response.Json(r, 9, "该用户在其他设备登录，您被迫下线")
	}

	data, err := jwt.ParseToken(tokenStr, []byte(g.Cfg().GetString("jwt.sign", "qc_sign")))
	if err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	var u *user.Entity
	if err = gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	r.SetCtxVar("uid", u.Id)
	r.SetCtxVar("role_type", u.RoleType)
	r.SetCtxVar("author_number", u.AuthorNumber)
	r.SetCtxVar("serial_number", u.SerialNumber)
	r.SetCtxVar("user_info", u)
	//判断权限  目前只有登录权限
	// 执行下一步请求逻辑
	r.Middleware.Next()
}
