package service

import (
	"aiyun_local_srv/app/model/case_info"
	"aiyun_local_srv/app/model/case_plane"
	"aiyun_local_srv/app/model/case_shot"
	"aiyun_local_srv/app/model/qc_upload"
	"aiyun_local_srv/app/model/qc_user"
	"aiyun_local_srv/library/utils"
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var QcService = new(qcService)

type qcService struct{}

func (s *qcService) GetUploadInfo(where interface{}) (info *qc_upload.Entity, err error) {
	err = qc_upload.M_alias.
		Fields("qu.*,c.name as check_name,c.type as check_type,c.level as check_level").
		LeftJoin("check c", "c.id = qu.check_id").
		Where(where).
		Scan(&info)
	return
}

//获取分页列表
func (s *qcService) UploadPage(req *qc_upload.PageReqParams) (total int, list []*qc_upload.Entity, err error) {
	M := qc_upload.M_alias

	M = M.Where("qu.status > ", 0)

	if req.KeyWord != "" {
	}

	if req.SerialNumber != "" {
		M = M.Where("qu.serial_number", req.SerialNumber)
	}

	if req.UserId > 0 {
		M = M.Where("qu.user_id", req.UserId)
	}

	if req.Status > 0 {
		M = M.Where("qu.status", req.Status)
	} else {
		M = M.WhereGT("qu.status", 0)
	}

	if req.TimeBegin != "" {
		M = M.WhereGTE("qu.create_at", req.TimeBegin)
	}

	if req.TimeEnd != "" {
		M = M.WhereLTE("qu.create_at", req.TimeEnd)
	}

	M = M.LeftJoin("check c", "c.id = qu.check_id")
	total, err = M.Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("qu." + req.Order + " " + req.Sort)
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("qu.*,c.name as check_name,c.type as check_type,c.level as check_level").All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*qc_upload.Entity, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	for _, v := range list {
		v.Password = s.PasswordConvert(v.Password)
		v.StatusText = qc_upload.UploadStatus[v.Status]
	}
	return
}

//密码转换算法
func (s *qcService) PasswordConvert(password string) string {
	return gmd5.MustEncryptString("zip_password:" + gmd5.MustEncryptString(password))
}

// 生成上传记录
func (s *qcService) CreateUpload(req *qc_upload.CreateUploadReq, uid int) (upload_id, password string, err error) {
	upload_id = utils.GenUUID()
	password = utils.GenUUID()
	_, err = qc_upload.M.Data(g.Map{
		"upload_id":     upload_id,
		"password":      password,
		"uid":           uid,
		"serial_number": req.SerialNumber,
		"client_type":   req.ClientType,
		"source_type":   req.SourceType,
		"check_id":      req.CheckId,
		"status":        0,
	}).Insert()
	return
}

//  设置上传状态
func (s *qcService) UpdateUploadStatus(req *qc_upload.SetUploadStatusReq) (err error) {
	data := g.Map{"status": req.Status}
	if req.Status == 1 {
		data["case_num"] = req.CaseNum
		data["img_num"] = req.ImgNum
		data["file_size"] = req.FileSize
		data["zip_at"] = gtime.Datetime()
	}

	time_map := map[int]string{
		3: "upload_at",
		5: "download_at",
		6: "case_create_at",
	}
	data[time_map[req.Status]] = gtime.Datetime()

	_, err = qc_upload.M.Where("upload_id", req.UploadId).Data(data).Update()
	return
}

//  删除记录
func (s *qcService) UploadDelete(req *qc_upload.UploadDeleteReq) (err error) {
	_, err = qc_upload.M.Where("upload_id", req.UploadIds).Delete()
	return
}

func (s *qcService) GetTimeLine(req *qc_upload.QcTimelineReq) (res []*qc_upload.TimeLine, err error) {
	M := qc_upload.M_alias

	M = M.Where("qu.status", 6)

	if req.CheckId != "" {
		M = M.Where("qu.check_id", req.CheckId)
	}

	if req.SourceType > 0 {
		M = M.Where("qu.source_type", req.SourceType)
	}

	if req.TimeBegin != "" {
		M = M.Where("qu.case_create_at >= ", req.TimeBegin+" 00:00:01")
	}

	if req.TimeEnd != "" {
		M = M.Where("qu.case_create_at <= ", req.TimeEnd+" 23:59:59")
	}

	if err = M.Fields("to_char(qu.case_create_at, 'yyyy-mm-dd') as date").
		LeftJoin("check c", "c.id = qu.check_id").
		Unscoped().
		Group("date").
		Order("date DESC").
		Scan(&res); err != nil {
		return
	}
	return
}

func (s *qcService) ReportList(req *qc_upload.QcTimelineReq) (list []*qc_upload.Entity, err error) {
	M := qc_upload.M_alias

	M = M.Where("qu.status", 6)

	if req.CheckId != "" {
		M = M.Where("qu.check_id", req.CheckId)
	}

	if req.SourceType > 0 {
		M = M.Where("qu.source_type", req.SourceType)
	}

	if req.TimeBegin != "" {
		M = M.Where("qu.case_create_at >= ", req.TimeBegin+" 00:00:01")
	}

	if req.TimeEnd != "" {
		M = M.Where("qu.case_create_at <= ", req.TimeEnd+" 23:59:59")
	}

	where := ""
	for i, v := range req.Dates {
		where += fmt.Sprintf("(qu.case_create_at >= '%s 00:00:01' AND qu.case_create_at <= '%s 23:59:59')", v, v)
		if i < len(req.Dates)-1 {
			where += " or "
		}
	}
	M = M.Where(where)

	data, err := M.Fields("qu.*,c.name as check_name,c.type as check_type,c.level as check_level").
		LeftJoin("check c", "c.id = qu.check_id").OrderDesc("qu.case_create_at").All()
	if err != nil {
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*qc_upload.Entity, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	return
}

func (s *qcService) CaseList(upload_id string) (list []*qc_upload.QcCaseList, err error) {
	M := case_info.M

	M = M.Where("upload_id", upload_id)

	data, err := M.Fields("id as case_id,doctor_name,score,create_at as case_time").Order("score,id").All()
	if err != nil {
		return
	}

	list = make([]*qc_upload.QcCaseList, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	return
}

// 病例详细
func (s *qcService) CaseDetail(case_id int64) (res *case_info.Entity, err error) {
	err = case_info.M_alias.
		Fields("ci.*,c.level as check_level,c.name as check_name,u.username").
		LeftJoin("check c", "c.id = ci.check_id").
		LeftJoin("user u", "u.id = ci.user_id").
		Where("ci.id", case_id).
		Where("ci.type", 4).
		Scan(&res)
	if err != nil || res == nil {
		return
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

func (s *qcService) GroupList(req *qc_upload.QcGroupListReq) (list []*qc_upload.QcGroupList, err error) {
	var res []*qc_upload.QcGroupList

	where := g.Map{"qu.status": 6}
	if req.HospitalId > 0 {
		where["u.hospital_id"] = req.HospitalId
	} else if req.RegionId != "" {
		where["u.region_id"] = req.RegionId
	}

	if req.TimeBegin != "" {
		where["qu.case_create_at >="] = req.TimeBegin + " 00:00:01"
	}

	if req.TimeEnd != "" {
		where["qu.case_create_at <="] = req.TimeEnd + " 23:59:59"
	}

	err = qc_upload.M_alias.
		Fields("qu.upload_id,qu.user_id,qu.check_id,qu.source_type,qu.case_create_at,"+
			"qu.case_num,qu.succ_case_num,qu.total_score,qu.pre_case_score,"+
			"u.real_name,u.region_id,u.region_name,u.hospital_id,u.hospital_name,c.name as check_name").
		LeftJoin("qc_user u", "u.id = qu.user_id").
		LeftJoin("check c", "c.id = qu.check_id").
		Unscoped().
		Where(where).
		Order("qu.case_create_at DESC").
		Scan(&res)

	list = []*qc_upload.QcGroupList{}
	hospital_map := map[int]struct {
		index1 int
		index2 int
	}{}

	for _, v := range res {
		value, ok := hospital_map[v.HospitalId]
		if ok {
			_, ok2 := list[value.index1].DateGroupList[v.CaseCreateAt.Format("Y年m月")]
			if !ok2 {
				list[value.index1].DateGroupList = make(map[string][]*qc_upload.QcGroupList)
				date_group_list := []*qc_upload.QcGroupList{}
				date_group_list = append(date_group_list, &qc_upload.QcGroupList{
					CaseCreateAt:     v.CaseCreateAt,
					UserId:           v.UserId,
					Realname:         v.Realname,
					CheckId:          v.CheckId,
					CheckName:        v.CheckName,
					RegionId:         v.RegionId,
					RegionName:       v.Realname,
					HospitalId:       v.HospitalId,
					HospitalName:     v.HospitalName,
					TotalNum:         v.TotalNum,
					TotalCaseNum:     v.TotalCaseNum,
					TotalAveragScore: v.TotalAveragScore,
					TotalMedianScore: v.TotalMedianScore,
				})
				list[value.index1].DateGroupList[v.CaseCreateAt.Format("Y年m月")] = date_group_list
			} else {
				list[value.index1].DateGroupList[v.CaseCreateAt.Format("Y年m月")] =
					append(list[value.index1].DateGroupList[v.CaseCreateAt.Format("Y年m月")], &qc_upload.QcGroupList{
						CaseCreateAt:     v.CaseCreateAt,
						UserId:           v.UserId,
						Realname:         v.Realname,
						CheckId:          v.CheckId,
						CheckName:        v.CheckName,
						RegionId:         v.RegionId,
						RegionName:       v.Realname,
						HospitalId:       v.HospitalId,
						HospitalName:     v.HospitalName,
						TotalNum:         v.TotalNum,
						TotalCaseNum:     v.TotalCaseNum,
						TotalAveragScore: v.TotalAveragScore,
						TotalMedianScore: v.TotalMedianScore,
					})
			}
		} else {
			len := len(list)
			hospital_map[v.HospitalId] = struct {
				index1 int
				index2 int
			}{len, 0}

			list = append(list, v)
			date_group_list := []*qc_upload.QcGroupList{}
			date_group_list = append(date_group_list, &qc_upload.QcGroupList{
				CaseCreateAt:     v.CaseCreateAt,
				UserId:           v.UserId,
				Realname:         v.Realname,
				CheckId:          v.CheckId,
				CheckName:        v.CheckName,
				RegionId:         v.RegionId,
				RegionName:       v.Realname,
				HospitalId:       v.HospitalId,
				HospitalName:     v.HospitalName,
				TotalNum:         v.TotalNum,
				TotalCaseNum:     v.TotalCaseNum,
				TotalAveragScore: v.TotalAveragScore,
				TotalMedianScore: v.TotalMedianScore,
			})
			list[len].DateGroupList = make(map[string][]*qc_upload.QcGroupList)
			list[len].DateGroupList[v.CaseCreateAt.Format("Y年m月")] = date_group_list

		}
	}

	return
}

// 质控中心密码生成规则
func (s *qcService) QcUserPasswordRule(password string) string {
	return utils.EncodeMD5(password, utils.EncodeMD5(password, "qc_user_password"))
}

//用户注册
func (s *qcService) Register(req *qc_user.Register, region_name, hospital_name string) (err error) {
	_, err = qc_user.M.Data(g.Map{
		"username":         req.Username,
		"password":         s.QcUserPasswordRule(req.Password),
		"role_type":        1,
		"avatar":           req.Avatar,
		"id_no":            req.IdNo,
		"real_name":        req.RealName,
		"region_id":        req.RegionId,
		"region_name":      region_name,
		"hospital_id":      req.HospitalId,
		"hospital_name":    hospital_name,
		"positional_title": req.PositionalTitle,
	}).Insert()
	return
}

//用户登录
func (s *qcService) GetUserInfo(where interface{}) (info *qc_user.Entity, err error) {
	err = qc_user.M.Where(where).Scan(&info)
	if err != nil {
		return
	}

	if info != nil {
		info.BaseUrl = g.Cfg().GetString("server.qc_server.Domain")
	}
	return
}

//获取分页列表
func (s *qcService) Page(req *qc_user.PageParams) (total int, list []*qc_user.Entity, err error) {
	M := qc_user.M_alias

	if req.Keyword != "" {
		M = M.WhereOrLike("u.username", "%"+req.Keyword+"%")
		M = M.WhereOrLike("u.realname", "%"+req.Keyword+"%")
	}

	if req.TimeBegin != "" {
		M = M.WhereGTE("u.create_at", req.TimeBegin)
	}

	if req.TimeEnd != "" {
		M = M.WhereLTE("u.create_at", req.TimeEnd)
	}

	total, err = M.Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("u." + req.Order + " " + req.Sort)
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*qc_user.Entity, len(data))
	err = data.Structs(&list)
	return
}

// 修改用户基本信息
func (s *qcService) UserModify(req *qc_user.ModifyReq, uid int, region_name, hospital_name string) (err error) {
	_, err = qc_user.M.
		Where("id", uid).
		Data(g.Map{
			"avatar":           req.Avatar,
			"id_no":            req.IdNo,
			"real_name":        req.RealName,
			"region_id":        req.RegionId,
			"region_name":      region_name,
			"hospital_id":      req.HospitalId,
			"hospital_name":    hospital_name,
			"positional_title": req.PositionalTitle,
		}).Update()
	return
}

// 修改手机号
func (s *qcService) UserModifyUsername(username string, uid int) (err error) {
	_, err = qc_user.M.Where("id", uid).Data(g.Map{"username": username}).Update()
	return
}

// 找回密码
func (s *qcService) FindPwd(req *qc_user.FindPwdReq) (err error) {
	_, err = qc_user.M.
		Where("username", req.Username).
		Data(g.Map{
			"password": s.QcUserPasswordRule(req.Password),
		}).Update()
	return
}
