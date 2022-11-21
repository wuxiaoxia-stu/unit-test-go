// +----------------------------------------------------------------------
// 远程质控-远程会诊相关接口
// @Author sinbook <778774780@qq.com>
// @Copyright 版权所有 广州爱孕记信息科技有限公司 [ https://www.aiyunji.cn/ ]
// @Date 2022-10-26 10:01:30
// +----------------------------------------------------------------------
package qc

import (
	"aiyun_local_srv/app/model/meeting_device"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

//远程会诊

var Meeting = meetingApi{}

type meetingApi struct{}

//
func (*meetingApi) DeviceReg(r *ghttp.Request) {
	var req *meeting_device.DeviceReg

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	uid := r.GetCtxVar("uid").Int()
	g.Dump(uid)

	err := service.MeetingService.DeviceReg(req, uid, r.GetClientIp())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

func (*meetingApi) DeviceList(r *ghttp.Request) {
	var req *meeting_device.DeviceListReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	list, err := service.MeetingService.DeviceList(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}
