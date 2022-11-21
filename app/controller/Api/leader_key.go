package Api

import (
	"aiyun_local_srv/app/model/leader_key"
	"aiyun_local_srv/app/model/user"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/curl"
	"encoding/base64"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

var LeaderKey = leaderKeyApi{}

type leaderKeyApi struct{}

//检查主任秘钥授权状态
func (*leaderKeyApi) Status(r *ghttp.Request) {
	leader_info, err := service.LeaderKeyService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_info == nil {
		response.Json(r, 0, "未绑定主任秘钥", g.Map{"state": 0})
	}

	response.Json(r, 0, "已绑定主任秘钥", g.Map{
		"state":                1,
		"leader_serial_number": leader_info.SerialNumber,
		"bind_time":            leader_info.CreateAt,
	})
}

//绑定主任秘钥
func (*leaderKeyApi) Bind(r *ghttp.Request) {
	var req *leader_key.LeaderBindReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//检查服务端是否授权
	licence, err := service.LicenceService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence == nil {
		response.Error(r, "服务端未授权")
	}

	leader_key_info, err := service.LeaderKeyService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if leader_key_info != nil {
		response.Error(r, "已经绑定主任秘钥，请勿重复操作")
	}

	//base64解码
	user_data_str, err := base64.StdEncoding.DecodeString(req.UserData)
	if err != nil {
		response.ErrorSys(r, err)
	}
	var user_data []*user.UserData
	if err := json.Unmarshal(user_data_str, &user_data); err != nil {
		response.ErrorSys(r, err)
	}

	//云端绑定主任秘钥
	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"leader-key/bind", leader_key.CloudLeaderBindReq{
		LeaderSerialNumber: req.LeaderSerialNumber,
		AuthorNumber:       req.AuthorNumber,
		ServerAuthorNumber: licence.AuthorNumber,
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

	//保存绑定记录
	if err = service.LeaderKeyService.Create(req.LeaderSerialNumber, user_data); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//重新绑定
func (*leaderKeyApi) Rebind(r *ghttp.Request) {
	var req *leader_key.LeaderBindReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//检查服务端是否授权
	licence, err := service.LicenceService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence == nil {
		response.Error(r, "服务端未授权")
	}

	//base64解码
	user_data_str, err := base64.StdEncoding.DecodeString(req.UserData)
	if err != nil {
		response.ErrorSys(r, err)
	}
	var user_data []*user.UserData
	if err := json.Unmarshal(user_data_str, &user_data); err != nil {
		response.ErrorSys(r, err)
	}

	//云端绑定主任秘钥
	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"leader-key/rebind", leader_key.CloudLeaderBindReq{
		LeaderSerialNumber: req.LeaderSerialNumber,
		AuthorNumber:       req.AuthorNumber,
		ServerAuthorNumber: licence.AuthorNumber,
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

	//保存绑定记录
	if err = service.LeaderKeyService.Create(req.LeaderSerialNumber, user_data); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
