package router

import (
	"aiyun_local_srv/app/controller/Admin"
	"aiyun_local_srv/app/controller/Api"
	"aiyun_local_srv/app/controller/Api/qc"
	"aiyun_local_srv/middleware"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()

	s.Group("/api/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.Middleware(middleware.LicenceAuth)

		group.ALL("hello", Api.Hello)
		group.ALL("public", Api.Public)

		//知识图谱
		group.ALL("kl", Api.Kl)

		//配对
		group.ALL("pair", Api.Pair)

		//主任秘钥
		group.ALL("leader-key", Api.LeaderKey)

		//鉴权中间件
		group.Middleware(middleware.ApiAuth)

		//病例
		group.ALL("case", Api.Case)

		//检查项
		group.ALL("check", Api.Check)

		//标签
		group.ALL("label", Api.Label)

		//用户
		group.ALL("user", Api.User)
	})

	// 远程质控
	s.Group("/api/qc", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.Middleware(middleware.LicenceAuth)
		group.Middleware(middleware.QcApiAuth)

		////远程质控
		group.ALL("/user", qc.QcUser)
		group.ALL("/", qc.Qc)
		group.ALL("/", qc.Meeting)
	})

	s.Group("/admin/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.ALL("public", Admin.Public)
		group.ALL("authorize", Admin.Authorize)

		//鉴权中间件
		group.Middleware(middleware.LicenceAuth)
		group.Middleware(middleware.Auth)

		group.ALL("case", Admin.Case)
		group.ALL("file", Admin.File)
		group.ALL("kl", Admin.Kl)
		group.ALL("sys-admin", Admin.SysAdmin)
		group.ALL("test", Admin.Test)
	})

	//远程质控服务
	//qc_s := g.Server("qc_server")
	//
	//qc_s.Group("/api/qc", func(group *ghttp.RouterGroup) {
	//	group.Middleware(middleware.CORS)
	//	group.Middleware(middleware.LicenceAuth)
	//	group.Middleware(middleware.QcApiAuth)
	//	//远程质控
	//	group.ALL("/user", qc.QcUser)
	//	group.ALL("/", qc.Qc)
	//})

	//qc_s.Start()
}
