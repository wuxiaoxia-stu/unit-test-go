package main

import (
	"aiyun_local_srv/app/model/case_label"
	"aiyun_local_srv/library/utils/convert"
	"aiyun_local_srv/library/utils/curl"
	"aiyun_local_srv/library/utils/device"
	"aiyun_local_srv/library/utils/rsa_crypt"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/fatih/color"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gosuri/uitable"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func Test4(t *testing.T) {
	g.Dump(device.GetDeviceUuid())

	//dmi, err := dmidecode.New()
	//if err != nil {
	//	g.Log().Error(err)
	//	return
	//}
	//
	////infos, err := dmi.BIOS()
	//// 支持以下类型的解析
	////infos, err := dmi.BaseBoard()
	////infos, err := dmi.Chassis()
	//// dmi.MemoryArray()
	//// dmi.MemoryDevice()
	////infos, err := dmi.Onboard()
	//// dmi.PortConnector()
	//// dmi.Processor()
	//// dmi.ProcessorCache()
	//// dmi.Slot()
	//infos, err := dmi.System()
	//if err != nil {
	//	g.Log().Error(err)
	//	return
	//}
	//
	//for i := range infos {
	//	fmt.Println(infos[i])
	//}
}

//GetBaseBoardID 获取主板的id
func TestGetDeviceUuid(t *testing.T) {
	systype := runtime.GOOS

	if systype == "windows" {
		cmd := exec.Command("CMD", "/C", "wmic csproduct get uuid")
		uuid_byte, err := cmd.Output()
		if err != nil {
			return
		}

		regexStr := regexp.MustCompile(`\r\r\n(.*)\r\r\n`)
		params := regexStr.FindStringSubmatch(convert.ToString(uuid_byte))
		if len(params) > 1 {
			uuid := strings.TrimSpace(params[1])
			g.Dump(uuid)
		}

	}

	if systype == "linux" {

	}

}

func Test3(t *testing.T) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return
	}
	g.Dump(cpuInfo)

	hostInfo, err := host.Info()
	if err != nil {
		return
	}
	g.Dump(hostInfo)

	netinfo, err := net.IOCounters(true)
	if err != nil {
		return
	}
	for _, v := range netinfo {
		g.Dump(v)
	}
}

type Hackers struct {
	Name     string
	Birthday string
	Bio      int
}

func Test2(t *testing.T) {
	table := uitable.New()
	table.MaxColWidth = 50

	hackers := []*Hackers{&Hackers{
		Name:     "zhangsna1231233333333333333333333",
		Birthday: "546545",
		Bio:      20,
	}, &Hackers{
		Name:     "zhangsna",
		Birthday: "546512333333333333333312333333333333345",
		Bio:      20,
	}}

	table.AddRow("NAME", "BIRTHDAY", "BIO")
	for _, hacker := range hackers {
		table.AddRow(hacker.Name, hacker.Birthday, hacker.Bio)
	}
	fmt.Println(table)
	//var list  []*user.Entity
	//s := []*user.Entity{
	//	&user.Entity{},
	//	&user.Entity{},
	//	&user.Entity{},
	//	&user.Entity{},
	//}
	//list = append(list, &user.Entity{})
	//list = append(list, &user.Entity{})
	//list = append(list, &user.Entity{})
}

// 阿里云短信
func TestAliyunSms(t *testing.T) {
	//accessKeyId := "LTAI4FzkincM1m4ZefnX7CKi"
	//accessKeySecret := "LLK0msb9gLMeBoex0gpoEZirgHePt1"
	//endpoint := "dysmsapi.aliyuncs.com"
	//c := &client.Config{AccessKeyId: &accessKeyId, AccessKeySecret: &accessKeySecret, Endpoint: &endpoint}
	//
	//newClient, err := dysmsapi20170525.NewClient(c)
	//if err != nil {
	//	panic(err)
	//}
	//phoneNumber := "18388178881"
	//templateCode := "SMS_251016135"
	//signName := "广州爱孕记"
	//code := fmt.Sprintf("{\"code\":%d}", grand.N(10000, 999999))
	//request := &dysmsapi20170525.SendSmsRequest{
	//	PhoneNumbers:  &phoneNumber,
	//	TemplateCode:  &templateCode,
	//	SignName:      &signName,
	//	TemplateParam: &code,
	//}
	//sms, err := newClient.SendSms(request)
	//if err != nil {
	//	panic(err)
	//}
	//
	//g.Dump(sms)

}

func TestColorPrinit(t *testing.T) {
	// Print with default helper functions
	color.Cyan("Prints text in cyan.")

	// A newline will be appended automatically
	color.Blue("Prints %s in blue.", "text")

	// These are using the default foreground colors
	color.Red("We have red")
	color.Magenta("And many others ..")
}

func TestAliOSS(t *testing.T) {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "oss-cn-chengdu.aliyuncs.com"
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := "LTAI5t8xtAp1pqymYbC5zqid"
	accessKeySecret := "OfD1IXqghhNKDzVRTne6Xjl2hmecS2"
	bucketName := "umi-sinbook-cn"

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	objectName := "main.exe"

	//上传文件
	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	//err = bucket.PutObjectFromFile(objectName, "main.exe")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}

	// 解冻
	//err = bucket.RestoreObject(objectName)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}

	// 下载文件
	// <yourObjectName>从OSS下载文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。

	downloadedFileName := "QGR311kZ_2622742983_uh111d.exe"
	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 列举文件。
	//marker := ""
	//for {
	//	lsRes, err := bucket.ListObjects(oss.Marker(marker))
	//	if err != nil {
	//			fmt.Println("Error:", err)
	//			os.Exit(-1)
	//	}
	//	// 打印列举文件，默认情况下一次返回100条记录。
	//	for _, object := range lsRes.Objects {
	//		fmt.Println("Bucket: ", object.Key)
	//	}
	//	if lsRes.IsTruncated {
	//		marker = lsRes.NextMarker
	//	} else {
	//		break
	//	}
	//}

	// 删除文件。
	//err = bucket.DeleteObject(objectName)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
}

// 导入系统标签
func TestImportLabel(t *testing.T) {
	//list, _ := service.KlFeatureService.Tree(g.Map{"status": 1, "invisible": 0})
	//
	//sys_label_list := g.Array{}
	//for _, v := range list {
	//	for _, v2 := range v.Children {
	//		for _, v3 := range v2.Children {
	//			for _, v4 := range v3.Children {
	//				g.Dump(v4.Name)
	//				sys_label_list = append(sys_label_list, g.Map{
	//					"type":       1,
	//					"user_id":    0,
	//					"group_name": v.Name,
	//					"name":       v4.Name,
	//				})
	//			}
	//		}
	//	}
	//}

	json_label := "[{\"id\":1,\"label_id\":\"M007-U182\",\"label_name\":\"多部位翼状胬肉(颈部/下巴/腋窝/肘窝/踝关节)\",\"group_name\":\"皮肤及软组织\"},{\"id\":2,\"label_id\":\"M007-U202\",\"label_name\":\"NF增厚\",\"group_name\":\"皮肤及软组织\"},{\"id\":3,\"label_id\":\"M007-U183\",\"label_name\":\"NT增厚\",\"group_name\":\"皮肤及软组织\"},{\"id\":4,\"label_id\":\"M007-U181\",\"label_name\":\"淋巴管囊肿\",\"group_name\":\"皮肤及软组织\"},{\"id\":5,\"label_id\":\"M007-U177\",\"label_name\":\"皮肤血管瘤\",\"group_name\":\"皮肤及软组织\"},{\"id\":6,\"label_id\":\"M007-U176\",\"label_name\":\"皮下组织肿胀\",\"group_name\":\"皮肤及软组织\"},{\"id\":7,\"label_id\":\"M007-U179\",\"label_name\":\"胎儿水肿\",\"group_name\":\"皮肤及软组织\"},{\"id\":8,\"label_id\":\"M003-U098\",\"label_name\":\"完全性房室传导阻滞\",\"group_name\":\"心血管系统\"},{\"id\":9,\"label_id\":\"M003-U099\",\"label_name\":\"心动过缓\",\"group_name\":\"心血管系统\"},{\"id\":10,\"label_id\":\"M003-U100\",\"label_name\":\"心律失常\",\"group_name\":\"心血管系统\"},{\"id\":11,\"label_id\":\"M003-U086\",\"label_name\":\"肺静脉异位引流\",\"group_name\":\"心血管系统\"},{\"id\":12,\"label_id\":\"M003-U095\",\"label_name\":\"下腔静脉离断\",\"group_name\":\"心血管系统\"},{\"id\":13,\"label_id\":\"M003-U094\",\"label_name\":\"双上腔静脉\",\"group_name\":\"心血管系统\"},{\"id\":14,\"label_id\":\"M003-U085\",\"label_name\":\"肺动脉狭窄\",\"group_name\":\"心血管系统\"},{\"id\":15,\"label_id\":\"M003-U097\",\"label_name\":\"动脉导管早闭\",\"group_name\":\"心血管系统\"},{\"id\":16,\"label_id\":\"M003-U088\",\"label_name\":\"大动脉共干畸形\",\"group_name\":\"心血管系统\"},{\"id\":17,\"label_id\":\"M003-U089\",\"label_name\":\"大动脉钙化\",\"group_name\":\"心血管系统\"},{\"id\":18,\"label_id\":\"M003-U087\",\"label_name\":\"大动脉转位\",\"group_name\":\"心血管系统\"},{\"id\":19,\"label_id\":\"M003-U091\",\"label_name\":\"主动脉弓异常\",\"group_name\":\"心血管系统\"},{\"id\":20,\"label_id\":\"M003-U090\",\"label_name\":\"主动脉弓缩窄\",\"group_name\":\"心血管系统\"},{\"id\":21,\"label_id\":\"M003-U092\",\"label_name\":\"法洛氏四联症\",\"group_name\":\"心血管系统\"},{\"id\":22,\"label_id\":\"M003-U093\",\"label_name\":\"右室双出口\",\"group_name\":\"心血管系统\"},{\"id\":23,\"label_id\":\"M003-U081\",\"label_name\":\"左心发育不良\",\"group_name\":\"心血管系统\"},{\"id\":24,\"label_id\":\"M003-U080\",\"label_name\":\"单心室\",\"group_name\":\"心血管系统\"},{\"id\":25,\"label_id\":\"M003-U082\",\"label_name\":\"室间隔缺损\",\"group_name\":\"心血管系统\"},{\"id\":26,\"label_id\":\"M003-U083\",\"label_name\":\"房间隔缺损\",\"group_name\":\"心血管系统\"},{\"id\":27,\"label_id\":\"M003-U084\",\"label_name\":\"心内膜垫缺损\",\"group_name\":\"心血管系统\"},{\"id\":28,\"label_id\":\"M003-U078\",\"label_name\":\"异位心\",\"group_name\":\"心血管系统\"},{\"id\":29,\"label_id\":\"M003-U077\",\"label_name\":\"右位心\",\"group_name\":\"心血管系统\"},{\"id\":30,\"label_id\":\"M003-U075\",\"label_name\":\"心轴异常\",\"group_name\":\"心血管系统\"},{\"id\":31,\"label_id\":\"M003-U076\",\"label_name\":\"心脏位置异常\",\"group_name\":\"心血管系统\"},{\"id\":32,\"label_id\":\"M003-U096\",\"label_name\":\"心脏灶性强回声\",\"group_name\":\"心血管系统\"},{\"id\":33,\"label_id\":\"M003-U079\",\"label_name\":\"心脏横纹肌瘤\",\"group_name\":\"心血管系统\"},{\"id\":34,\"label_id\":\"M003-U074\",\"label_name\":\"心包积液\",\"group_name\":\"心血管系统\"},{\"id\":35,\"label_id\":\"M003-U073\",\"label_name\":\"心脏扩大\",\"group_name\":\"心血管系统\"},{\"id\":36,\"label_id\":\"M003-U072\",\"label_name\":\"心脏肥大\",\"group_name\":\"心血管系统\"},{\"id\":37,\"label_id\":\"M003-U071\",\"label_name\":\"心脏增大\",\"group_name\":\"心血管系统\"},{\"id\":38,\"label_id\":\"M003-U070\",\"label_name\":\"心脏畸形\",\"group_name\":\"心血管系统\"},{\"id\":39,\"label_id\":\"M004-U122\",\"label_name\":\"胎粪性腹膜炎\",\"group_name\":\"胸腹部\"},{\"id\":40,\"label_id\":\"M004-U121\",\"label_name\":\"肠管回声增强\",\"group_name\":\"胸腹部\"},{\"id\":41,\"label_id\":\"M004-U120\",\"label_name\":\"消化道异常\",\"group_name\":\"胸腹部\"},{\"id\":42,\"label_id\":\"M004-U117\",\"label_name\":\"单脐动脉\",\"group_name\":\"胸腹部\"},{\"id\":43,\"label_id\":\"M004-U104\",\"label_name\":\"脐膨出\",\"group_name\":\"胸腹部\"},{\"id\":44,\"label_id\":\"M004-U103\",\"label_name\":\"腹壁肿块\",\"group_name\":\"胸腹部\"},{\"id\":45,\"label_id\":\"M004-U102\",\"label_name\":\"腹壁缺损\",\"group_name\":\"胸腹部\"},{\"id\":46,\"label_id\":\"M004-U101\",\"label_name\":\"肢体体壁缺陷\",\"group_name\":\"胸腹部\"},{\"id\":47,\"label_id\":\"M004-U116\",\"label_name\":\"腹水\",\"group_name\":\"胸腹部\"},{\"id\":48,\"label_id\":\"M004-U123\",\"label_name\":\"肝内异常\",\"group_name\":\"胸腹部\"},{\"id\":49,\"label_id\":\"M004-U118\",\"label_name\":\"胃泡增大\",\"group_name\":\"胸腹部\"},{\"id\":50,\"label_id\":\"M004-U115\",\"label_name\":\"内脏增大\",\"group_name\":\"胸腹部\"},{\"id\":51,\"label_id\":\"M004-U119\",\"label_name\":\"内脏反位\",\"group_name\":\"胸腹部\"},{\"id\":52,\"label_id\":\"M004-U114\",\"label_name\":\"膈肌反向\",\"group_name\":\"胸腹部\"},{\"id\":53,\"label_id\":\"M004-U113\",\"label_name\":\"膈疝\",\"group_name\":\"胸腹部\"},{\"id\":54,\"label_id\":\"M004-U107\",\"label_name\":\"肺部回声增高\",\"group_name\":\"胸腹部\"},{\"id\":55,\"label_id\":\"M004-U106\",\"label_name\":\"肺发育不良\",\"group_name\":\"胸腹部\"},{\"id\":56,\"label_id\":\"M004-U108\",\"label_name\":\"双肺增大\",\"group_name\":\"胸腹部\"},{\"id\":57,\"label_id\":\"M004-U111\",\"label_name\":\"胸腺异常\",\"group_name\":\"胸腹部\"},{\"id\":58,\"label_id\":\"M004-U105\",\"label_name\":\"胸腔积液\",\"group_name\":\"胸腹部\"},{\"id\":59,\"label_id\":\"M004-U112\",\"label_name\":\"胎儿甲状腺肿\",\"group_name\":\"胸腹部\"},{\"id\":60,\"label_id\":\"M004-U110\",\"label_name\":\"食道畸形\",\"group_name\":\"胸腹部\"},{\"id\":61,\"label_id\":\"M004-U109\",\"label_name\":\"气道畸形\",\"group_name\":\"胸腹部\"},{\"id\":62,\"label_id\":\"M008-U187\",\"label_name\":\"胎动减少\",\"group_name\":\"生长发育\"},{\"id\":63,\"label_id\":\"M008-U186\",\"label_name\":\"器官肥大\",\"group_name\":\"生长发育\"},{\"id\":64,\"label_id\":\"M008-U185\",\"label_name\":\"巨大胎儿\",\"group_name\":\"生长发育\"},{\"id\":65,\"label_id\":\"M008-U184\",\"label_name\":\"IUGR\",\"group_name\":\"生长发育\"},{\"id\":66,\"label_id\":\"M006-U175\",\"label_name\":\"髂骨翼角增宽\",\"group_name\":\"骨骼系统\"},{\"id\":67,\"label_id\":\"M006-U174\",\"label_name\":\"髂骨异常\",\"group_name\":\"骨骼系统\"},{\"id\":68,\"label_id\":\"M006-U171\",\"label_name\":\"骶尾椎异常\",\"group_name\":\"骨骼系统\"},{\"id\":69,\"label_id\":\"M006-U173\",\"label_name\":\"肩胛骨异常\",\"group_name\":\"骨骼系统\"},{\"id\":70,\"label_id\":\"M006-U172\",\"label_name\":\"锁骨异常\",\"group_name\":\"骨骼系统\"},{\"id\":71,\"label_id\":\"M006-U168\",\"label_name\":\"颈椎异常\",\"group_name\":\"骨骼系统\"},{\"id\":72,\"label_id\":\"M006-U170\",\"label_name\":\"脊柱侧弯畸形\",\"group_name\":\"骨骼系统\"},{\"id\":73,\"label_id\":\"M006-U169\",\"label_name\":\"椎体骨化异常\",\"group_name\":\"骨骼系统\"},{\"id\":74,\"label_id\":\"M006-U155\",\"label_name\":\"骨骺回声异常\",\"group_name\":\"骨骼系统\"},{\"id\":75,\"label_id\":\"M006-U161\",\"label_name\":\"髋关节异常\",\"group_name\":\"骨骼系统\"},{\"id\":76,\"label_id\":\"M006-U162\",\"label_name\":\"膝关节异常\",\"group_name\":\"骨骼系统\"},{\"id\":77,\"label_id\":\"M006-U153\",\"label_name\":\"并腿畸形\",\"group_name\":\"骨骼系统\"},{\"id\":78,\"label_id\":\"M006-U154\",\"label_name\":\"截肢样肢体畸形\",\"group_name\":\"骨骼系统\"},{\"id\":79,\"label_id\":\"M006-U152\",\"label_name\":\"海豹肢畸形\",\"group_name\":\"骨骼系统\"},{\"id\":80,\"label_id\":\"M006-U166\",\"label_name\":\"足姿势异常\",\"group_name\":\"骨骼系统\"},{\"id\":81,\"label_id\":\"M006-U164\",\"label_name\":\"手姿势异常\",\"group_name\":\"骨骼系统\"},{\"id\":82,\"label_id\":\"M006-U148\",\"label_name\":\"肢体姿势异常\",\"group_name\":\"骨骼系统\"},{\"id\":83,\"label_id\":\"M006-U146\",\"label_name\":\"肢体完整且不对称\",\"group_name\":\"骨骼系统\"},{\"id\":84,\"label_id\":\"M006-U167\",\"label_name\":\"小指第二指节短小\",\"group_name\":\"骨骼系统\"},{\"id\":85,\"label_id\":\"M006-U163\",\"label_name\":\"多并缺指（趾）畸形\",\"group_name\":\"骨骼系统\"},{\"id\":86,\"label_id\":\"M006-U165\",\"label_name\":\"手脚水肿\",\"group_name\":\"骨骼系统\"},{\"id\":87,\"label_id\":\"M006-U160\",\"label_name\":\"单纯腓骨畸形\",\"group_name\":\"骨骼系统\"},{\"id\":88,\"label_id\":\"M006-U159\",\"label_name\":\"单纯股骨畸形\",\"group_name\":\"骨骼系统\"},{\"id\":89,\"label_id\":\"M006-U156\",\"label_name\":\"单纯肱骨畸形\",\"group_name\":\"骨骼系统\"},{\"id\":90,\"label_id\":\"M006-U158\",\"label_name\":\"单纯桡骨畸形\",\"group_name\":\"骨骼系统\"},{\"id\":91,\"label_id\":\"M006-U157\",\"label_name\":\"单纯尺骨畸形\",\"group_name\":\"骨骼系统\"},{\"id\":92,\"label_id\":\"M006-U151\",\"label_name\":\"长骨弯曲\",\"group_name\":\"骨骼系统\"},{\"id\":93,\"label_id\":\"M006-U150\",\"label_name\":\"四肢增粗\",\"group_name\":\"骨骼系统\"},{\"id\":94,\"label_id\":\"M006-U149\",\"label_name\":\"四肢细长\",\"group_name\":\"骨骼系统\"},{\"id\":95,\"label_id\":\"M006-U147\",\"label_name\":\"四肢短小\",\"group_name\":\"骨骼系统\"},{\"id\":96,\"label_id\":\"M006-U145\",\"label_name\":\"肋骨异常\",\"group_name\":\"骨骼系统\"},{\"id\":97,\"label_id\":\"M006-U144\",\"label_name\":\"肋骨短小\",\"group_name\":\"骨骼系统\"},{\"id\":98,\"label_id\":\"M006-U200\",\"label_name\":\"胸廓畸形\",\"group_name\":\"骨骼系统\"},{\"id\":99,\"label_id\":\"M006-U143\",\"label_name\":\"胸廓狭窄\",\"group_name\":\"骨骼系统\"},{\"id\":100,\"label_id\":\"M006-U142\",\"label_name\":\"颅骨发育不良\",\"group_name\":\"骨骼系统\"},{\"id\":101,\"label_id\":\"M006-U141\",\"label_name\":\"大头畸形\",\"group_name\":\"骨骼系统\"},{\"id\":102,\"label_id\":\"M006-U140\",\"label_name\":\"无颅畸形\",\"group_name\":\"骨骼系统\"},{\"id\":103,\"label_id\":\"M006-U139\",\"label_name\":\"头型异常\",\"group_name\":\"骨骼系统\"},{\"id\":104,\"label_id\":\"M001-U032\",\"label_name\":\"颜面异常\",\"group_name\":\"颜面\"},{\"id\":105,\"label_id\":\"M001-U031\",\"label_name\":\"前额异常\",\"group_name\":\"颜面\"},{\"id\":106,\"label_id\":\"M001-U030\",\"label_name\":\"下颌发育不良\",\"group_name\":\"颜面\"},{\"id\":107,\"label_id\":\"M001-U028\",\"label_name\":\"上颌发育不良\",\"group_name\":\"颜面\"},{\"id\":108,\"label_id\":\"M001-U027\",\"label_name\":\"小下颌畸形\",\"group_name\":\"颜面\"},{\"id\":109,\"label_id\":\"M001-U029\",\"label_name\":\"颌骨囊肿\",\"group_name\":\"颜面\"},{\"id\":110,\"label_id\":\"M001-U022\",\"label_name\":\"多发舌肿瘤\",\"group_name\":\"颜面\"},{\"id\":111,\"label_id\":\"M001-U021\",\"label_name\":\"舌后坠\",\"group_name\":\"颜面\"},{\"id\":112,\"label_id\":\"M001-U020\",\"label_name\":\"巨舌畸形\",\"group_name\":\"颜面\"},{\"id\":113,\"label_id\":\"M001-U019\",\"label_name\":\"鱼嘴\",\"group_name\":\"颜面\"},{\"id\":114,\"label_id\":\"M001-U015\",\"label_name\":\"面横裂畸形\",\"group_name\":\"颜面\"},{\"id\":115,\"label_id\":\"M001-U014\",\"label_name\":\"正中面裂\",\"group_name\":\"颜面\"},{\"id\":116,\"label_id\":\"M001-U018\",\"label_name\":\"腭裂\",\"group_name\":\"颜面\"},{\"id\":117,\"label_id\":\"M001-U016\",\"label_name\":\"假唇裂\",\"group_name\":\"颜面\"},{\"id\":118,\"label_id\":\"M001-U017\",\"label_name\":\"唇裂\",\"group_name\":\"颜面\"},{\"id\":119,\"label_id\":\"M001-U009\",\"label_name\":\"鼻发育异常\",\"group_name\":\"颜面\"},{\"id\":120,\"label_id\":\"M001-U012\",\"label_name\":\"鼻尖分裂\",\"group_name\":\"颜面\"},{\"id\":121,\"label_id\":\"M001-U010\",\"label_name\":\"鼻梁扁平\",\"group_name\":\"颜面\"},{\"id\":122,\"label_id\":\"M001-U011\",\"label_name\":\"鼻骨发育不良\",\"group_name\":\"颜面\"},{\"id\":123,\"label_id\":\"M001-U013\",\"label_name\":\"喙鼻\",\"group_name\":\"颜面\"},{\"id\":124,\"label_id\":\"M001-U024\",\"label_name\":\"耳前皮赘\",\"group_name\":\"颜面\"},{\"id\":125,\"label_id\":\"M001-U025\",\"label_name\":\"耳低置\",\"group_name\":\"颜面\"},{\"id\":126,\"label_id\":\"M001-U026\",\"label_name\":\"耳畸形\",\"group_name\":\"颜面\"},{\"id\":127,\"label_id\":\"M001-U023\",\"label_name\":\"小耳畸形\",\"group_name\":\"颜面\"},{\"id\":128,\"label_id\":\"M001-U008\",\"label_name\":\"白内障\",\"group_name\":\"颜面\"},{\"id\":129,\"label_id\":\"M001-U005\",\"label_name\":\"眼距窄\",\"group_name\":\"颜面\"},{\"id\":130,\"label_id\":\"M001-U004\",\"label_name\":\"眼距宽\",\"group_name\":\"颜面\"},{\"id\":131,\"label_id\":\"M001-U007\",\"label_name\":\"眼球突出\",\"group_name\":\"颜面\"},{\"id\":132,\"label_id\":\"M001-U002\",\"label_name\":\"眼眶囊肿\",\"group_name\":\"颜面\"},{\"id\":133,\"label_id\":\"M001-U001\",\"label_name\":\"眼发育异常\",\"group_name\":\"颜面\"},{\"id\":134,\"label_id\":\"M001-U006\",\"label_name\":\"独眼畸形\",\"group_name\":\"颜面\"},{\"id\":135,\"label_id\":\"M001-U003\",\"label_name\":\"小眼畸形\",\"group_name\":\"颜面\"},{\"id\":136,\"label_id\":\"M002-U045\",\"label_name\":\"皮质异位\",\"group_name\":\"中枢神经系统\"},{\"id\":137,\"label_id\":\"M002-U069\",\"label_name\":\"脊柱裂\",\"group_name\":\"中枢神经系统\"},{\"id\":138,\"label_id\":\"M002-U067\",\"label_name\":\"神经管畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":139,\"label_id\":\"M002-U047\",\"label_name\":\"巨头畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":140,\"label_id\":\"M002-U046\",\"label_name\":\"小头畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":141,\"label_id\":\"M002-U065\",\"label_name\":\"脂肪瘤\",\"group_name\":\"中枢神经系统\"},{\"id\":142,\"label_id\":\"M002-U066\",\"label_name\":\"颅内出血\",\"group_name\":\"中枢神经系统\"},{\"id\":143,\"label_id\":\"M002-U044\",\"label_name\":\"颅内钙化灶\",\"group_name\":\"中枢神经系统\"},{\"id\":144,\"label_id\":\"M002-U064\",\"label_name\":\"小脑囊肿\",\"group_name\":\"中枢神经系统\"},{\"id\":145,\"label_id\":\"M002-U062\",\"label_name\":\"脑中线囊肿\",\"group_name\":\"中枢神经系统\"},{\"id\":146,\"label_id\":\"M002-U061\",\"label_name\":\"脉络膜丛小囊\",\"group_name\":\"中枢神经系统\"},{\"id\":147,\"label_id\":\"M002-U060\",\"label_name\":\"室管膜下囊肿\",\"group_name\":\"中枢神经系统\"},{\"id\":148,\"label_id\":\"M002-U063\",\"label_name\":\"蛛网膜囊肿\",\"group_name\":\"中枢神经系统\"},{\"id\":149,\"label_id\":\"M002-U043\",\"label_name\":\"外侧裂形态异常\",\"group_name\":\"中枢神经系统\"},{\"id\":150,\"label_id\":\"M002-U042\",\"label_name\":\"无脑回畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":151,\"label_id\":\"M002-U041\",\"label_name\":\"无脑畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":152,\"label_id\":\"M002-U040\",\"label_name\":\"露脑畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":153,\"label_id\":\"M002-U037\",\"label_name\":\"半叶全前脑\",\"group_name\":\"中枢神经系统\"},{\"id\":154,\"label_id\":\"M002-U036\",\"label_name\":\"全前脑\",\"group_name\":\"中枢神经系统\"},{\"id\":155,\"label_id\":\"M002-U057\",\"label_name\":\"后颅窝异常\",\"group_name\":\"中枢神经系统\"},{\"id\":156,\"label_id\":\"M002-U056\",\"label_name\":\"后颅窝池增宽\",\"group_name\":\"中枢神经系统\"},{\"id\":157,\"label_id\":\"M002-U055\",\"label_name\":\"Dandy-Walker畸形\",\"group_name\":\"中枢神经系统\"},{\"id\":158,\"label_id\":\"M002-U059\",\"label_name\":\"香蕉小脑\",\"group_name\":\"中枢神经系统\"},{\"id\":159,\"label_id\":\"M002-U058\",\"label_name\":\"小脑发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":160,\"label_id\":\"M002-U054\",\"label_name\":\"小脑蚓部发育不全\",\"group_name\":\"中枢神经系统\"},{\"id\":161,\"label_id\":\"M002-U053\",\"label_name\":\"小脑蚓部发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":162,\"label_id\":\"M002-U048\",\"label_name\":\"胼胝体缺失\",\"group_name\":\"中枢神经系统\"},{\"id\":163,\"label_id\":\"M002-U049\",\"label_name\":\"胼胝体发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":164,\"label_id\":\"M002-U050\",\"label_name\":\"透明隔发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":165,\"label_id\":\"M002-U051\",\"label_name\":\"视隔发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":166,\"label_id\":\"M002-U039\",\"label_name\":\"脑发育不良\",\"group_name\":\"中枢神经系统\"},{\"id\":167,\"label_id\":\"M002-U052\",\"label_name\":\"丘脑融合\",\"group_name\":\"中枢神经系统\"},{\"id\":168,\"label_id\":\"M002-U034\",\"label_name\":\"脑室扩张\",\"group_name\":\"中枢神经系统\"},{\"id\":169,\"label_id\":\"M002-U038\",\"label_name\":\"脑穿通\",\"group_name\":\"中枢神经系统\"},{\"id\":170,\"label_id\":\"M002-U035\",\"label_name\":\"脑积水\",\"group_name\":\"中枢神经系统\"},{\"id\":171,\"label_id\":\"M002-U068\",\"label_name\":\"脑膨出\",\"group_name\":\"中枢神经系统\"},{\"id\":172,\"label_id\":\"M009-U188\",\"label_name\":\"脐动脉血流阻力增高\",\"group_name\":\"附属结构\"},{\"id\":173,\"label_id\":\"M009-U190\",\"label_name\":\"脐血管增宽\",\"group_name\":\"附属结构\"},{\"id\":174,\"label_id\":\"M009-U189\",\"label_name\":\"脐带水肿\",\"group_name\":\"附属结构\"},{\"id\":175,\"label_id\":\"M009-U197\",\"label_name\":\"羊膜腔内漂浮物\",\"group_name\":\"附属结构\"},{\"id\":176,\"label_id\":\"M009-U199\",\"label_name\":\"羊水过少\",\"group_name\":\"附属结构\"},{\"id\":177,\"label_id\":\"M009-U198\",\"label_name\":\"羊水过多\",\"group_name\":\"附属结构\"},{\"id\":178,\"label_id\":\"M009-U191\",\"label_name\":\"葡萄胎样胎盘\",\"group_name\":\"附属结构\"},{\"id\":179,\"label_id\":\"M009-U192\",\"label_name\":\"胎盘较小\",\"group_name\":\"附属结构\"},{\"id\":180,\"label_id\":\"M009-U193\",\"label_name\":\"胎盘较大\",\"group_name\":\"附属结构\"},{\"id\":181,\"label_id\":\"M009-U195\",\"label_name\":\"胎盘水肿\",\"group_name\":\"附属结构\"},{\"id\":182,\"label_id\":\"M009-U194\",\"label_name\":\"胎盘过早钙化\",\"group_name\":\"附属结构\"},{\"id\":183,\"label_id\":\"M009-U196\",\"label_name\":\"胎盘增厚\",\"group_name\":\"附属结构\"},{\"id\":184,\"label_id\":\"M005-U138\",\"label_name\":\"骶尾部畸胎瘤\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":185,\"label_id\":\"M005-U137\",\"label_name\":\"肛门闭锁\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":186,\"label_id\":\"M005-U133\",\"label_name\":\"阴茎短小\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":187,\"label_id\":\"M005-U134\",\"label_name\":\"生殖器发育不良\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":188,\"label_id\":\"M005-U135\",\"label_name\":\"尿道异常\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":189,\"label_id\":\"M005-U136\",\"label_name\":\"尿道下裂\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":190,\"label_id\":\"M005-U131\",\"label_name\":\"膀胱未显示\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":191,\"label_id\":\"M005-U130\",\"label_name\":\"膀胱扩张\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":192,\"label_id\":\"M005-U132\",\"label_name\":\"膀胱外翻\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":193,\"label_id\":\"M005-U129\",\"label_name\":\"输尿管异常\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":194,\"label_id\":\"M005-U127\",\"label_name\":\"肾脏回声增强\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":195,\"label_id\":\"M005-U125\",\"label_name\":\"肾脏肿瘤\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":196,\"label_id\":\"M005-U126\",\"label_name\":\"肾脏畸形\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":197,\"label_id\":\"M005-U128\",\"label_name\":\"肾囊性病变\",\"group_name\":\"泌尿生殖会阴\"},{\"id\":198,\"label_id\":\"M005-U124\",\"label_name\":\"肾盂扩张\",\"group_name\":\"泌尿生殖会阴\"}]"

	type SysLabel struct {
		LabelName string `json:"label_name"`
		GroupName string `json:"group_name"`
	}

	var label_list []*SysLabel
	json.Unmarshal([]byte(json_label), &label_list)

	group_sort_list := []string{
		"中枢神经系统",
		"心血管系统",
		"泌尿生殖会阴",
		"生长发育",
		"胸腹部",
		"皮肤及软组织",
		"附属结构",
		"颜面",
		"骨骼系统",
	}

	sys_label_list := g.Array{}
	for k, v := range group_sort_list {
		for _, v2 := range label_list {
			if v == v2.GroupName {
				sys_label_list = append(sys_label_list, g.Map{
					"type":       1,
					"user_id":    0,
					"sort":       k + 100,
					"group_name": v2.GroupName,
					"name":       v2.LabelName,
				})
			}
		}
	}

	case_label.M.Where("type", 1).Delete()
	case_label.M.Insert(sys_label_list)
}

//更新数据
func TestAll(t *testing.T) {
	g.Dump(gtime.Datetime())
}

func TestAll2(t *testing.T) {
	s := g.Map{"2": 2, "1": 1, "3": 3, "4": 4}

	for _, v := range s {
		g.Dump(v)
	}
}

func TestGetUuid(t *testing.T) {
	uuid, err := device.GetDeviceUuid()
	fmt.Println(uuid, err)
}

func TestCurl(t *testing.T) {
	list, err := curl.Post("http://127.0.0.1:9005/api/hospital/list", g.Map{
		"page_size": 1000,
		"region_id": []string{"1", "0010"},
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(string(list))
}

func TestVerifySign(t *testing.T) {
	signature := "EgWI3y6DSZs7Asj0dmKVL322gNw7CT6WFeNXo/wjFTLyEWHGmIqnU+ZxHpFF7jpItTSMX66+GWSTxrCSaW1DKhWogJdwhGvHq7QowjM9myPX8YRJbsvlNhpfqJXqiwOih8VHL/yKRXAnbJEMHaTawUMg7d5ULuJg73ZbAVlUjZE="

	public_key := "MIGJAoGBAMbiZRqO6i6B53crApS0CI3SqbNMZwgp7Ka0rIo0DXX/6uqTd4uaHUOaR2Vq5BRTReejuBkVbaNm7j2GGGei7bTEOOo4ckRhQ63a8JMmnaxaFaIB6/Lcr0HIPek2t+vJutEzPbBlrwWS8Xq6wbaaht76zLbyF2hznx6aIJTVRPc3AgMBAAE="

	var pubKey *rsa.PublicKey
	pubKey, err := rsa_crypt.LoadPublicKeyBase64(public_key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	g.Dump(pubKey)

	bool := rsa_crypt.Verify(pubKey, "1", signature)
	fmt.Println("验签结果:", bool)
	//
	//rsa_pair := rsa_crypt.RsaKeyPair{}
	//rsa_crypt.RSAGenKey(1024, &rsa_pair)
	//
	//var priKey *rsa.PrivateKey
	//priKey, err := rsa_crypt.LoadPrivateKeyBase64(rsa_pair.PriKey)
	//if err != nil {
	//	return
	//}
	//
	//var pubKey *rsa.PublicKey
	//pubKey, err = rsa_crypt.LoadPublicKeyBase64(rsa_pair.PubKey)
	//if err != nil {
	//	return
	//}
	//
	//g.Dump("公私钥对：", rsa_pair)
	//
	//sign := rsa_crypt.Sign(priKey, "1")
	//
	//g.Dump("签名：", sign)
	//
	//bool := rsa_crypt.Verify(pubKey, "1", sign)
	//
	//g.Dump("验签结果：", bool)
}

type PlaneData struct {
	PlaneId     string
	PlaneNameCh string
	PlaneNameEn string
	PlaneHaseId int
}

type QcPlaneData struct {
	ExamineplaneId     string `json:"examineplane_id"`
	ExamineplaneNameCh string `json:"examineplane_name_ch"`
	ExamineplaneNameEn string `json:"examineplane_name_en"`
	ExamineplaneHaseId int    `json:"examineplane_hase_id"`
}

type QcPlaneGroupData struct {
	GroupId           string
	GroupNameCh       string
	GroupNameEn       string
	GroupHaseId       string
	GroupContainPlean []int
}

type PartsItem struct {
	PlaneId     string `json:"plane_id"`
	PlaneNameCh string `json:"plane_name_ch"`
	PlaneNameEn string `json:"plane_name_en"`
}

type GroupItem struct {
	GroupHaseId       int    `json:"group_hase_id"`
	GroupId           string `json:"group_id"`
	GroupNameCh       string `json:"group_name_ch"`
	GroupNameEn       string `json:"group_name_en"`
	GroupContainPlean []int  `json:"group_contain_plean"`
}

// 算法切面配置信息提取
func TestPlaneGen(t *testing.T) {
	file, err := os.Open("./data/study.dat")
	if err != nil {
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	var parts_all map[string]map[string]interface{}
	err = json.Unmarshal(buffer, &parts_all)
	if err != nil {
		g.Log().Error(err)
		return
	}

	var plane_list0 []*PlaneData
	if err = gconv.Struct(parts_all["week_type_0"]["array_plane"], &plane_list0); err != nil {
		return
	}

	var plane_list1 []*PlaneData
	if err = gconv.Struct(parts_all["week_type_1"]["array_plane"], &plane_list1); err != nil {
		return
	}

	var plane_list2 []*PlaneData
	if err = gconv.Struct(parts_all["week_type_2"]["array_plane"], &plane_list2); err != nil {
		return
	}

	var part_map = map[string][]*PartsItem{}

	list0 := []*PartsItem{}
	for _, v := range plane_list0 {
		list0 = append(list0, &PartsItem{
			PlaneId:     v.PlaneId,
			PlaneNameCh: v.PlaneNameCh,
			PlaneNameEn: v.PlaneNameEn,
		})
	}
	part_map["0"] = list0

	list1 := []*PartsItem{}
	for _, v := range plane_list1 {
		list1 = append(list1, &PartsItem{
			PlaneId:     v.PlaneId,
			PlaneNameCh: v.PlaneNameCh,
			PlaneNameEn: v.PlaneNameEn,
		})
	}
	part_map["1"] = list1

	list2 := []*PartsItem{}
	for _, v := range plane_list2 {
		list2 = append(list2, &PartsItem{
			PlaneId:     v.PlaneId,
			PlaneNameCh: v.PlaneNameCh,
			PlaneNameEn: v.PlaneNameEn,
		})
	}
	part_map["2"] = list1

	var qc_plane_list0 []*QcPlaneData
	if err = gconv.Struct(parts_all["week_type_0"]["array_examineplane"], &qc_plane_list0); err != nil {
		return
	}
	var qc_plane_list1 []*QcPlaneData
	if err = gconv.Struct(parts_all["week_type_1"]["array_examineplane"], &qc_plane_list1); err != nil {
		return
	}
	var qc_plane_list2 []*QcPlaneData
	if err = gconv.Struct(parts_all["week_type_2"]["array_examineplane"], &qc_plane_list2); err != nil {
		return
	}
	var qc_part_map = map[string][]*PartsItem{}

	qc_list0 := []*PartsItem{}
	for _, v := range qc_plane_list0 {
		qc_list0 = append(qc_list0, &PartsItem{
			PlaneId:     v.ExamineplaneId,
			PlaneNameCh: v.ExamineplaneNameCh,
			PlaneNameEn: v.ExamineplaneNameEn,
		})
	}
	qc_part_map["0"] = qc_list0

	qc_list1 := []*PartsItem{}
	for _, v := range qc_plane_list1 {
		qc_list1 = append(qc_list1, &PartsItem{
			PlaneId:     v.ExamineplaneId,
			PlaneNameCh: v.ExamineplaneNameCh,
			PlaneNameEn: v.ExamineplaneNameEn,
		})
	}
	qc_part_map["1"] = qc_list1

	qc_list2 := []*PartsItem{}
	for _, v := range qc_plane_list2 {
		qc_list2 = append(qc_list2, &PartsItem{
			PlaneId:     v.ExamineplaneId,
			PlaneNameCh: v.ExamineplaneNameCh,
			PlaneNameEn: v.ExamineplaneNameEn,
		})
	}
	qc_part_map["2"] = qc_list1

	var group_options0 []*GroupItem
	if err = gconv.Struct(parts_all["week_type_0"]["array_examinegroup"], &group_options0); err != nil {
		return
	}
	var group_options1 []*GroupItem
	if err = gconv.Struct(parts_all["week_type_1"]["array_examinegroup"], &group_options1); err != nil {
		return
	}
	var group_options2 []*GroupItem
	if err = gconv.Struct(parts_all["week_type_2"]["array_examinegroup"], &group_options2); err != nil {
		return
	}

	var part_group_map = map[string][]*GroupItem{}
	part_group_map["0"] = group_options0
	part_group_map["1"] = group_options1
	part_group_map["2"] = group_options2

	g.Dump("--------------")
	g.Dump(g.Map{"plane": part_map, "qc_plane": qc_part_map, "qc_plane_group": part_group_map})

	jsonFile, _ := os.Create("./data/plane.json")
	defer jsonFile.Close()

	b, err := json.Marshal(g.Map{"plane": part_map, "qc_plane": qc_part_map, "qc_plane_group": part_group_map})

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")

	jsonFile.Write(out.Bytes())
}
