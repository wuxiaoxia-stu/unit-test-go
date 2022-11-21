package Admin

import (
	"aiyun_local_srv/app/model/licence"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/curl"
	"aiyun_local_srv/library/utils/device"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

//授权相关接口

var Authorize = authorizeApi{}

type authorizeApi struct{}

//检测是否存在授权记录
//查询授权记录是否存， 如果不存在则进入授权流程
func (*authorizeApi) Check(r *ghttp.Request) {
	licence, err := service.LicenceService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if licence == nil {
		response.Json(r, 2, "未授权")
	}

	response.Success(r)
}

func (*authorizeApi) GetUuid(r *ghttp.Request) {
	uuid, err := device.GetDeviceUuid()
	if err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r, uuid)
}

//获取医院列表
func (*authorizeApi) HospitalList(r *ghttp.Request) {
	var req *licence.QueryHospital

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"hospital/list", g.Map{
		"page_size": 1000,
		"region_id": req.RegionId,
	})
	if err != nil {
		response.ErrorSys(r, err)
	}

	r.Response.WriteJson(ret)
}

//检查授权状态是否正常  1 可授权  0 授权无效
func (*authorizeApi) Status(r *ghttp.Request) {
	var req *licence.QueryUkey

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"authorize/status-srv", req)
	if err != nil {
		response.ErrorSys(r, err)
	}

	r.Response.WriteJson(ret)
}

//创建服务端授权
func (*authorizeApi) Create(r *ghttp.Request) {
	var req *licence.QueryUkey
	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.HospitalId <= 0 {
		response.Error(r, "被授权单位必填")
	}

	//获取硬件信息
	cpu_name, physical_id, product_name, err := device.GetDeviceInfo()
	if err != nil {
		response.ErrorSys(r, err)
	}

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"authorize/create-srv-licence", licence.LicenceCreate{
		UkeySerialNumber: req.SerialNumber,
		Uuid:             req.Uuid,
		Signature:        req.Signature,
		HospitalId:       req.HospitalId,
		State:            1,
		Role:             2,
		GpuName:          "",
		CpuName:          cpu_name,
		CpuId:            physical_id,
		CpuDeviceId:      physical_id,
		BaseboardID:      product_name,
		ClientVersion:    g.Cfg().GetString("server.Version"),
	})
	if err != nil {
		response.ErrorSys(r, err)
	}

	var rsp *response.Response
	if err = gconv.Struct(ret, &rsp); err != nil {
		response.ErrorSys(r, err)
	}
	if rsp.Code != 0 {
		response.Error(r, rsp.Msg)
	}

	var cloud_rsp *licence.LicenceCreateRsp
	if err = gconv.Struct(rsp.Data, &cloud_rsp); err != nil {
		response.ErrorSys(r, err)
	}

	_, err = service.LicenceService.Create(cloud_rsp)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
