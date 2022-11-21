package device

import (
	"aiyun_local_srv/library/utils/convert"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/yumaojun03/dmidecode"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

//获取设备uuid
func GetDeviceUuid() (uuid string, err error) {
	systype := runtime.GOOS

	if systype == "windows" {
		cmd := exec.Command("CMD", "/C", "wmic csproduct get uuid")
		uuid_byte, err := cmd.Output()
		if err != nil {
			return "", err
		}

		regexStr := regexp.MustCompile(`\r\r\n(.*)\r\r\n`)
		params := regexStr.FindStringSubmatch(convert.ToString(uuid_byte))
		if len(params) > 1 {
			uuid = strings.TrimSpace(params[1])
		}
	}

	if systype == "linux" {
		dmi, err := dmidecode.New()
		if err != nil {
			return "", err
		}

		infos, err := dmi.System()
		if err != nil {
			return "", err
		}

		if len(infos) > 0 {
			uuid = infos[0].UUID
		} else {
			err = fmt.Errorf("Device uuid not found")
		}
	}
	return
}

// 获取设备其他信息
func GetDeviceInfo() (cup_name, physical_id, product_name string, err error) {
	dmi, err := dmidecode.New()
	if err != nil {
		return
	}

	boardInfo, err := dmi.BaseBoard()
	if err != nil {
		return
	}
	if len(boardInfo) > 0 {
		product_name = boardInfo[0].ProductName
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return
	}

	if len(cpuInfo) > 0 {
		cup_name = cpuInfo[0].ModelName
		physical_id = cpuInfo[0].PhysicalID
	}

	return
}
