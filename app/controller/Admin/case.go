package Admin

import (
	"aiyun_local_srv/app/model/case_info"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Case = caseApi{}

type caseApi struct{}

func (*caseApi) List(r *ghttp.Request) {
	var req *case_info.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.CaseService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

func (*caseApi) Delete(r *ghttp.Request) {
	var req *case_info.DelReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CaseService.Delete(g.Map{"id": req.Ids}); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
