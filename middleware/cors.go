package middleware

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gookit/color"
)

// CORS 跨域处理中间件
func CORS(r *ghttp.Request) {
	color.Println("")
	color.Println("")
	color.Println("<fg=CC00FF>路由地址：" + r.GetUrl() + "</>")
	if g.Cfg().GetBool("server.Debug") {
		color.Println("<fg=FF0066>提交数据：</><fg=CCFF33>" + string(r.GetBody()) + "</>")
	}
	r.Response.CORSDefault()
	r.Middleware.Next()
}
