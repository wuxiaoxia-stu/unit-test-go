// +----------------------------------------------------------------------
// 远程质控相关接口
// @Author sinbook <778774780@qq.com>
// @Copyright 版权所有 广州爱孕记信息科技有限公司 [ https://www.aiyunji.cn/ ]
// @Date 2022-10-26 10:01:30
// +----------------------------------------------------------------------
package qc

import (
	"aiyun_local_srv/app/model/qc_upload"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"strings"
)

//远程质控

var Qc = qcApi{}

type qcApi struct{}

//oss配置
func (*qcApi) OssOptions(r *ghttp.Request) {
	response.Success(r, g.Cfg().Get("AliOSS"))
}

//创建病例
func (*qcApi) CheckOptions(r *ghttp.Request) {
	parts_all, err := service.CheckService.InitPartAllData(2)
	if err != nil {
		response.Error(r, err.Error())
	}

	parts_defualt_check, err := service.CheckService.InitPartDefaultData()
	if err != nil {
		response.Error(r, err.Error())
	}

	list, err := service.CheckService.List(g.Map{"type": 2})
	if err != nil {
		response.Error(r, err.Error())
	}

	for _, v := range list {
		v.CurrCheckItem = strings.Split(v.CheckItem, ",")
		for _, v2 := range parts_defualt_check {
			if v.Level == v2.Level && v.CheckType == v2.CheckType && v.Type == v2.Type {
				v.DefaultCheckItem = v2.Default
			}
		}
	}

	response.Success(r, g.Map{
		"parts_all":  parts_all,
		"check_list": list,
	})
	response.Success(r)
}

//上传任务列表
func (*qcApi) UploadList(r *ghttp.Request) {
	var req *qc_upload.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.QcService.UploadPage(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

// 获取解压密码
func (*qcApi) GetUploadInfo(r *ghttp.Request) {
	upload_id := r.GetQueryString("upload_id")

	info, err := service.QcService.GetUploadInfo(g.Map{"upload_id": upload_id})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "上传任务不存在")
	}

	info.Password = service.QcService.PasswordConvert(info.Password)
	info.StatusText = qc_upload.UploadStatus[info.Status]
	response.Success(r, info)
}

//创建上传记录
func (*qcApi) CreateUpload(r *ghttp.Request) {
	var req *qc_upload.CreateUploadReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	upload_id, password, err := service.QcService.CreateUpload(req, 1)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{"upload_id": upload_id, "password": service.QcService.PasswordConvert(password)})
}

// 更新上传状态
func (*qcApi) UpdateUploadStatus(r *ghttp.Request) {
	var req *qc_upload.SetUploadStatusReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.QcService.UpdateUploadStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"upload_id":   req.UploadId,
		"status":      req.Status,
		"status_text": qc_upload.UploadStatus[req.Status],
	})
}

// 更新上传状态
func (*qcApi) UploadDelete(r *ghttp.Request) {
	var req *qc_upload.UploadDeleteReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.QcService.UploadDelete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 获取时间线
func (*qcApi) TimeLine(r *ghttp.Request) {
	var req *qc_upload.QcTimelineReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	time_line, err := service.QcService.GetTimeLine(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	var tree []*qc_upload.Year

	curr_year := ""
	curr_month := ""
	for _, v := range time_line {
		t_len := len(tree)
		d := strings.Split(v.Date, "-")
		if curr_year == d[0] && curr_month == d[1] {
			tree[t_len-1].Children[len(tree[t_len-1].Children)-1].Children =
				append(tree[t_len-1].Children[len(tree[t_len-1].Children)-1].Children, d[2])
		} else if curr_year == d[0] {
			month := &qc_upload.Month{
				Month:    d[1],
				Children: []string{d[2]},
			}
			tree[t_len-1].Children = append(tree[t_len-1].Children, month)
		} else {
			month := []*qc_upload.Month{
				&qc_upload.Month{
					Month:    d[1],
					Children: []string{d[2]},
				},
			}
			tree = []*qc_upload.Year{
				&qc_upload.Year{
					Year:     d[0],
					Children: month,
				},
			}
		}

		curr_year = d[0]
		curr_month = d[1]
	}

	response.Success(r, tree)
}

// 获取时间线
func (*qcApi) ReportList(r *ghttp.Request) {
	var req *qc_upload.QcTimelineReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	list, err := service.QcService.ReportList(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

// 单个批次病例列表
func (*qcApi) CaseList(r *ghttp.Request) {
	upload_id := r.GetQueryString("upload_id")

	list, err := service.QcService.CaseList(upload_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

// 单个批次病例列表
func (*qcApi) CaseDetail(r *ghttp.Request) {
	case_id := r.GetQueryInt64("case_id")

	info, err := service.QcService.CaseDetail(case_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, info)
}

// 获取质控中心分组数据列表
func (*qcApi) GroupList(r *ghttp.Request) {
	var req *qc_upload.QcGroupListReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	list, err := service.QcService.GroupList(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}
