package service

import (
	"aiyun_local_srv/app/model"
	"aiyun_local_srv/app/model/sys_admin"
	"aiyun_local_srv/library/utils"
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
)

var SysAdminService = new(sysAdminService)

type sysAdminService struct{}

// 登录验证
func (s *sysAdminService) GetByUsername(username string) (res *sys_admin.Entity, err error) {
	err = sys_admin.M.Where(g.Map{"username": username}).Scan(&res)
	return
}

// 查询单个信息
func (s *sysAdminService) FindById(id int) (res *sys_admin.Entity, err error) {
	err = sys_admin.M.Where("id=?", id).Scan(&res)
	return
}

//注册新用户
func (s *sysAdminService) Add(req *sys_admin.AddReq) (result sql.Result, err error) {
	salt := grand.Letters(4)
	result, err = sys_admin.M.Data(g.Map{
		"username": req.Username,
		"password": utils.GenPwd(req.Password, salt),
		"salt":     salt,
		//"department_id": strings.Join(req.DepartmentId, ","),
		"email":  req.Email,
		"phone":  req.Phone,
		"status": req.Status,
	}).Insert()
	return
}

//更新用户信息
func (s *sysAdminService) Edit(req *sys_admin.EditReq) (result sql.Result, err error) {
	data := g.Map{
		//"role_id": req.RoleId,
		//"department_id": strings.Join(req.DepartmentId, ","),
		"email":  req.Email,
		"phone":  req.Phone,
		"status": req.Status,
	}

	if req.Password != "" {
		salt := grand.Letters(4)
		data["salt"] = salt
		data["password"] = utils.GenPwd(req.Password, salt)
	}

	result, err = sys_admin.M.Where("id", req.Id).Data(data).Update()
	return
}

//设置数据状态
func (s *sysAdminService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = sys_admin.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *sysAdminService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = sys_admin.M.WhereIn("id", req.Ids).Delete()
	return
}

//设置用户角色
func (s *sysAdminService) SetRole(data *sys_admin.SetRoleReq) (result sql.Result, err error) {
	result, err = sys_admin.M.Where("id", data.Id).Data(g.Map{
		"role_id": data.RoleId,
	}).Update()
	return
}

//获取分页列表
func (s *sysAdminService) Page(req *model.PageReqParams) (total int, list []*sys_admin.Entity, err error) {
	M := sys_admin.M

	if req.KeyWord != "" {
		M = M.WhereLike("username", "%"+req.KeyWord+"%")
	}

	if req.Status != -1 {
		M = M.Where("status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order(req.Order + " " + req.Sort)
	}

	total, err = M.Group("id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("id,role_id,username,department_id,email,phone,avatar,create_at,update_at,status").All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*sys_admin.Entity, len(data))
	err = data.Structs(&list)
	return
}

//获取分页列表
func (s *sysAdminService) GetDepartmentTree() (tree []*sys_admin.DepartmentTree, err error) {
	tree, err = sys_admin.GetDepartmentTree()
	return
}

func (s *sysAdminService) CheckDepartmentId(departmentId []string) (b bool, err error) {
	if len(departmentId) != 2 {
		return
	}

	tree, err := s.GetDepartmentTree()
	if err != nil {
		return
	}

	for _, v := range tree {
		for _, v2 := range v.Children {
			if v.Value == departmentId[0] && v2.Value == departmentId[1] {
				b = true
				return
			}
		}
	}
	return
}
