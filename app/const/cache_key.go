package _const

import "fmt"

//token缓存
func TOKEN_CACHE_KEY(token string) string {
	return fmt.Sprintf("token:%s", token)
}

// 手机短信验证码缓存
func PHONE_CODE_CACHE_KEY(scene, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", scene, phone)
}

// 图形验证码缓存
func CAPTCHA_CODE_CACHE_KEY(capacha_id string) string {
	return fmt.Sprintf("captcha_code:%s", capacha_id)
}
