package Api

import (
	"aiyun_local_srv/app/model/case_info"
	"aiyun_local_srv/app/model/case_label_relate"
	"aiyun_local_srv/app/model/case_measured"
	"aiyun_local_srv/app/model/case_plane"
	"aiyun_local_srv/app/model/case_shot"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"strings"
)

//病例相关接口

var Case = caseApi{}

type caseApi struct{}

//创建病例
func (*caseApi) Create(r *ghttp.Request) {
	var req *case_info.CreateCaseReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.CaseType == 3 && req.MultiGroupId == "" {
		response.Error(r, "此病例未多人质控病例，多人质控分组标识必填")
	}

	//检测检查是否存在
	check_info, err := service.CheckService.Info(g.Map{"id": req.CheckId})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if check_info == nil {
		response.Error(r, "产筛不存在")
	}

	if (check_info.Type == 2 && req.CaseType == 1) || (check_info.Type == 1 && req.CaseType != 1) {
		response.Error(r, "产筛数据和病例类型不匹配")
	}

	uid := r.GetCtxVar("uid").Int()
	serial_number := r.GetCtxVar("serial_number").String()
	case_id, err := service.CaseService.Create(req, uid, serial_number, check_info)
	if err != nil {
		response.ErrorDb(r, err)
	}

	info, err := service.CaseService.Info(g.Map{"ci.id": case_id})
	if err != nil {
		return
	}

	response.Success(r, info)
}

// 保存病例评分标准
func (*caseApi) SaveScoreStandard(r *ghttp.Request) {
	var req *case_info.SaveScoreStandardReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	serial_number := r.GetCtxVar("serial_number").String()
	if err := service.CaseService.SaveScoreStandard(req, serial_number); err != nil {
		response.ErrorDb(r, err)
	}

	info, err := service.CaseService.Info(g.Map{"ci.id": req.CaseId})
	if err != nil {
		return
	}

	response.Success(r, info)
}

// 更新病例
func (*caseApi) Update(r *ghttp.Request) {
	var req *case_info.UpdateCaseReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CaseService.Update(req); err != nil {
		response.ErrorDb(r, err)
	}

	info, _ := service.CaseService.Detail(req.CaseId)
	response.Success(r, info)
}

// 删除无效的病例
func (*caseApi) Delete(r *ghttp.Request) {
	case_id := r.GetQueryInt("case_id")
	if err := service.CaseService.Delete(g.Map{"id": case_id}); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//结束病例
func (*caseApi) Finish(r *ghttp.Request) {
	var req *case_info.FinishCaseReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CaseService.Finish(req); err != nil {
		response.ErrorDb(r, err)
	}

	info, _ := service.CaseService.Detail(req.CaseId)
	response.Success(r, info)
}

//病例时间线， 用于病例搜索
func (*caseApi) TimeLine(r *ghttp.Request) {
	var req *case_info.CaseTimelineReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	time_line, case_count, case_curr_count, err := service.CaseService.GetTimeLine(req, r.GetCtxVar("serial_number").String())
	if err != nil {
		response.ErrorDb(r, err)
	}

	var tree []*case_info.Year

	curr_year := ""
	curr_month := ""
	for _, v := range time_line {
		t_len := len(tree)
		d := strings.Split(v.Date, "-")
		if curr_year == d[0] && curr_month == d[1] {
			tree[t_len-1].Children[len(tree[t_len-1].Children)-1].Children =
				append(tree[t_len-1].Children[len(tree[t_len-1].Children)-1].Children, &case_info.Day{
					Day: d[2],
				})
		} else if curr_year == d[0] {
			tree[t_len-1].Children = append(tree[t_len-1].Children, &case_info.Month{
				Month: d[1],
				Children: []*case_info.Day{
					&case_info.Day{
						Day: d[2],
					},
				},
			})
		} else {
			tree = []*case_info.Year{
				&case_info.Year{
					Year: d[0],
					Children: []*case_info.Month{
						&case_info.Month{
							Month: d[1],
							Children: []*case_info.Day{
								&case_info.Day{
									Day: d[2],
								},
							},
						},
					},
				},
			}
		}

		curr_year = d[0]
		curr_month = d[1]
	}

	response.Success(r, g.Map{
		"case_count":      case_count,
		"case_curr_count": case_curr_count,
		"time_line":       tree,
	})
}

//获取病例列表
func (*caseApi) List(r *ghttp.Request) {
	var req *case_info.CaseListReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	serial_number := r.GetCtxVar("serial_number").String()
	list, err := service.CaseService.List(req, serial_number, false)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

//获取续接病例列表
func (*caseApi) ListJoin(r *ghttp.Request) {
	var req *case_info.CaseListJoinReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	serial_number := r.GetCtxVar("serial_number").String()
	list, err := service.CaseService.ListJoin(serial_number, req)
	if err != nil {
		response.ErrorDb(r, err)
	}
	response.Success(r, list)
}

//重置病例
func (*caseApi) Reset(r *ghttp.Request) {
	response.Success(r)
}

//病例详情
func (*caseApi) Detail(r *ghttp.Request) {
	case_id := r.GetQueryInt64("case_id")
	//serial_number := r.GetCtxVar("serial_number").String()
	info, err := service.CaseService.Detail(case_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "病例不存在")
	}

	response.Success(r, info)
}

//病例添加标签
func (*caseApi) AddLabel(r *ghttp.Request) {
	var req *case_label_relate.CaseAddLabelReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	label, err := service.LabelService.Info(g.Map{"id": req.LabelId})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if label == nil {
		response.Error(r, "标签不存在")
	}

	//唯一性校验  标签和病例联合唯一
	var info *case_label_relate.Entity
	if err = case_label_relate.M.Where(g.Map{"case_id": req.CaseId, "label_id": req.LabelId}).Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}

	if info != nil {
		response.Error(r, "此病例已经关联了该标签")
	}

	count, err := case_label_relate.M.Where("case_id", req.CaseId).Count()
	if err != nil {
		response.ErrorDb(r, err)
	}
	label_count := g.Cfg().GetInt("case.LabelCount", 5)
	if count >= label_count {
		response.Error(r, fmt.Sprintf("单个病例至多关联%d个标签", label_count))
	}

	if err := service.CaseService.AddLabel(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, label)
}

//病例删除标签
func (*caseApi) DeleteLabel(r *ghttp.Request) {
	var req *case_label_relate.CaseAddLabelReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	label, err := service.LabelService.Info(g.Map{"id": req.LabelId})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if err := service.CaseService.DelLabel(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, label)
}

//病例设置测量值
func (*caseApi) SaveMeasured(r *ghttp.Request) {
	var req []*case_measured.CaseSetMeasuredReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CaseService.SetMeasured(req, r.GetCtxVar("serial_number").String()); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//病例截图
func (*caseApi) SaveShot(r *ghttp.Request) {
	var req *case_shot.CaseSaveShotReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}
	serial_number := r.GetCtxVar("serial_number").String()

	if req.ShotPts > 0 {
		var shot_info *case_shot.Entity
		err := case_shot.M.Where(g.Map{"serial_number": serial_number, "case_id": req.CaseId, "pts": req.ShotPts}).Scan(&shot_info)
		if err != nil {
			response.ErrorDb(r, err)
		}

		if shot_info != nil {
			response.Error(r, "该截图已存在")
		}
	}

	if err := service.CaseService.SaveShot(req, serial_number); err != nil {
		response.ErrorDb(r, err)
	}

	var info *case_shot.Entity
	case_shot.M.Where(g.Map{"shot_id": req.ShotId}).Scan(&info)
	response.Success(r, info)
}

// 删除截图
func (*caseApi) DeleteShot(r *ghttp.Request) {
	var req *case_shot.CaseDeleteShotReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CaseService.DeleteShot(req.ShotIds); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 保存截图分数
//func (*caseApi) SaveShotScore(r *ghttp.Request) {
//	var req *case_shot_plane.ShotScoreSaveReq
//
//	if err := r.Parse(&req); err != nil {
//		response.Error(r, err.(gvalid.Error).FirstString())
//	}
//
//	if len(req.PlaneList) <= 0 {
//		response.Error(r, "数据异常，未检测到切面数据")
//	}
//
//	serial_number := r.GetCtxVar("serial_number").String()
//	err := service.CaseService.SaveShotScore(req, serial_number)
//	if err != nil {
//		response.ErrorDb(r, err)
//	}
//
//	response.Success(r)
//}

// 图片使用
func (*caseApi) ShotRelatePlane(r *ghttp.Request) {
	var req *case_plane.ShotUsedReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	err := service.CaseService.ShotUsed(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//切面分数详情
func (*caseApi) PlaneScoreDetail(r *ghttp.Request) {
	case_id := r.GetQueryInt64("case_id")

	plane_score_list, err := service.CaseService.PlaneScoreDetail(case_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, plane_score_list)
}

// 获取多人质控 病例列表
func (*caseApi) MultiGroupList(r *ghttp.Request) {
	multi_group_id := r.GetQueryString("multi_group_id")
	if multi_group_id == "" {
		response.Error(r, "病例分组ID必填")
	}

	case_list, max_score, min_score, averag_score, median_score, err := service.CaseService.MultiList(multi_group_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"case_list":    case_list,
		"max_score":    max_score,
		"min_score":    min_score,
		"averag_score": averag_score,
		"median_score": median_score,
	})
}

// 修改分数
func (*caseApi) ModifyScore(r *ghttp.Request) {
	var req *case_plane.ModifyScoreReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.PlaneId == "" && req.GroupId == "" {
		response.Error(r, "切面ID或分组ID必填一项")
	}

	if len(req.ScoreItem) <= 0 {
		response.Error(r, "改分项必填")
	}

	where := g.Map{"case_id": req.CaseId}
	if req.GroupId != "" {
		where = g.Map{"group_id": req.GroupId}
	}
	if req.PlaneId != "" {
		where = g.Map{"plane_id": req.PlaneId}
	}

	// 查询切面
	plane_info, err := service.CaseService.GetPlaneList(where)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if plane_info == nil {
		response.Error(r, "切面不存在")
	}

	structure_ids := []string{}
	for _, v := range req.ScoreItem {
		structure_ids = append(structure_ids, v.StructureId)
	}
	where = g.Map{"structure_id": structure_ids}

	structure_list, err := service.CaseService.GetStructureList(where)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if len(structure_list) != len(req.ScoreItem) {
		response.Error(r, "分数数据异常")
	}

	for _, v := range req.ScoreItem {
		for _, v2 := range structure_list {
			if v.StructureId == v2.StructureId {
				if v2.StructureScoreTotal > 0 {
					if v.Score < 0 || v.Score > v2.StructureScoreTotal {
						response.Error(r, "分数异常")
					}
				} else {
					if v.Score > 0 || v.Score < v2.StructureScoreTotal {
						response.Error(r, "分数异常")
					}
				}
			}
		}
	}

	if err = service.CaseService.ModifyScore(req, plane_info, structure_list, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 隐藏病例
func (*caseApi) Hidden(r *ghttp.Request) {
	var req *case_info.CaseHiddenStatusReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	// 修改病例隐藏状态
	err := service.CaseService.ModifyCaseHidden(req.CaseIds, 1)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 隐藏病例列表
func (*caseApi) HiddenList(r *ghttp.Request) {
	var req *case_info.HiddenListReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	// 获取影藏病例列表
	total, list, err := service.CaseService.CaseHiddenList(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

// 显示病例列表
func (*caseApi) Display(r *ghttp.Request) {
	var req *case_info.CaseHiddenStatusReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	// 修改病例隐藏状态
	err := service.CaseService.ModifyCaseHidden(req.CaseIds, 0)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 永久隐藏
func (*caseApi) Sealing(r *ghttp.Request) {
	var req *case_info.CaseHiddenStatusReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//检查病例隐藏状态是否正常

	// 修改病例隐藏状态
	err := service.CaseService.ModifyCaseHidden(req.CaseIds, 2)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
