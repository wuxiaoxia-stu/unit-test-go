package service

import (
	"aiyun_local_srv/app/model/check"
	"aiyun_local_srv/library/utils"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"os"
	"strings"
)

var CheckService = new(checkService)

type checkService struct{}

//保存配置
func (s *checkService) Info(where interface{}) (res *check.Entity, err error) {
	err = check.M.Where(where).Scan(&res)
	return
}

//保存配置
func (s *checkService) Save(req []*check.SaveCheckOptions) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	check_ids := []string{}
	for _, v := range req {
		check_ids = append(check_ids, v.CheckId)
	}

	var check_list []*check.Entity
	err = tx.Model(check.Table).Where(g.Map{"id": check_ids}).Scan(&check_list)
	if err != nil {
		tx.Rollback()
		return
	}

	//如果修改了检查名称或者检查切面项，那么删除旧的检查，新建新的检查，如果只是修改状态，那么直接修改检查状态即可
	for _, v := range req {
		for _, v2 := range check_list {
			if v.CheckId == v2.Id {
				if v.CheckName != v2.Name || strings.Join(v.CurrCheckItem, ",") != v2.CheckItem {
					_, err = tx.Model(check.Table).Where(g.Map{"id": v.CheckId}).Delete()
					if err != nil {
						tx.Rollback()
						return
					}

					_, err = tx.Model(check.Table).
						Data(g.Map{
							"id":         utils.GenUUID(),
							"name":       v.CheckName,
							"type":       v.Type,
							"check_type": v.CheckType,
							"level":      v.CheckLevel,
							"check_item": strings.Join(v.CurrCheckItem, ","),
							"sort":       v.Sort,
							"status":     v.Status,
						}).Insert()
					if err != nil {
						tx.Rollback()
						return
					}
				} else if v.Status != v2.Status {
					_, err = tx.Model(check.Table).
						Where(g.Map{"id": v.CheckId}).
						Data(g.Map{
							"status": v.Status,
						}).Update()
					if err != nil {
						tx.Rollback()
						return
					}
				}
			}
		}
	}

	return tx.Commit()
}

//先按照孕期排序，在在按照检查级别排序
func (s *checkService) List(where interface{}) (list []*check.Entity, err error) {
	err = check.M.Where(where).Order("sort,check_type,level DESC").Scan(&list)
	return
}

//读取全部检查项配置文件
func (s *checkService) InitPartAllData(t int) (plane_list map[string][]*check.PartsItem, serr error) {
	file, err := os.Open("./data/plane.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	var parts_all map[string]map[string][]*check.PartsItem
	err = json.Unmarshal(buffer, &parts_all)
	if err != nil {
		return
	}

	if t == 1 {
		plane_list = parts_all["plane"]
	} else {
		plane_list = parts_all["qc_plane"]
	}

	return
}

//读取默认检查项配置文件
func (s *checkService) InitPartDefaultData() (parts_defalut_check []*check.PartsDefaultCheck, serr error) {
	file, err := os.Open("./data/parts_default_check.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	err = json.Unmarshal(buffer, &parts_defalut_check)
	if err != nil {
		return
	}

	err = genCheckList(parts_defalut_check)

	return
}

//生成检查
func genCheckList(parts_defalut_check []*check.PartsDefaultCheck) (err error) {
	//判断检查表里数据是否存在，如果不存在 则通过默认配置数据生成检查数据
	var curr_check_list []*check.Entity
	err = check.M.Scan(&curr_check_list)
	if err != nil {
		return
	}

	list := g.List{}
	for _, v := range parts_defalut_check {
		exist_check := false
		for _, v2 := range curr_check_list {
			if v.Type == v2.Type && v.CheckType == v2.CheckType && v.Level == v2.Level {
				exist_check = true
			}
		}

		if !exist_check {
			list = append(list, g.Map{
				"id":         utils.GenUUID(),
				"name":       v.Name,
				"type":       v.Type,
				"check_type": v.CheckType,
				"level":      v.Level,
				"check_item": strings.Join(v.Default, ","),
				"sort":       v.Sort,
				"status":     1,
			})
		}
	}

	if len(list) > 0 {
		_, err = check.M.Data(list).Insert()
	}

	check_ids := g.Array{}
	for _, v := range curr_check_list {
		exist_check := false
		for _, v2 := range parts_defalut_check {
			if v.Type == v2.Type && v.CheckType == v2.CheckType && v.Level == v2.Level {
				exist_check = true
			}
		}

		if !exist_check {
			check_ids = append(check_ids, v.Id)
		}
	}

	if len(check_ids) > 0 {
		_, err = check.M.Where("id", check_ids).Delete()
	}

	return
}
