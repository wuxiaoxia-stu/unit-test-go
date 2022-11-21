package Api

import (
	"aiyun_local_srv/app/model/pre_pair"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils"
	"aiyun_local_srv/library/utils/curl"
	"aiyun_local_srv/library/utils/rsa_crypt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

var Pair = pairApi{}

type pairApi struct{}

// 处理配对
func (*pairApi) Apply(r *ghttp.Request) {
	var req *pre_pair.PairClientReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询配对记录是否存在
	per_pair, err := service.PairService.Info(g.Map{"status": 1, "author_number": req.AuthorNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if per_pair != nil {
		response.Error(r, "已经配对过，请勿重复操作")
	}

	//检查服务端是否授权
	licence, err := service.LicenceService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence == nil {
		response.Error(r, "服务端未授权")
	}

	pubKey, err := rsa_crypt.LoadPublicKeyBase64(licence.LicenceKey)
	if err != nil {
		response.Error(r, "服务端授权异常")
	}

	post_data := pre_pair.PairReq{
		ClientAuthorNumber: req.AuthorNumber,
		ClientSerialNumber: req.SerialNumber,
		ServerAuthorNumber: licence.AuthorNumber,
	}
	post_data.Signature = rsa_crypt.Crypt(pubKey, utils.GenSignMsg(post_data))

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"pair/apply", post_data)
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

	var data *pre_pair.Entity
	if err = gconv.Struct(rsp.Data, &data); err != nil {
		response.ErrorSys(r, err)
	}

	//配对成功,添加配对记录
	data.Ip = r.GetClientIp()

	if _, err := service.PairService.Create(data); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"server_author_number": licence.AuthorNumber,
		"state":                1,
		"pair_time":            gtime.Datetime(),
	})
}

// 解除配对
func (*pairApi) Break(r *ghttp.Request) {
	var req *pre_pair.BreakPairClientReq

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

	//查询配对记录是否存在
	per_pair, err := service.PairService.Info(g.Map{"status": 1, "author_number": req.AuthorNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if per_pair == nil {
		response.Error(r, "未检测到本地服务端配对信息")
	}

	pubKey, err := rsa_crypt.LoadPublicKeyBase64(licence.LicenceKey)
	if err != nil {
		response.Error(r, "服务端授权异常")
	}

	post_data := pre_pair.PairReq{
		ClientAuthorNumber: req.AuthorNumber,
		ServerAuthorNumber: licence.AuthorNumber,
	}
	post_data.Signature = rsa_crypt.Crypt(pubKey, utils.GenSignMsg(post_data))

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"pair/break", post_data)
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

	if err := service.PairService.Update(g.Map{"status": 0}, g.Map{"author_number": req.AuthorNumber}); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{"state": 0})
}

//检查配对状态
func (*pairApi) Status(r *ghttp.Request) {
	author_number := r.GetQueryString("author_number")

	licence, err := service.LicenceService.Info(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if licence == nil {
		response.Json(r, 2, "服务端未授权")
	}

	info, err := service.PairService.Info(g.Map{"author_number": author_number, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Json(r, 0, "未配对", g.Map{"state": 0, "server_author_number": nil})
	}

	response.Json(r, 0, "已配对", g.Map{
		"state":                1,
		"server_author_number": licence.AuthorNumber,
		"pair_time":            info.CreateAt,
	})
}
