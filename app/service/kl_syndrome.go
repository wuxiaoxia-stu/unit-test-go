package service

import (
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_syndrome"
	"aiyun_local_srv/app/model/kl_syndrome_feature"
	"encoding/json"
)

//综合征
var KlSyndromeService = new(klSyndromeService)

type klSyndromeService struct{}

//获取列表数据
func (s *klSyndromeService) Info(syndrome_id int) (res *kl_syndrome.Entity, err error) {
	err = kl_syndrome.M.Where("id", syndrome_id).Limit(1).Scan(&res)
	return
}

//获取列表数据
func (s *klSyndromeService) List(where interface{}, order string) (res []*kl_syndrome.Entity, err error) {
	err = kl_syndrome.M.Where(where).Order(order).Scan(&res)
	return
}

//获取当前节点及其所有子节点数据
func (s *klSyndromeService) Tree(where interface{}) (tree []*kl_syndrome.SyndromeListRsp, err error) {
	var list []*kl_syndrome.SyndromeListRsp

	err = kl_syndrome.M.Fields("id,type,sub_type,id as syndrome_uuid,name as syndrome_name,name_en as syndrome_name_en").Order("sort,id").Where(where).Scan(&list)
	if err != nil {
		return
	}

	data, _ := json.Marshal(kl_syndrome.TypeTree)
	if err = json.Unmarshal(data, &tree); err != nil {
		return
	}

	for _, v := range tree {
		for _, v2 := range list {
			if v2.SubType == 0 && v.Type == v2.Type {
				v.Children = append(v.Children, v2)
			}
		}
	}

	for _, v := range tree {
		if v.Type == 1 {
			v.SyndromeSequence = 1
		}
		for _, v2 := range v.Children {
			for _, v3 := range list {
				if v3.SubType > 0 && v3.SubType == v2.SubType {
					v2.Children = append(v2.Children, v3)
				}
			}
		}
	}

	return
}

//获取id路径地址
func (s *klSyndromeService) Detail(id int) (res *kl_syndrome.Entity, err error) {
	if err = kl_syndrome.M.Where("id", id).Scan(&res); err != nil {
		return
	}

	var syndrome_feature_list []*kl_feature.Entity
	if err = kl_syndrome_feature.M_alias.
		Fields("kf.id,kf.name,kf.name_en,ksf.type,ksf.feature_root_id as root_id").
		LeftJoin("kl_feature kf", "kf.id = ksf.feature_id").
		Where("ksf.syndrome_id", id).
		Scan(&syndrome_feature_list); err != nil {
		return
	}

	if res != nil {
		res.Features = syndrome_feature_list
	}
	return
}
