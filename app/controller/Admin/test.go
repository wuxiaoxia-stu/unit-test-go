package Admin

import (
	"aiyun_local_srv/app/model/leader_key"
	"aiyun_local_srv/app/model/licence"
	"aiyun_local_srv/app/model/pre_pair"
	"aiyun_local_srv/app/model/user"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils"
	"aiyun_local_srv/library/utils/curl"
	"aiyun_local_srv/library/utils/rsa_crypt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

//测试使用接口

var Test = testApi{}

type testApi struct{}

//删除服务授权
func (*testApi) DelSrvAuthorize(r *ghttp.Request) {
	_, err := licence.M.Delete(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

type BindLeaderReq struct {
	LeaderSerialNumber string `json:"leader_serial_number" p:"leader_serial_number" v:"required#参数错误,主任秘钥系列号必填"` // 主任密钥唯一编号
	ServerAuthorNumber string `json:"server_author_number" p:"server_author_number" v:"required#参数错误：服务端授权码必填"`  // 服务端授权id
}

//删除主任秘钥绑定
func (*testApi) DelLeaderBind(r *ghttp.Request) {

	var bind_info *leader_key.Entity
	if err := leader_key.M.Limit(1).Scan(&bind_info); err != nil {
		response.ErrorDb(r, err)
	}
	if bind_info == nil {
		response.Error(r, "未绑定主任秘钥")
	}

	var licence_info *licence.Entity
	if err := licence.M.Limit(1).Scan(&licence_info); err != nil {
		response.ErrorDb(r, err)
	}
	if licence_info == nil {
		response.Error(r, "服务未授权")
	}

	// 删除云端绑定记录
	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Post(base_url+"leader-key/unbind", BindLeaderReq{
		LeaderSerialNumber: bind_info.SerialNumber,
		ServerAuthorNumber: licence_info.AuthorNumber,
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

	/**
		注意：此接口未测试使用接口，代码不严谨
		删除本地绑定记录
	 	删除主任秘钥用户数据
	*/
	if _, err := leader_key.M.Delete("status", 1); err != nil {
		response.ErrorDb(r, err)
	}

	if _, err := user.M.Delete(); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//解除配对
func (*testApi) BreakPair(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")
	if serial_number == "" {
		response.Error(r, "需要解除配对的客户端设备序列号必填")
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
	per_pair, err := service.PairService.Info(g.Map{"status": 1, "device_number": serial_number})
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
		ClientAuthorNumber: per_pair.AuthorNumber,
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

	if err := service.PairService.Update(g.Map{"status": 0}, g.Map{"author_number": per_pair.AuthorNumber}); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//取消激活
func (*testApi) CancelActivate(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")
	if serial_number == "" {
		response.Error(r, "需要解除配对的客户端设备序列号必填")
	}

	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Get(base_url+"authorize/cancel-activate", g.Map{"serial_number": serial_number})
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

	response.Success(r)
}

//取消激活
func (*testApi) GetPairList(r *ghttp.Request) {
	list := []*pre_pair.Entity{}
	if err := pre_pair.M.Where(g.Map{"status": 1}).Scan(&list); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}
