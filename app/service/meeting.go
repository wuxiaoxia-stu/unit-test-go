package service

import (
	"aiyun_local_srv/app/model/meeting_device"
	"github.com/gogf/gf/frame/g"
	"strings"
)

var MeetingService = new(meetingService)

type meetingService struct{}

// 设备注册
func (s *meetingService) DeviceReg(req *meeting_device.DeviceReg, uid int, ip string) (err error) {
	var info *meeting_device.Entity
	err = meeting_device.M.Where(g.Map{
		"serial_number": req.SerialNumber,
		"device_type":   req.DeviceType,
	}).Scan(&info)
	if err != nil {
		return
	}

	if info == nil {
		_, err = meeting_device.M.Data(g.Map{
			"user_id":        uid,
			"author_number":  req.AuthorNumber,
			"serial_number":  req.SerialNumber,
			"device_type":    req.DeviceType,
			"cover_image":    req.CoverImage,
			"doctor_name":    req.DoctorName,
			"region_id":      req.RegionId,
			"hospital_id":    req.HospitalId,
			"hospital_name":  req.HospitalName,
			"ip":             ip,
			"machine_brand":  req.MachineBrand,
			"floor":          req.Floor,
			"room":           req.Room,
			"machine_number": req.MachineNumber,
			"online":         1,
			"case_status":    1,
			"assist_status":  0,
		}).Insert()
	} else {
		_, err = meeting_device.M.Where(g.Map{
			"serial_number": req.SerialNumber,
			"device_type":   req.DeviceType,
		}).Data(g.Map{
			"user_id":        uid,
			"cover_image":    req.CoverImage,
			"doctor_name":    req.DoctorName,
			"region_id":      req.RegionId,
			"hospital_id":    req.HospitalId,
			"hospital_name":  req.HospitalName,
			"ip":             ip,
			"machine_brand":  req.MachineBrand,
			"floor":          req.Floor,
			"room":           req.Room,
			"machine_number": req.MachineNumber,
			"online":         1,
			"case_status":    1,
			"assist_status":  0,
		}).Update()
	}
	return
}

// 设备列表
func (s *meetingService) DeviceList(req *meeting_device.DeviceListReq) (list []*meeting_device.Entity, err error) {
	M := meeting_device.M_alias

	if req.Online != -1 {
		M = M.Where("md.online", req.Online)
	}

	if req.HospitalId > 0 {
		M = M.Where("md.hospital_id", req.HospitalId)
	} else {
		if req.HospitalName != "" {
			M = M.WhereOrLike("md.hospital_name", "%"+req.HospitalName+"%")
		} else if req.RegionId != "" {
			_region_id := strings.Split(req.RegionId, "")
			if len(_region_id) >= 6 {
				if _region_id[2] == "0" && _region_id[3] == "0" && _region_id[4] == "0" && _region_id[5] == "0" {
					M = M.WhereOrLike("md.region_id", _region_id[0]+_region_id[1]+"%")
				} else if _region_id[4] == "0" && _region_id[5] == "0" {
					M = M.WhereOrLike("md.region_id", _region_id[0]+_region_id[1]+_region_id[2]+_region_id[3]+"%")
				} else {
					M = M.Where("md.region_id", req.RegionId)
				}
			}
		}
	}

	err = M.Scan(&list)
	if err != nil {
		return
	}

	for _, v := range list {
		v.BaseUrl = g.Cfg().GetString("server.qc_server.Domain")
		if v.UserId == 0 {
			v.Avatar = "avatar/qc/default.jpeg"
		}
	}
	return
}
