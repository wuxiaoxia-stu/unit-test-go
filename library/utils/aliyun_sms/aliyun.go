package aliyun_sms

import (
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/gogf/gf/frame/g"
)

var Sign = g.MapStrStr{
	"SMS_251016135": "广州爱孕记", //【广州爱孕记】您正在申请手机注册，验证码为：602541，5分钟内有效！
}

func SMS(phone_number, template_code string, code_map g.MapStrAny) error {
	accessKeyId := g.Cfg().GetString("AliSMS.accessKeyId")
	accessKeySecret := g.Cfg().GetString("AliSMS.accessKeySecret")
	endpoint := g.Cfg().GetString("AliSMS.endpoint")
	c := &client.Config{AccessKeyId: &accessKeyId, AccessKeySecret: &accessKeySecret, Endpoint: &endpoint}

	newClient, err := dysmsapi20170525.NewClient(c)
	if err != nil {
		g.Log().Error(err)
		return err
	}

	code_byte, err := json.Marshal(code_map)
	if err != nil {
		g.Log().Error(err)
		return err
	}
	code := string(code_byte)
	sign := Sign[template_code]
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phone_number,
		TemplateCode:  &template_code,
		SignName:      &sign,
		TemplateParam: &code,
	}
	sms, err := newClient.SendSms(request)
	if err != nil {
		g.Log().Error(err)
		return err
	}

	if *sms.Body.Code != "OK" {
		g.Log().Error(sms)
		return fmt.Errorf("短信发送失败：%s", *sms.Body.Message)
	}

	return nil
}
