package Admin

import (
	"aiyun_local_srv/app/model"
	"aiyun_local_srv/app/model/sys_admin"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var SysAdmin = sysAdminApi{}

type sysAdminApi struct{}

// 获取系统用户列表
func (*sysAdminApi) List(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.SysAdminService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加管理员
func (*sysAdminApi) Add(r *ghttp.Request) {
	var req *sys_admin.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//bool, err := service.SysAdminService.CheckDepartmentId(req.DepartmentId)
	//if err != nil {
	//	response.ErrorDb(r, err)
	//}
	//if !bool {
	//	response.Error(r, "参数错误，部门信息异常")
	//}

	_, err := service.SysAdminService.Add(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//修改管理员信息
func (*sysAdminApi) Edit(r *ghttp.Request) {
	var req *sys_admin.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//bool, err := service.SysAdminService.CheckDepartmentId(req.DepartmentId)
	//if err != nil {
	//	response.ErrorDb(r, err)
	//}
	//if !bool {
	//	response.Error(r, "参数错误，部门信息异常")
	//}

	_, err := service.SysAdminService.Edit(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "编辑成功")
}

//设置管理员状态
func (*sysAdminApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.Id == 1 {
		response.Error(r, "禁止状态操作")
	}

	_, err := service.SysAdminService.SetStatus(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除管理员
func (*sysAdminApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.SysAdminService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//获取部门结构数据
func (*sysAdminApi) GetDepartmentTree(r *ghttp.Request) {
	//tree, err := service.SysAdminService.GetDepartmentTree()
	//if err != nil {
	//	response.ErrorDb(r, err)
	//}

	response.Success(r)
}

//获取路由列表
func (*sysAdminApi) GetAsyncRoutes(r *ghttp.Request) {
	response.Success(r, g.Array{})
}
