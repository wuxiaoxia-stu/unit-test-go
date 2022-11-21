package service

import (
	"aiyun_local_srv/app/model/case_info"
	"aiyun_local_srv/app/model/case_label_relate"
	"aiyun_local_srv/app/model/case_measured"
	"aiyun_local_srv/app/model/case_patient"
	"aiyun_local_srv/app/model/case_plane"
	"aiyun_local_srv/app/model/case_plane_structure"
	qc_score_log "aiyun_local_srv/app/model/case_score_log"
	"aiyun_local_srv/app/model/case_shot"
	"aiyun_local_srv/app/model/case_shot_plane"
	"aiyun_local_srv/app/model/case_shot_plane_structure"
	"aiyun_local_srv/app/model/check"
	"aiyun_local_srv/app/model/licence"
	"aiyun_local_srv/app/model/user"
	"aiyun_local_srv/library/utils"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"strconv"
	"strings"
)

var CaseService = new(caseService)

type caseService struct{}

func (s *caseService) Info(where interface{}) (res *case_info.Entity, err error) {
	err = case_info.M_alias.
		Fields("ci.*,c.name as check_name").
		LeftJoin("check c", "c.id = ci.check_id").
		Where(where).
		Scan(&res)

	if res != nil {
		if res.StopAt.String() != "" {
			res.StopTime = res.StopAt.Timestamp()
		}
		res.CreateTime = res.CreateAt.Timestamp()
		res.UpdateTime = res.UpdateAt.Timestamp()
	}
	return
}

func (s *caseService) Create(req *case_info.CreateCaseReq, uid int, serial_number string, check_info *check.Entity) (case_id int64, err error) {
	var licence_info *licence.Entity
	err = licence.M.Where("status", 1).Scan(&licence_info)
	if err != nil {
		return
	}
	if licence_info == nil {
		return 0, fmt.Errorf("未授权")
	}

	var user_info *user.Entity
	err = user.M.Where("id", uid).Scan(&user_info)
	if err != nil {
		return 0, err
	}

	if user_info == nil {
		return 0, fmt.Errorf("用户不存在")
	}

	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	hash_id := utils.GenUUID()
	_, err = tx.Model(case_info.Table).Data(g.Map{
		"serial_number":  serial_number,
		"type":           req.CaseType,
		"multi_group_id": req.MultiGroupId,
		"multi_index":    req.MultiIndex,
		"check_id":       req.CheckId,
		"check_type":     check_info.CheckType,
		"check_level":    check_info.Level,
		"hash_id":        hash_id,
		"user_id":        uid,
	}).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	var info case_info.Entity
	err = tx.Model(case_info.Table).Where(g.Map{"hash_id": hash_id}).Order("id desc").Scan(&info)
	if err != nil {
		tx.Rollback()
		return
	}

	srv_number, err := strconv.ParseInt(strings.TrimLeft(licence_info.AuthorNumber, "S_"), 10, 64)
	if err != nil {
		tx.Rollback()
		return
	}

	case_id = int64(srv_number*100000000 + info.Id)
	_, err = tx.Model(case_info.Table).Where(g.Map{"id": info.Id}).Data(g.Map{"id": case_id}).Update()
	if err != nil {
		tx.Rollback()
		return
	}
	return case_id, tx.Commit()
}

func (s *caseService) SaveScoreStandard(req *case_info.SaveScoreStandardReq, serial_number string) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Model(case_plane.Table).Where("case_id", req.CaseId).Delete()
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Model(case_plane_structure.Table).Where("case_id", req.CaseId).Delete()
	if err != nil {
		tx.Rollback()
		return
	}

	total_score := 0
	case_plane_list, case_plane_structure_list := g.Array{}, g.Array{}
	for index, v := range req.PlaneList {
		case_plane_list = append(case_plane_list, g.Map{
			"serial_number":     serial_number,
			"case_id":           req.CaseId,
			"shot_id":           "",
			"group_id":          v.GroupId,
			"plane_id":          v.PlaneId,
			"plane_index":       v.PlaneIndex,
			"plane_hash":        v.PlaneHash,
			"plane_name_ch":     v.PlaneNameCh,
			"plane_name_en":     v.PlaneNameEn,
			"plane_score":       0,
			"plane_total_score": v.PlaneScoreTotal,
			"sort":              index,
		})

		total_score += v.PlaneScoreTotal

		for index2, v2 := range v.StructureList {
			case_plane_structure_list = append(case_plane_structure_list, g.Map{
				"serial_number":         serial_number,
				"case_id":               req.CaseId,
				"shot_id":               "",
				"group_id":              v.GroupId,
				"plane_id":              v.PlaneId,
				"plane_hash":            v.PlaneHash,
				"plane_index":           v.PlaneIndex,
				"structure_id":          v2.StructureId,
				"structure_hash":        v2.StructureHash,
				"structure_name_ch":     v2.StructureNameCh,
				"structure_name_en":     v2.StructureNameEn,
				"structure_score":       0,
				"structure_score_total": v2.StructureScoreTotal,
				"sort":                  index2,
			})
		}
	}

	_, err = tx.Model(case_info.Table).Where(g.Map{"id": req.CaseId}).Data(g.Map{"total_score": total_score}).Update()
	if err != nil {
		tx.Rollback()
		return
	}

	// 保存切面
	if len(case_plane_list) > 0 {
		if _, err = tx.Model(case_plane.Table).Data(case_plane_list).Insert(); err != nil {
			tx.Rollback()
			return
		}

		// 保存切面部位
		if len(case_plane_structure_list) > 0 {
			if _, err = tx.Model(case_plane_structure.Table).Data(case_plane_structure_list).Insert(); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	return tx.Commit()
}

func (s *caseService) Update(req *case_info.UpdateCaseReq) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	data := g.Map{}
	if req.DoctorName != "" {
		data["doctor_name"] = req.DoctorName
	}
	if req.QcCaseDate != "" {
		date, _ := gtime.StrToTimeFormat(req.QcCaseDate, "Ymd")
		data["qc_case_date"] = date.Format("Y-m-d")
	}
	if req.DepartmentName != "" {
		data["department_name"] = req.DepartmentName
	}
	if req.CasePath != "" {
		data["path"] = req.CasePath
	}
	if req.CasePathCloud != "" {
		data["path_cloud"] = req.CasePathCloud
	}
	if req.JoinFlag != -1 {
		data["join_flag"] = req.JoinFlag
	}
	if req.CaseLength != -1 {
		data["length"] = req.CaseLength
	}
	if req.GaAua != -1 {
		data["ga_aua"] = req.GaAua
	}
	if req.GaLmp != -1 {
		data["ga_lmp"] = req.GaLmp
	}
	if req.LmpDate != -1 {
		data["lmp_date"] = req.LmpDate
	}

	if req.PatientId != "" {
		data["patient_id"] = req.PatientId
		//查询病人信息是否存在，存在则修改，不存在则添加
		var patient_info *case_patient.Entity
		if err := tx.Model(case_patient.Table).Where("patient_id", req.PatientId).Scan(&patient_info); err != nil {
			tx.Rollback()
			return err
		}

		if patient_info == nil {
			_, err = tx.Model(case_patient.Table).Data(g.Map{
				"patient_id":       req.PatientId,
				"name":             req.PatientName,
				"age":              req.PatientAge,
				"sex":              req.PatientSex,
				"accession_no":     req.PatientAccessionNo,
				"study_uid":        req.PatientStudyUid,
				"pat_local_id":     req.PatientPatLoaclId,
				"last_menses_date": req.PatientLastMensesDate,
			}).Insert()
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err = tx.Model(case_patient.Table).
				Where("patient_id", req.PatientId).
				Data(g.Map{
					"name":             req.PatientName,
					"age":              req.PatientAge,
					"sex":              req.PatientSex,
					"accession_no":     req.PatientAccessionNo,
					"study_uid":        req.PatientStudyUid,
					"pat_local_id":     req.PatientPatLoaclId,
					"last_menses_date": req.PatientLastMensesDate,
				}).Update()
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}
	if req.FetusName != "" {
		data["fetus_name"] = req.FetusName
	}

	if req.ResultStatus != -1 {
		data["result_status"] = req.ResultStatus
	}

	_, err = tx.Model(case_info.Table).Where("id", req.CaseId).Data(data).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新胎儿数量
	//if req.MultiGroupId != "" {
	//	fetus_num, err := tx.Model(case_info.Table).Where("multi_fetus_group_id", req.MultiFetusGroupId).Count()
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//
	//	_, err = tx.Model(case_info.Table).Where("multi_fetus_group_id", req.MultiFetusGroupId).Data(g.Map{"fetus_num": fetus_num}).Update()
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//}
	return tx.Commit()
}

// 删除病例相关数据
// 判断病例是否存在，如果存在 则判断是否为组病例（多胎病例、多人质控病例），如果为多胎病例 按照组ID删除
// 删除病例截图数据
// 删除病例切面、结构数据
// 删除病例得分数据
// 删除病例测量值数据
// 删除病例标签关联数据
func (s *caseService) Delete(where interface{}) (err error) {
	var case_list []*case_info.Entity
	if err = case_info.M.Unscoped().Where(where).Scan(&case_list); err != nil {
		return
	}

	if len(case_list) <= 0 {
		return
	}

	case_ids := []int64{}
	multi_group_ids := []string{}
	for _, v := range case_list {
		case_ids = append(case_ids, v.Id)
		if v.MultiGroupId != "" {
			multi_group_ids = append(multi_group_ids, v.MultiGroupId)
		}
	}

	if len(multi_group_ids) > 0 {
		var multi_case_list []*case_info.Entity
		if err = case_info.M.Unscoped().Where("multi_group_id", multi_group_ids).Scan(&multi_case_list); err != nil {
			return
		}

		for _, v := range multi_case_list {
			case_ids = append(case_ids, v.Id)
		}
	}

	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// 删除相关病例
	if _, err = tx.Model(case_info.Table).Unscoped().Where("id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除病例截图数据
	if _, err = tx.Model(case_shot.Table).Unscoped().Where("case_id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除病例切面、结构数据数据
	if _, err = tx.Model(case_plane.Table).Unscoped().Where("case_id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除病例切面、结构数据数据
	if _, err = tx.Model(case_plane_structure.Table).Unscoped().Where("case_id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	//  删除病例测量值数据
	if _, err = tx.Model(case_measured.Table).Unscoped().Where("case_id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除病例标签关联数据
	if _, err = tx.Model(case_label_relate.Table).Unscoped().Where("case_id", case_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

func (s *caseService) Finish(req *case_info.FinishCaseReq) (err error) {
	data := g.Map{"length": req.CaseLength, "stop_at": gtime.Datetime()}
	if req.CaseResult != -1 {
		data["result_status"] = req.CaseResult
	}
	_, err = case_info.M.Where("id", req.CaseId).Data(data).Update()
	return
}

var W = make(map[string]*gdb.Model)

func (s *caseService) GetTimeLine(req *case_info.CaseTimelineReq, serial_number string) (res []*case_info.TimeLine, case_count, case_curr_count int, err error) {
	M := case_info.M_alias

	M = M.Where("ci.is_hidden", 0)
	M = M.Where("ci.type != ", 4)
	M = M.Where("ci.multi_index", 0)

	M = M.Where(fmt.Sprintf("ci.serial_number = '%s' OR (ci.serial_number != '%s' AND ci.stop_at IS NOT NULL AND ci.length > 0)", serial_number, serial_number))

	case_count, err = M.Count()
	if err != nil {
		return
	}

	if req.CheckLevel != "" {
		M = M.Where("ci.check_level IN(" + strings.Trim(req.CheckLevel, ",") + ")")
	}

	if req.JoinFlag != -1 {
		M = M.Where("ci.join_flag", req.JoinFlag)
	}

	if req.CaseResult != -1 {
		M = M.Where("ci.result_status", req.CaseResult)
	}

	if req.CaseType != "" {
		M = M.Where("ci.type IN(" + strings.Trim(req.CaseType, ",") + ")")
	}

	if req.CheckType != "" {
		M = M.Where("ci.check_type IN(" + strings.Trim(req.CheckType, ",") + ")")
	}

	if req.TimeBegin != "" {
		M = M.Where("ci.create_at >= ", req.TimeBegin+" 00:00:01")
	}

	if req.TimeEnd != "" {
		M = M.Where("ci.create_at <= ", req.TimeEnd+" 23:59:59")
	}

	if req.CaseLengthMin != -1 {
		M = M.Where("ci.length >= ", req.CaseLengthMin)
	}

	if req.CaseLengthMax != -1 {
		M = M.Where("ci.length <= ", req.CaseLengthMax)
	}

	if len(req.CaseLabels) > 0 {
		case_label_list, err := LabelService.RelateList(g.Map{"label_id": req.CaseLabels})
		if err != nil {
			g.Log().Error(err)
			return nil, 0, 0, err
		}
		if len(case_label_list) > 0 {
			case_ids := []int64{}
			for _, v := range case_label_list {
				case_ids = append(case_ids, v.CaseId)
			}

			M = M.WhereIn("ci.id", case_ids)
		}
	}

	if len(req.UserIds) > 0 {
		M = M.WhereIn("ci.length", req.UserIds)
	}

	M = M.LeftJoin("check c", "c.id = ci.check_id").Unscoped()
	W[serial_number] = M //缓存查询模型, 用以病例查询

	if err = M.Fields("to_char(ci.create_at, 'yyyy-mm-dd') as date").
		Group("date").
		Order("date DESC").
		Scan(&res); err != nil {
		return
	}

	case_curr_count, err = M.Count()
	return
}

//获取分页列表
func (s *caseService) List(req *case_info.CaseListReq, serial_number string, is_finish bool) (list []*case_info.Entity, err error) {
	M := W[serial_number]
	//if is_finish {
	//	M = M.WhereNotNull("ci.stop_at")
	//}
	// 'serial_number' == serial_number Or ('serial_number' != serial_number AND stop_at is not null AND length > 0)

	//M = M.Where(fmt.Sprintf("ci.serial_number = '%s' OR (ci.serial_number != '%s' AND ci.stop_at IS NOT NULL AND ci.length > 0)", serial_number, serial_number))
	if len(req.Dates) <= 0 {
		list = make([]*case_info.Entity, 0)
		return
	}

	where := ""
	for i, v := range req.Dates {
		where += fmt.Sprintf("(ci.create_at >= '%s 00:00:01' AND ci.create_at <= '%s 23:59:59')", v, v)
		if i < len(req.Dates)-1 {
			where += " or "
		}
	}

	if where != "" {
		M = M.Where(where)
	}

	M = M.LeftJoin("case_patient cp", "cp.patient_id = ci.patient_id")
	data, err := M.Fields("ci.*,u.username,c.level as check_level,c.name as check_name,cp.name as patient_name,cp.sex as patient_sex,cp.age as patient_age").
		LeftJoin("user u", "u.id = ci.user_id").Order("ci.id DESC").All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*case_info.Entity, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	for _, v := range list {
		if v.StopAt.String() != "" {
			v.StopTime = v.StopAt.Timestamp()
		}
		v.CreateTime = v.CreateAt.Timestamp()
		v.UpdateTime = v.UpdateAt.Timestamp()
	}
	return
}

//获取分页列表
func (s *caseService) Page(req *case_info.PageReqParams) (total int, list []*case_info.Entity, err error) {
	M := case_info.M_alias

	if req.KeyWord != "" {
		M = M.WhereOrLike("ci.serial_number", "%"+req.KeyWord+"%")
		M = M.WhereOrLike("ci.id", "%"+req.KeyWord+"%")
		M = M.WhereOrLike("ci.path", "%"+req.KeyWord+"%")
		M = M.WhereOrLike("ci.doctor_name", "%"+req.KeyWord+"%")
	}

	if req.CaseId > 0 {
		M = M.WhereGTE("ci.id", req.CaseId)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("ci.create_at", req.StartTime)
	}

	if req.EndTime != "" {
		M = M.WhereLTE("ci.create_at", req.EndTime)
	}

	if req.CheckLevel != "" {
		M = M.Where("ci.check_level IN(" + strings.Trim(req.CheckLevel, ",") + ")")
	}

	if req.CaseResult != -1 {
		M = M.Where("ci.result_status", req.CaseResult)
	}

	if req.CaseType != "" {
		M = M.Where("ci.type IN(" + strings.Trim(req.CaseType, ",") + ")")
	}

	if req.CheckType != "" {
		M = M.Where("ci.check_type IN(" + strings.Trim(req.CheckType, ",") + ")")
	}

	if req.TimeBegin != "" {
		M = M.Where("ci.create_at >= ", req.TimeBegin+" 00:00:01")
	}

	if req.TimeEnd != "" {
		M = M.Where("ci.create_at <= ", req.TimeEnd+" 23:59:59")
	}

	if req.CaseLengthMin != -1 {
		M = M.Where("ci.length >= ", req.CaseLengthMin)
	}

	if req.CaseLengthMax != -1 {
		M = M.Where("ci.length <= ", req.CaseLengthMax)
	}

	if len(req.UserIds) > 0 {
		M = M.WhereIn("ci.length", req.UserIds)
	}

	total, err = M.LeftJoin("check c", "c.id = ci.check_id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("ci." + req.Order + " " + req.Sort)
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("ci.*,c.level as check_level,c.name as check_name").
		LeftJoin("check c", "c.id = ci.check_id").Unscoped().
		All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*case_info.Entity, len(data))
	err = data.Structs(&list)
	return
}

//当天病例
func (s *caseService) ListJoin(serial_number string, req *case_info.CaseListJoinReq) (list []*case_info.Entity, err error) {
	M := case_info.M_alias

	M = M.Where("ci.is_hidden", 0)
	M = M.Where("ci.type", 1)

	M = M.Where(fmt.Sprintf("ci.serial_number = '%s' OR (ci.serial_number != '%s' AND ci.stop_at IS NOT NULL AND ci.length > 0)", serial_number, serial_number))

	if req.Date != "" {
		M = M.Where("ci.create_at >= ", req.Date+" 00:00:01")
		M = M.Where("ci.create_at <= ", req.Date+" 23:59:59")
	}

	if req.JoinFlag {
		M = M.Where("ci.join_flag", 1)
	}

	if req.ExistLabel {
		M = M.Where("ci.flag_count > ", 0)
	}

	M = M.LeftJoin("check c", "c.id = ci.check_id").Unscoped()
	M = M.LeftJoin("case_patient cp", "cp.patient_id = ci.patient_id")
	data, err := M.Fields("ci.*,c.level as check_level,c.name as check_name,cp.name as patient_name,cp.sex as patient_sex,cp.age as patient_age").
		Order("ci.id DESC").All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*case_info.Entity, len(data))
	err = data.Structs(&list)
	case_id_arr := []int64{}
	for _, v := range list {
		v.StopTime = v.StopAt.Timestamp()
		v.CreateTime = v.CreateAt.Timestamp()
		v.UpdateTime = v.UpdateAt.Timestamp()
		case_id_arr = append(case_id_arr, v.Id)
	}

	var case_label_list []*case_label_relate.Entity
	case_label_relate.M_alias.
		Fields("clr.*,cl.name as label_name,cl.type as label_type").
		LeftJoin("case_label cl", "cl.id = clr.label_id").
		WhereIn("clr.case_id", case_id_arr).
		Order("clr.id").
		Scan(&case_label_list)

	for _, v := range list {
		for _, v2 := range case_label_list {
			if v2.CaseId == v.Id {
				v.LabelList = append(v.LabelList, &case_info.CaseLabel{
					LabelId:   v2.LabelId,
					LabelName: v2.LabelName,
					LabelType: v2.LabelType,
				})
			}
		}
	}

	return
}

//病例添加标签
func (s *caseService) AddLabel(req *case_label_relate.CaseAddLabelReq, uid int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if _, err = tx.Model(case_label_relate.Table).Data(g.Map{
		"case_id":  req.CaseId,
		"label_id": req.LabelId,
		"user_id":  uid,
	}).Insert(); err != nil {
		tx.Rollback()
		return
	}

	if _, err = tx.Model(case_info.Table).
		Where("id", req.CaseId).
		Data(g.Map{"flag_count": gdb.Raw("flag_count + 1")}).
		Update(); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

//病例删除标签
func (s *caseService) DelLabel(req *case_label_relate.CaseAddLabelReq) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if _, err = tx.Model(case_label_relate.Table).Where(g.Map{
		"case_id":  req.CaseId,
		"label_id": req.LabelId,
	}).Delete(); err != nil {
		tx.Rollback()
		return
	}

	if _, err = tx.Model(case_info.Table).
		Where("id", req.CaseId).
		Data(g.Map{"flag_count": gdb.Raw("flag_count - 1")}).
		Update(); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

// 病例详细
func (s *caseService) Detail(case_id int64) (res *case_info.Entity, err error) {
	err = case_info.M_alias.
		Fields("ci.*,c.level as check_level,c.name as check_name,u.username").
		LeftJoin("check c", "c.id = ci.check_id").
		LeftJoin("user u", "u.id = ci.user_id").
		Where("ci.id", case_id).
		Scan(&res)
	if err != nil || res == nil {
		return
	}

	//如果是实时分析的病例则查询标签、测量值等附加信息
	if res.Type == 1 {
		var label_list []*case_label_relate.Entity
		err = case_label_relate.M_alias.
			Fields("clr.*,cl.name as label_name,cl.type as label_type").
			LeftJoin("case_label cl", "clr.label_id = cl.id").
			Where("case_id", case_id).
			Order("clr.id").Scan(&label_list)
		if err != nil {
			return
		}

		//获取标签
		case_label_list := []*case_info.CaseLabel{}
		for _, v := range label_list {
			case_label_list = append(case_label_list, &case_info.CaseLabel{
				LabelId:   v.LabelId,
				LabelName: v.LabelName,
				LabelType: v.LabelType,
			})
		}
		res.LabelList = case_label_list

		//获取测量值
		var measured_list []*case_measured.Entity
		if err := case_measured.M.Where("case_id", case_id).Order("sort,id").Scan(&measured_list); err != nil {
			return nil, err
		}
		if len(measured_list) > 0 {
			res.MeasuredList = measured_list
		}
	}

	//病例截图
	var shot_list []*case_shot.Entity
	if err := case_shot.M.Where("case_id", case_id).Order("sort,create_at").Scan(&shot_list); err != nil {
		return nil, err
	}
	if len(shot_list) > 0 {
		var case_plane_list []*case_plane.Entity
		if err := case_plane.M.Where("case_id", case_id).Scan(&case_plane_list); err != nil {
			return nil, err
		}

		auto_shot_count := 0
		for _, v := range shot_list {
			for _, v2 := range case_plane_list {
				if v.ShotId == v2.ShotId {
					v.PlaneList = append(v.PlaneList, v2)
					auto_shot_count++
				}
			}
		}
		res.ShotList = shot_list
		res.AutoShotCount = auto_shot_count
	}

	if res.StopAt.String() != "" {
		res.StopTime = res.StopAt.Timestamp()
	}
	res.CreateTime = res.CreateAt.Timestamp()
	res.UpdateTime = res.UpdateAt.Timestamp()
	return
}

// 设置测量值
func (s *caseService) SetMeasured(req []*case_measured.CaseSetMeasuredReq, serial_number string) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	for _, v := range req {
		var info *case_measured.Entity
		if err = tx.Model(case_measured.Table).Where(g.Map{
			"serial_number": serial_number,
			"case_id":       v.CaseId,
			"type":          v.MeasuredType,
			"name":          v.MeasuredName,
		}).Scan(&info); err != nil {
			tx.Rollback()
			return
		}

		if info == nil {
			if _, err = tx.Model(case_measured.Table).Data(g.Map{
				"serial_number": serial_number,
				"case_id":       v.CaseId,
				"type":          v.MeasuredType,
				"name":          v.MeasuredName,
				"value":         v.MeasuredValue,
				"value_pts":     v.MeasuredValuePts,
				"min":           v.MeasuredMin,
				"max":           v.MeasuredMax,
				"sort":          v.Sort,
			}).Insert(); err != nil {
				tx.Rollback()
				return
			}
		} else {
			if _, err = tx.Model(case_measured.Table).Where(g.Map{
				"serial_number": serial_number,
				"case_id":       v.CaseId,
				"type":          v.MeasuredType,
				"name":          v.MeasuredName,
			}).Data(g.Map{
				"value":     v.MeasuredValue,
				"value_pts": v.MeasuredValuePts,
				"min":       v.MeasuredMin,
				"max":       v.MeasuredMax,
				"sort":      v.Sort,
			}).Update(); err != nil {
				tx.Rollback()
				return
			}
		}
	}

	return tx.Commit()
}

//保存截图
func (s *caseService) SaveShot(req *case_shot.CaseSaveShotReq, serial_number string) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Model(case_shot.Table).Data(g.Map{
		"shot_id":       req.ShotId,
		"serial_number": serial_number,
		"case_id":       req.CaseId,
		"type":          req.ShotType,
		"name":          req.ShotName,
		"path":          req.ShotPath,
		"pts":           req.ShotPts,
		"plane_num":     len(req.PlaneList),
		"used_times":    0,
		"sort":          req.ShotSort,
	}).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	shot_plane_list := g.Array{}
	shot_plane_part_list := g.Array{}
	for _, v := range req.PlaneList {
		shot_plane_list = append(shot_plane_list, g.Map{
			"serial_number":     serial_number,
			"case_id":           req.CaseId,
			"shot_id":           req.ShotId,
			"group_id":          v.GroupId,
			"plane_id":          v.PlaneId,
			"plane_hash":        v.PlaneHash,
			"plane_name_ch":     v.PlaneNameCh,
			"plane_name_en":     v.PlaneNameEn,
			"plane_score":       v.PlaneScore,
			"plane_score_total": v.PlaneScoreTotal,
		})

		for _, v2 := range v.StructureList {
			shot_plane_part_list = append(shot_plane_part_list, g.Map{
				"serial_number":         serial_number,
				"case_id":               req.CaseId,
				"shot_id":               req.ShotId,
				"group_id":              v.GroupId,
				"plane_id":              v.PlaneId,
				"plane_hash":            v.PlaneHash,
				"structure_id":          v2.StructureId,
				"structure_hash":        v2.StructureHash,
				"structure_name_ch":     v2.StructureNameCh,
				"structure_name_en":     v2.StructureNameEn,
				"structure_score":       v2.StructureScore,
				"structure_score_total": v2.StructureScoreTotal,
			})
		}
	}

	// 保存截图分数
	if len(shot_plane_list) > 0 {
		if _, err = tx.Model(case_shot_plane.Table).Data(shot_plane_list).Insert(); err != nil {
			tx.Rollback()
			return
		}

		// 保存切面部位分数
		if len(shot_plane_part_list) > 0 {
			if _, err = tx.Model(case_shot_plane_structure.Table).Data(shot_plane_part_list).Insert(); err != nil {
				tx.Rollback()
				return
			}
		}

	}

	return tx.Commit()
}

// 删除截图
func (s *caseService) DeleteShot(shot_ids []string) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// 删除截图
	if _, err = tx.Model(case_shot.Table).Unscoped().Where("shot_id", shot_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除截图切面得分
	if _, err = tx.Model(case_shot_plane.Table).Unscoped().Where("shot_id", shot_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	// 删除截图结构得分
	if _, err = tx.Model(case_shot_plane_structure.Table).Unscoped().Where("shot_id", shot_ids).Delete(); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

//保存截图  废弃
func (s *caseService) SaveShotScore(req *case_shot_plane.ShotScoreSaveReq, serial_number string) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	shot_plane_list := g.Array{}
	shot_plane_part_list := g.Array{}
	for _, v := range req.PlaneList {
		shot_plane_list = append(shot_plane_list, g.Map{
			"serial_number":     serial_number,
			"case_id":           req.CaseId,
			"shot_id":           req.ShotId,
			"group_id":          v.GroupId,
			"plane_id":          v.PlaneId,
			"plane_hash":        v.PlaneHash,
			"plane_name_ch":     v.PlaneNameCh,
			"plane_name_en":     v.PlaneNameEn,
			"plane_score":       v.PlaneScore,
			"plane_score_total": v.PlaneScoreTotal,
		})

		for _, v2 := range v.StructureList {
			shot_plane_part_list = append(shot_plane_part_list, g.Map{
				"serial_number":         serial_number,
				"case_id":               req.CaseId,
				"shot_id":               req.ShotId,
				"group_id":              v.GroupId,
				"plane_id":              v.PlaneId,
				"plane_hash":            v.PlaneHash,
				"structure_id":          v2.StructureId,
				"structure_hash":        v2.StructureHash,
				"structure_name_ch":     v2.StructureNameCh,
				"structure_name_en":     v2.StructureNameEn,
				"structure_score":       v2.StructureScore,
				"structure_score_total": v2.StructureScoreTotal,
			})
		}
	}

	// 保存截图分数
	if _, err = tx.Model(case_shot_plane.Table).Data(shot_plane_list).Insert(); err != nil {
		tx.Rollback()
		return
	}

	// 保存切面部位分数
	if _, err = tx.Model(case_shot_plane_structure.Table).Data(shot_plane_part_list).Insert(); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

//截图使用
func (s *caseService) ShotUsed(req *case_plane.ShotUsedReq) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var shot_palce_score_list []*case_shot_plane.Entity
	if err = tx.Model(case_shot_plane.Table).Where("case_id", req.CaseId).Scan(&shot_palce_score_list); err != nil {
		tx.Rollback()
		return
	}

	var shot_plane_structure_score_list []*case_shot_plane_structure.Entity
	if err = tx.Model(case_shot_plane_structure.Table).Where("case_id", req.CaseId).Scan(&shot_plane_structure_score_list); err != nil {
		tx.Rollback()
		return
	}

	if req.IsReset {
		//重置病例切面数据
		if _, err = tx.Model(case_plane.Table).
			Where(g.Map{"case_id": req.CaseId}).
			Data(g.Map{"shot_id": "", "plane_score": 0}).
			Update(); err != nil {
			tx.Rollback()
			return
		}

		if _, err = tx.Model(case_plane_structure.Table).
			Where(g.Map{"case_id": req.CaseId}).
			Data(g.Map{"shot_id": "", "structure_score": 0}).
			Update(); err != nil {
			tx.Rollback()
			return
		}
	}

	shot_ids := []string{}
	for _, v := range req.PlaneShotList {
		shot_ids = append(shot_ids, v.ShotId)
		plane_score := 0
		for _, v2 := range shot_palce_score_list {
			if v.ShotId == v2.ShotId && v.PlaneId == v2.PlaneId {
				plane_score = v2.PlaneScore
			}
		}

		// 更新切面得分
		if _, err = tx.Model(case_plane.Table).
			Where(g.Map{"case_id": req.CaseId, "plane_id": v.PlaneId, "plane_index": v.PlaneIndex}).
			Data(g.Map{"shot_id": v.ShotId, "plane_score": plane_score}).
			Update(); err != nil {
			tx.Rollback()
			return
		}

		//更新结构得分
		for _, v2 := range shot_plane_structure_score_list {
			if v2.PlaneId == v.PlaneId && v.ShotId == v2.ShotId {
				if _, err = tx.Model(case_plane_structure.Table).
					Where(g.Map{"case_id": req.CaseId, "plane_id": v.PlaneId, "structure_id": v2.StructureId}).
					Data(g.Map{"shot_id": v.ShotId, "structure_score": v2.StructureScore}).
					Update(); err != nil {
					tx.Rollback()
					return
				}
			}
		}
	}

	// 更新截图使用次数
	var case_plane_count_list []*struct {
		ShotId string
		Count  int
	}
	if err = tx.Model(case_plane.Table).
		Fields("shot_id,Count(*) as count").
		Where(g.Map{"case_id": req.CaseId, "shot_id": shot_ids}).
		Group("shot_id").Scan(&case_plane_count_list); err != nil {
		tx.Rollback()
		return
	}

	for _, v := range case_plane_count_list {
		if _, err = tx.Model(case_shot.Table).
			Where(g.Map{"case_id": req.CaseId, "shot_id": v.ShotId}).
			Data(g.Map{"used_times": v.Count}).Update(); err != nil {
			tx.Rollback()
			return
		}
	}

	case_score, err := tx.Model(case_plane.Table).Where(g.Map{"case_id": req.CaseId}).Sum("plane_score")
	if err != nil {
		tx.Rollback()
		return
	}

	// 修改病例总分
	if _, err = tx.Model(case_info.Table).
		Where(g.Map{"id": req.CaseId}).
		Data(g.Map{"score": case_score}).
		Update(); err != nil {
		tx.Rollback()
		return
	}

	// todo
	// 更新远程质控批次总分

	return tx.Commit()
}

//病例切面列表
func (s *caseService) PlaneScoreDetail(case_id int64) (list []*case_plane.Entity, err error) {
	if err = case_plane.M.Where("case_id", case_id).Order("sort").Scan(&list); err != nil {
		return
	}

	var structure_list []*case_plane_structure.Entity
	if err = case_plane_structure.M.Where("case_id", case_id).Order("sort").Scan(&structure_list); err != nil {
		return
	}

	var shot_list []*case_shot.Entity
	if err = case_shot.M.Where("case_id", case_id).Scan(&shot_list); err != nil {
		return
	}

	for _, v := range list {
		for _, v2 := range structure_list {
			if v.PlaneId == v2.PlaneId {
				v.StructureList = append(v.StructureList, v2)
			}
		}
		for _, v2 := range shot_list {
			if v.ShotId == v2.ShotId {
				v.ShotPath = v2.Path
				v.ShotName = v2.Name
			}
		}
	}
	return
}

// 组病例
func (s *caseService) MultiList(multi_group_id string) (list []*case_info.Entity, max_score, min_score, averag_score, median_score int, err error) {
	M := case_info.M_alias
	M = M.Where("ci.multi_group_id", multi_group_id)
	M = M.LeftJoin("case_patient cp", "cp.patient_id = ci.patient_id")
	M = M.LeftJoin("check c", "c.id = ci.check_id").Unscoped()
	M = M.LeftJoin("user u", "u.id = ci.user_id").Unscoped()
	data, err := M.Fields("ci.*,u.username,c.level as check_level,c.name as check_name,cp.name as patient_name,cp.sex as patient_sex,cp.age as patient_age").
		Order("ci.score DESC,ci.id").All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*case_info.Entity, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	len := len(list)
	if len > 0 {
		max_score = list[0].Score
		min_score = list[len-1].Score
		total_score := 0
		for _, v := range list {
			total_score += v.Score
			if v.StopAt.String() != "" {
				v.StopTime = v.StopAt.Timestamp()
			}
			v.CreateTime = v.CreateAt.Timestamp()
			v.UpdateTime = v.UpdateAt.Timestamp()
		}

		averag_score = total_score / len

		m := len % 2
		if m == 0 {
			median_score = (list[len/2-1].Score + list[len/2].Score) / 2
		} else {
			median_score = list[len/2].Score
		}
	}

	return
}

//查询切面
func (s *caseService) GetPlaneList(where interface{}) (info *case_plane.Entity, err error) {
	err = case_plane.M.Where(where).Scan(&info)
	return
}

//查询结构
func (s *caseService) GetStructureList(where interface{}) (list []*case_plane_structure.Entity, err error) {
	err = case_plane_structure.M.Where(where).Scan(&list)
	return
}

//修改分数
func (s *caseService) ModifyScore(req *case_plane.ModifyScoreReq, plane_info *case_plane.Entity, structure_list []*case_plane_structure.Entity, uid int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var case_info_ *case_info.Entity
	if err = tx.Model(case_info.Table).Where("case_id", req.CaseId).Scan(&case_info_); err != nil {
		tx.Rollback()
		return
	}

	//添加日志记录
	data := g.Array{g.Map{
		"type":        1,
		"case_id":     req.CaseId,
		"score":       case_info_.Score,
		"total_score": case_info_.TotalScore,
		"user_id":     uid,
		"remark":      req.Remark,
	}, g.Map{
		"type":          2,
		"case_id":       req.CaseId,
		"group_id":      req.GroupId,
		"plane_id":      req.PlaneId,
		"plane_hash":    plane_info.PlaneHash,
		"plane_index":   plane_info.PlaneIndex,
		"plane_name_ch": plane_info.PlaneNameCh,
		"plane_name_en": plane_info.PlaneNameEn,
		"score":         plane_info.PlaneScore,
		"total_score":   plane_info.PlaneTotalScore,
		"user_id":       uid,
		"remark":        req.Remark,
	}}

	for _, v := range structure_list {
		data = append(data, g.Map{
			"type":              3,
			"case_id":           req.CaseId,
			"group_id":          req.GroupId,
			"plane_id":          req.PlaneId,
			"plane_hash":        plane_info.PlaneHash,
			"plane_index":       plane_info.PlaneIndex,
			"plane_name_ch":     plane_info.PlaneNameCh,
			"plane_name_en":     plane_info.PlaneNameEn,
			"structure_id":      v.StructureId,
			"structure_hash":    v.StructureHash,
			"structure_name_ch": v.StructureNameCh,
			"structure_name_en": v.StructureNameEn,
			"score":             v.StructureScore,
			"total_score":       v.StructureScoreTotal,
			"user_id":           uid,
			"remark":            req.Remark,
		})
	}
	if _, err = tx.Model(qc_score_log.Table).Data(data).Insert(); err != nil {
		tx.Rollback()
		return
	}

	// 修改结构分数
	for _, v := range req.ScoreItem {
		if _, err = tx.Model(case_plane_structure.Table).Where(g.Map{
			"case_id":      req.CaseId,
			"group_id":     req.GroupId,
			"plane_id":     req.PlaneId,
			"structure_id": v.StructureId,
		}).Data(g.Map{"structure_score": v.Score}).Update(); err != nil {
			tx.Rollback()
			return
		}
	}

	// 修改切面分数
	structure_score, err := tx.Model(case_plane_structure.Table).Where(g.Map{
		"case_id":  req.CaseId,
		"group_id": req.GroupId,
		"plane_id": req.PlaneId,
	}).Sum("structure_score")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = tx.Model(case_plane.Table).Where(g.Map{
		"case_id":  req.CaseId,
		"group_id": req.GroupId,
		"plane_id": req.PlaneId,
	}).Data(g.Map{"plane_score": structure_score}).Update(); err != nil {
		tx.Rollback()
		return
	}

	// 修改病例分数
	plane_score, err := tx.Model(case_plane.Table).Where(g.Map{
		"case_id":     req.CaseId,
		"plane_index": 0,
	}).Sum("plane_score")
	if err != nil {
		tx.Rollback()
		return
	}
	if _, err = tx.Model(case_info.Table).Where(g.Map{
		"case_id": req.CaseId,
	}).Data(g.Map{"score": plane_score}).Update(); err != nil {
		tx.Rollback()
		return
	}

	return
}

// 修改病例影藏状态
func (s *caseService) ModifyCaseHidden(case_ids []int64, is_hidden int) (err error) {
	if !utils.InArray(is_hidden, []int{0, 1, 2}) {
		err = fmt.Errorf("状态值异常")
		return err
	}

	_, err = case_info.M.Where(g.Map{"case_id": case_ids}).Data(g.Map{"is_hidden": is_hidden}).Update()
	return
}

// 获取隐藏病例列表
func (s *caseService) CaseHiddenList(req *case_info.HiddenListReq) (total int, list []*case_info.Entity, err error) {
	M := case_info.M_alias

	M = M.Where("ci.is_hidden", 1)

	total, err = M.LeftJoin("check c", "c.id = ci.check_id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("ci." + req.Order + " " + req.Sort)
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("ci.*,c.level as check_level,c.name as check_name").
		LeftJoin("check c", "c.id = ci.check_id").Unscoped().
		All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*case_info.Entity, len(data))
	err = data.Structs(&list)
	return
}
