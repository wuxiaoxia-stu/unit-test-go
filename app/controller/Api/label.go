package Api

import (
	"aiyun_local_srv/app/model/case_label"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

//标签相关接口
var Label = labelApi{}

type labelApi struct{}

//创建标签
func (*labelApi) Create(r *ghttp.Request) {
	var req *case_label.CreateLabelReq

	if err := r.Parse(&req); err != nil {
		if err.(gvalid.Error).FirstString() == "标签已存在" {
			var res *case_label.Entity
			if err := case_label.M.Where(g.Map{"name": req.LabelName}).Scan(&res); err != nil {
				response.ErrorDb(r, err)
			}
			response.Success(r, res)
		}
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	info, err := service.LabelService.Create(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, info)
}

//删除标签
func (*labelApi) Delete(r *ghttp.Request) {
	id := r.GetQueryInt("label_id")

	lable, err := service.LabelService.Info(g.Map{"id": id})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if lable == nil {
		response.Error(r, "标签不存在")
	}

	if err := service.LabelService.Delete(g.Map{"id": id}); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, lable)
}

//我的标签
func (*labelApi) UserList(r *ghttp.Request) {
	limit := r.GetQueryInt("limit", 5)
	case_id := r.GetQueryInt64("case_id")

	used_label, err := service.LabelService.RecentUesdList(case_id, r.GetCtxVar("uid").Int(), limit)
	if err != nil {
		response.ErrorDb(r, err)
	}

	create_label, err := service.LabelService.UserList(case_id, r.GetCtxVar("uid").Int(), limit, g.Map{"type": 2})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"used_label":   used_label,
		"create_label": create_label,
	})
}

//系统标签
func (*labelApi) SystemList(r *ghttp.Request) {
	case_id := r.GetQueryInt64("case_id")

	label_list, err := service.LabelService.SystemList(case_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, label_list)
}

//标签搜索
func (*labelApi) Search(r *ghttp.Request) {
	keywords := r.GetQueryString("key_words")
	uid := r.GetCtxVar("uid").Int()

	label_list, err := service.LabelService.Search(keywords, uid)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, label_list)
}
