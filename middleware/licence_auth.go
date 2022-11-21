package middleware

import (
	"aiyun_local_srv/app/model/licence"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/net/ghttp"
)

// 服务授权鉴定处理中间件
func LicenceAuth(r *ghttp.Request) {
	var info *licence.Entity
	if err := licence.M.Where("status", 1).OrderDesc("id").Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Json(r, 7, "服务未授权")
	}

	r.SetCtxVar("srv_author_number", info.AuthorNumber)
	r.SetCtxVar("srv_ukey_code", info.UkeyCode)
	r.Middleware.Next()
}
