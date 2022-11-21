package Api

import (
	"aiyun_local_srv/app/model/check"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"strconv"
	"strings"
)

//产筛相关接口

var Check = checkApi{}

type checkApi struct{}

//保存检查项
func (*checkApi) Save(r *ghttp.Request) {
	var req []*check.SaveCheckOptions

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.CheckService.Save(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//获取检查配置项
func (*checkApi) Options(r *ghttp.Request) {
	t := r.GetQueryInt("type", 1)

	parts_all, err := service.CheckService.InitPartAllData(t)
	if err != nil {
		response.Error(r, err.Error())
	}

	parts_defualt_check, err := service.CheckService.InitPartDefaultData()
	if err != nil {
		response.Error(r, err.Error())
	}

	list, err := service.CheckService.List(g.Map{"type": t})
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
}

// 获取单个检查数据
func (*checkApi) Info(r *ghttp.Request) {
	check_id := r.GetQueryString("check_id")

	info, err := service.CheckService.Info(g.Map{"id": check_id})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Error(r, "检查项不存在")
	}

	parts_all, err := service.CheckService.InitPartAllData(info.Type)
	if err != nil {
		response.Error(r, err.Error())
	}

	check_item := strings.Split(info.CheckItem, ",")
	part := parts_all[strconv.Itoa(info.CheckType)]
	for _, v := range part {
		for _, v2 := range check_item {
			if v.PlaneId == v2 {
				v.IsCheck = true
				break
			}
		}
	}

	response.Success(r, g.Map{
		"check_id":         info.Id,
		"check_name":       info.Name,
		"check_type":       info.CheckType,
		"check_level":      info.Level,
		"plane_check_list": part,
	})
}
