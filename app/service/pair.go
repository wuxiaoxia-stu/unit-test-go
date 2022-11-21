package service

import (
	"aiyun_local_srv/app/model/pre_pair"
	"database/sql"
	"github.com/gogf/gf/frame/g"
)

var PairService = new(pairService)

type pairService struct{}

// 检查授权记录
func (s *pairService) Info(where g.Map) (res *pre_pair.Entity, err error) {
	err = pre_pair.M.Where(where).Order("id DESC").Limit(1).Scan(&res)
	return
}

//创建授权信息
func (s *pairService) Create(req *pre_pair.Entity) (result sql.Result, err error) {
	result, err = pre_pair.M.Data(g.Map{
		"ukey_code":         req.UkeyCode,
		"author_number":     req.AuthorNumber,
		"device_number":     req.DeviceNumber,
		"uuid":              req.Uuid,
		"public_key":        req.PublicKey,
		"ip":                req.Ip,
		"floor":             req.Floor,
		"room":              req.Room,
		"machine_type":      req.MachineType,
		"client_version":    req.ClientVersion,
		"algorithm_version": req.AlgorithmVersion,
		"status":            1,
	}).Insert()
	return
}

func (s *pairService) Update(data g.Map, where g.Map) (err error) {
	_, err = pre_pair.M.Where(where).Data(data).Update()
	return
}
