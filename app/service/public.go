package service

import (
	"aiyun_local_srv/library/response"
	"aiyun_local_srv/library/utils/curl"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

var PublicService = new(publicService)

type publicService struct{}

type Hospital struct {
	Id         int    `json:"hospital_id"`
	Name       string `json:"hospital_name"`
	RegionId   string `json:"region_id"`
	RegionName string `json:"region_name"`
}

// 获取云端医院地址
func (s *publicService) CloudHospitalList(region_id string) (list []*Hospital, err error) {
	base_url := g.Cfg().GetString("cloud_server.BaseUrl")
	ret, err := curl.Get(base_url+"public/hospital", g.Map{
		"region_id": region_id,
	})
	if err != nil {
		return
	}

	var rsp *response.Response
	if err = gconv.Struct(ret, &rsp); err != nil {
		return
	}

	if rsp.Code != 0 {
		err = fmt.Errorf(rsp.Msg)
		return
	}

	if err = gconv.Struct(rsp.Data, &list); err != nil {
		return
	}
	return
}
