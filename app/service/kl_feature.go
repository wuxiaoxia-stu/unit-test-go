package service

import (
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_feature_atlas"
	"strconv"
)

var KlFeatureService = new(klFeatureService)

type klFeatureService struct{}

//获取列表数据
func (s *klFeatureService) Info(where interface{}) (res *kl_feature.Entity, err error) {
	err = kl_feature.M.Where(where).Scan(&res)
	return
}

//获取列表数据
func (s *klFeatureService) GetAtlas(where interface{}) (res []*kl_feature_atlas.Entity, err error) {
	err = kl_feature_atlas.M.Where(where).Scan(&res)
	return
}

//获取列表数据
func (s *klFeatureService) List(where interface{}, order string) (res []*kl_feature.Entity, err error) {
	err = kl_feature.M.Where(where).Order(order).Scan(&res)
	return
}

//获取当前节点及其所有子节点数据
func (s *klFeatureService) Tree(where interface{}) (tree []*kl_feature.FeatureListRsp, err error) {
	var list []*kl_feature.FeatureListRsp

	err = kl_feature.M.Fields("id,pid,id as uuid,id as serial,name,name_en,level as content_type").Order("level,sort,id").Where(where).Scan(&list)
	if err != nil {
		return
	}

	tree = s.ListToTree(list, 0, tree)
	return
}

//获取当前节点及其所有子节点数据
func (s *klFeatureService) ListToTree(source []*kl_feature.FeatureListRsp, id int, result []*kl_feature.FeatureListRsp) []*kl_feature.FeatureListRsp {
	if source == nil {
		return result
	}
	if len(source) == 0 {
		return result
	}

	var otherData []*kl_feature.FeatureListRsp //多余的数据，用于下一次递归
	for _, v := range source {
		if v.Pid == id {
			result = append(result, v)
		} else {
			otherData = append(otherData, v)
		}

	}

	if result != nil {
		for i := 0; i < len(result); i++ {
			result[i].Children = s.ListToTree(otherData, result[i].Id, result[i].Children)
		}
	}

	return result
}

//获取id路径地址
func (s *klFeatureService) FullPath(id int, path string) (full_path string, err error) {
	full_path = path
	var info *kl_feature.Entity
	err = kl_feature.M.Where("id", id).Scan(&info)
	if err != nil {
		return path, err
	}

	if info != nil {
		if info.Pid > 0 {
			if full_path != "" {
				full_path = strconv.Itoa(info.Id) + "," + full_path
			} else {
				full_path = strconv.Itoa(info.Id)
			}
			return s.FullPath(info.Pid, full_path)
		} else {
			if full_path != "" {
				full_path = strconv.Itoa(info.Id) + "," + full_path
			} else {
				full_path = strconv.Itoa(info.Id)
			}
		}
	}

	return
}

//获取id路径地址
func (s *klFeatureService) Detail(id int) (res *kl_feature.Entity, err error) {
	if err = kl_feature.M.Where("id", id).Scan(&res); err != nil {
		return
	}

	if res != nil {
		var feature_atlas []*kl_feature_atlas.Entity
		if err = kl_feature_atlas.M.Where("feature_id", id).Scan(&feature_atlas); err != nil {
			return
		}
		res.Atlas = feature_atlas
	}

	return
}

//获取部位列表
func (s *klFeatureService) PartList() (res []*kl_feature.Entity, err error) {
	err = kl_feature.M.Where("level", 1).Order("sort,id").Scan(&res)
	return
}
