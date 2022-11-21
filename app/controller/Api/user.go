package Api

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/model/pre_pair"
	"aiyun_local_srv/app/model/user"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/cache"
	"aiyun_local_srv/library/utils/jwt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"time"
)

//授权相关接口

var User = userApi{}

type userApi struct{}

func (*userApi) Login(r *ghttp.Request) {
	var req *user.LoginReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	var pair_info *pre_pair.Entity
	err := pre_pair.M.Where(g.Map{"status": 1, "author_number": req.AuthorNumber}).Scan(&pair_info)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if pair_info == nil {
		response.Error(r, "客户端未配对")
	}

	//验证登录
	info, err := service.UserService.Info(g.Map{"id": req.UserId})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "用户不存在")
	}

	//校验密码
	if req.Password != info.Password {
		response.Error(r, "密码错误")
	}

	// 生成token
	token, err := jwt.GenerateCustomerLoginToken(info, pair_info.AuthorNumber, pair_info.DeviceNumber, r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("jwt.expires")

	//清除旧的token
	old_token, err := cache.Get(_const.TOKEN_CACHE_KEY(fmt.Sprintf("user:%d", info.Id)))
	if err != nil {
		response.ErrorSys(r, err)
	}
	if old_token != "" {
		if err := cache.Set(_const.TOKEN_CACHE_KEY("user:"+old_token), "0", time.Second*time.Duration(token_expires)); err != nil {
			response.ErrorSys(r, err)
		}

	}

	if err := cache.Set(_const.TOKEN_CACHE_KEY("user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	//记录用户id和token的映射关系
	if err := cache.Set(_const.TOKEN_CACHE_KEY(fmt.Sprintf("user:%d", info.Id)), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r, g.Map{
		"token":    token,
		"expire":   time.Now().Unix() + token_expires,
		"user_id":  info.Id,
		"username": info.Username,
	})
}

//刷新token
func (*userApi) RefreshToken(r *ghttp.Request) {
	//验证登录
	info, err := service.UserService.Info(g.Map{"id": r.GetCtxVar("uid").Int()})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "用户不存在")
	}

	// 生成token
	token, err := jwt.GenerateCustomerLoginToken(info, r.GetCtxVar("author_number").String(), r.GetCtxVar("serial_number").String(), r.GetClientIp())
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt64("jwt.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY("user:"+token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r, g.Map{
		"token":  token,
		"expire": time.Now().Unix() + token_expires,
	})
}

//退出登录
func (*userApi) Logout(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if err := cache.Del(_const.TOKEN_CACHE_KEY("user:" + token)); err != nil {
		response.ErrorSys(r, err)
	}
	response.Success(r)
}

//用户列表
func (*userApi) List(r *ghttp.Request) {
	user_list, err := service.UserService.List(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	list := []*user.UserList{}
	for _, v := range user_list {
		delete_at := 0
		if v.DeleteAt.String() != "" {
			delete_at = int(v.DeleteAt.Unix())
		}
		list = append(list, &user.UserList{
			Id:       v.Id,
			Username: v.Username,
			RoleType: v.RoleType,
			RoleName: user.RoleTypeOptions[v.RoleType],
			Password: v.Password,
			UpdateAt: int(v.UpdateAt.Unix()),
			DeleteAt: delete_at,
		})
	}

	response.Success(r, g.Map{"user_list": list})
}

//保存用户信息
func (*userApi) Save(r *ghttp.Request) {
	var req *user.SaveUserDataReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//base64解码
	user_data_str, err := base64.StdEncoding.DecodeString(req.UserData)
	if err != nil {
		response.ErrorSys(r, err)
	}

	var user_data []*user.UserData
	if err := json.Unmarshal(user_data_str, &user_data); err != nil {
		response.Error(r, "用户数据解析失败")
	}

	//校验用户名是否有重复
	var count map[string]int
	for _, v := range user_data {
		if v.Username == "" {
			response.Error(r, "用户名必填")
		}

		if v.DeleteAt == 0 {
			_, ok := count[v.Username]
			if ok {
				response.Error(r, "存在重复用户名")
			} else {
				count[v.Username] = 1
			}
		}
	}

	if err := service.UserService.Save(user_data); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

////修改用户信息  用户名 密码 状态
//func (*customerApi) Update(r *ghttp.Request) {
//
//}
//
////设置角色
//func (*customerApi) SetRole(r *ghttp.Request) {
//
//}

//获取角色清单
func (*userApi) RoleOptions(r *ghttp.Request) {
	response.Success(r, user.RoleTypeOptions)
}

//删除用户
func (*userApi) Delete(r *ghttp.Request) {
	user_id := r.GetQueryInt("user_id")
	_, err := service.UserService.Delete(user_id)
	if err != nil {
		response.ErrorDb(r, err)
	}
	response.Success(r)
}
