package jwt

import (
	"aiyun_local_srv/app/model/qc_user"
	"aiyun_local_srv/app/model/sys_admin"
	"aiyun_local_srv/app/model/user"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
)

//
//  admin 管理员信息信息
//  roleIds 用户所属的角色id
//  isRefresh 是否是刷新token
//  exp 过期时间
//  @return string 返回token
//  @return error
//
func GenerateLoginToken(admin *sys_admin.Entity) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       admin.Id,
		"role_id":  admin.RoleId,
		"username": admin.Username,
		"avatar":   admin.Avatar,
		"status":   admin.Status,
		"rand":     grand.Letters(20),
	})

	tokenString, err := token.SignedString([]byte(g.Cfg().GetString("jwt.sign", "jwt_sign")))
	return tokenString, err
}

func GenerateCustomerLoginToken(customer *user.Entity, author_number, serial_number, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            customer.Id,
		"role_type":     customer.RoleType,
		"username":      customer.Username,
		"ip":            ip,
		"author_number": author_number,
		"serial_number": serial_number,
		"rand":          grand.Letters(20),
	})

	tokenString, err := token.SignedString([]byte(g.Cfg().GetString("jwt.sign", "jwt_sign")))
	return tokenString, err
}

func GenerateQcUserLoginToken(user *qc_user.Entity, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            user.Id,
		"role_type":     user.RoleType,
		"username":      user.Username,
		"realname":      user.RealName,
		"region_id":     user.RegionId,
		"region_name":   user.RegionName,
		"hospital_id":   user.HospitalId,
		"hospital_name": user.HospitalName,
		"ip":            ip,
		"rand":          grand.Letters(20),
	})

	tokenString, err := token.SignedString([]byte(g.Cfg().GetString("qc.sign", "jwt_sign")))
	return tokenString, err
}

// 解析token
// claims["UserId"]这样使用
func ParseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	if tokenString == "" {
		err := gerror.New("token 为空")
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
