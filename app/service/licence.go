package service

import (
	"aiyun_local_srv/app/model"
	"aiyun_local_srv/app/model/licence"
	"database/sql"
	"github.com/gogf/gf/frame/g"
)

var LicenceService = new(licenceService)

type licenceService struct{}

// 检查授权记录
func (s *licenceService) Info(where g.Map) (res *licence.Entity, err error) {
	err = licence.M.Where(where).Order("id DESC").Limit(1).Scan(&res)
	return
}

//创建授权信息
func (s *licenceService) Create(req *licence.LicenceCreateRsp) (result sql.Result, err error) {
	result, err = licence.M.Data(g.Map{
		"author_number": req.AuthorNumber,
		"licence":       req.Licence,
		"licence_key":   req.LicenceKey,
		"ukey_code":     req.UkeyCode,
		"ukey_crypt":    req.UkeyCrypt,
	}).Insert()
	return
}

//设置数据状态
func (s *licenceService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = licence.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *licenceService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = licence.M.WhereIn("id", req.Ids).Delete()
	return
}
