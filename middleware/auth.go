package middleware

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/model/sys_admin"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/cache"
	"aiyun_local_srv/library/utils/jwt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

// Auth 权限判断处理中间件
func Auth(r *ghttp.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		response.Json(r, 9, "Token错误")
	}
	token, err := cache.Get(_const.TOKEN_CACHE_KEY(tokenStr))
	if err != nil || token == "" {
		response.Json(r, 9, "Token无效或已过期")
	}

	data, err := jwt.ParseToken(tokenStr, []byte(g.Cfg().GetString("jwt.sign", "qc_sign")))
	if err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	var u *sys_admin.Entity
	if err = gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	r.SetCtxVar("uid", u.Id)

	//判断权限  目前只有登录权限
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

//获取角色菜单数据
func GetRoleMenu() {

}
