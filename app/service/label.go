package service

import (
	"aiyun_local_srv/app/model/case_label"
	"aiyun_local_srv/app/model/case_label_relate"
	"aiyun_local_srv/library/utils"
	"github.com/gogf/gf/frame/g"
)

var LabelService = new(labelService)

type labelService struct{}

func (s *labelService) Info(where interface{}) (res *case_label.Entity, err error) {
	err = case_label.M.Where(where).Limit(1).Scan(&res)
	return
}

// 获取列表
func (s *labelService) RelateList(where interface{}) (list []*case_label_relate.Entity, err error) {
	err = case_label_relate.M.Where(where).Scan(&list)
	return
}

func (s *labelService) Create(req *case_label.CreateLabelReq, uid int) (res *case_label.Entity, err error) {
	_, err = case_label.M.Data(g.Map{
		"user_id": uid,
		"type":    2,
		"name":    req.LabelName,
		"status":  1,
	}).Insert()
	if err != nil {
		return
	}

	err = case_label.M.Where(g.Map{"name": req.LabelName, "type": 2}).Scan(&res)
	return
}

func (s *labelService) Delete(where interface{}) (err error) {
	_, err = case_label.M.Where(where).Data(g.Map{"status": 0}).Update()
	return
}

func (s *labelService) RecentUesdList(case_id int64, uid, limit int) (res []*case_label.Entity, err error) {
	var case_label_relate_list []*case_label_relate.Entity
	err = case_label_relate.M.Where("user_id", uid).OrderDesc("id").Scan(&case_label_relate_list)
	if err != nil {
		return
	}

	label_id_arr := []int{}
	label_relate_arr := []*case_label_relate.Entity{}
	for _, v := range case_label_relate_list {
		if len(label_id_arr) <= limit && !utils.InArray(v.LabelId, label_id_arr) {
			label_id_arr = append(label_id_arr, v.LabelId)
			label_relate_arr = append(label_relate_arr, &case_label_relate.Entity{
				CaseId:  v.CaseId,
				LabelId: v.LabelId,
			})
		}

	}

	if len(label_id_arr) > 0 {
		var label_list []*case_label.Entity
		err = case_label.M.WhereIn("id", label_id_arr).Where("status", 1).Scan(&label_list)
		if err != nil {
			return
		}

		for _, v := range label_relate_arr {
			for _, v2 := range label_list {
				if v2.Id == v.LabelId {
					var case_id_ int64 = 0
					if v.CaseId == case_id {
						case_id_ = case_id
					}
					res = append(res, &case_label.Entity{
						Id:     v.LabelId,
						UserId: v2.UserId,
						Type:   v2.Type,
						Name:   v2.Name,
						CaseId: case_id_,
					})
				}
			}
		}
	}

	return
}

func (s *labelService) UserList(case_id int64, uid, limit int, where interface{}) (res []*case_label.Entity, err error) {
	err = case_label.M.Where(g.Map{"user_id": uid, "status": 1}).Where(where).OrderDesc("id").Limit(limit).Scan(&res)
	if err != nil {
		return
	}

	var relate_list []*case_label_relate.Entity
	err = case_label_relate.M.Where(g.Map{"case_id": case_id}).Scan(&relate_list)
	if err != nil {
		return
	}

	for _, v := range res {
		for _, v2 := range relate_list {
			if v2.LabelId == v.Id {
				v.CaseId = v2.CaseId
			}
		}
	}

	return
}

func (s *labelService) SystemList(case_id int64) (list []*case_label.GroupLabel, err error) {
	var res []*case_label.Entity
	err = case_label.M.Where(g.Map{"type": 1, "status": 1}).Order("sort,id").Scan(&res)
	if err != nil {
		return
	}

	var relate_list []*case_label_relate.Entity
	err = case_label_relate.M.Where(g.Map{"case_id": case_id}).Scan(&relate_list)
	if err != nil {
		return
	}

	for _, v := range res {
		exist := false
		for _, v2 := range list {
			if v2.GroupName == v.GroupName {
				exist = true
				var case_id_n int64 = 0
				for _, v3 := range relate_list {
					if v3.LabelId == v.Id {
						case_id_n = v3.CaseId
					}
				}
				v2.LabelList = append(v2.LabelList, &case_label.LabelList{
					Id:     v.Id,
					Name:   v.Name,
					Type:   v.Type,
					CaseId: case_id_n,
				})
			}
		}

		if !exist {
			label_list := []*case_label.LabelList{}
			var case_id_n int64 = 0
			for _, v3 := range relate_list {
				if v3.LabelId == v.Id {
					case_id_n = v3.CaseId
				}
			}
			list = append(list, &case_label.GroupLabel{
				GroupName: v.GroupName,
				LabelList: append(label_list, &case_label.LabelList{
					Id:     v.Id,
					Name:   v.Name,
					Type:   v.Type,
					CaseId: case_id_n,
				}),
			})
		}
	}
	return
}

//搜索标签
func (s *labelService) Search(keywords string, uid int) (list []*case_label.Entity, err error) {
	err = case_label.M.
		Where(g.Map{"name like": "%" + keywords + "%"}).
		//Where(g.Map{"type": 1, "name like": "%"+keywords+"%"}).
		//WhereOr(g.Map{"type": 2, "user_id": uid, "name like": "%"+keywords+"%"}).
		Order("sort,id").
		Scan(&list)
	return
}
