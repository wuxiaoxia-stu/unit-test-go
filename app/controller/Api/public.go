package Api

import (
	_const "aiyun_local_srv/app/const"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/cache"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/mojocn/base64Captcha"
	"time"
)

var Public = publicApi{}

type publicApi struct{}

//获取地区列表数据 is_tree为ture 返回树结构数据
func (*publicApi) Region(r *ghttp.Request) {
	is_tree := r.GetQueryBool("is_tree")
	region_id := r.GetQueryString("region_id")
	if is_tree {
		list, err := service.RegionService.Tree()
		if err != nil {
			response.ErrorSys(r, err)
		}
		response.Success(r, list)
	} else {
		l, err := service.RegionService.List()
		if err != nil {
			response.ErrorSys(r, err)
		}

		list := []*service.Region{}
		if region_id != "" {
			if region_id == "0" {
				region_id = ""
			}

			for _, v := range l {
				if v.ParentCode == region_id {
					list = append(list, &service.Region{
						Value:  v.Code,
						Label:  v.Name,
						Pinyin: v.Pinyin,
					})
				}
			}
		} else {
			for _, v := range l {
				list = append(list, &service.Region{
					Value:  v.Code,
					Label:  v.Name,
					Pinyin: v.Pinyin,
				})
			}
		}
		response.Success(r, list)
	}
}

//通过地区检索医院数据
func (*publicApi) Hospital(r *ghttp.Request) {
	region_id := r.GetQueryString("region_id", 0)

	list, err := service.PublicService.CloudHospitalList(region_id)
	if err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r, list)
}

// 获取验证码
var store = base64Captcha.DefaultMemStore

func (*publicApi) Captcha(r *ghttp.Request) {
	device := &base64Captcha.DriverDigit{
		Height:   60,
		Width:    150,
		Length:   4,
		MaxSkew:  0.1,
		DotCount: 25,
	}
	c := base64Captcha.NewCaptcha(device, store)
	captcha_id, base64, err := c.Generate()
	if err != nil {
		response.ErrorSys(r, err)
	}

	if err = cache.Set(_const.CAPTCHA_CODE_CACHE_KEY(captcha_id), c.Store.Get(captcha_id, true), time.Minute*5); err != nil {
		response.ErrorSys(r, err)
	}
	//if err = cache.Set(_const.CAPTCHA_CODE_CACHE_KEY("captcha_id"), "1234", time.Minute*5); err != nil {
	//	response.ErrorSys(r, err)
	//}
	response.Success(r, g.Map{
		"captcha_id": captcha_id,
		"base64":     base64,
	})
}
