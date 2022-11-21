package Admin

import (
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/net/ghttp"
)

var Kl = klApi{}

type klApi struct{}

// 版本数据导入
func (*klApi) Import(r *ghttp.Request) {
	data_path := r.GetQueryString("data_path")
	if err := service.KlService.Import("public/" + data_path); err != nil {
		response.ErrorSys(r, err)
	}
	response.Success(r)
}
