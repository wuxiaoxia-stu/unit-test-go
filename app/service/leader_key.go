package service

import (
	"aiyun_local_srv/app/model/leader_key"
	"aiyun_local_srv/app/model/user"
	"github.com/gogf/gf/frame/g"
)

var LeaderKeyService = new(leaderKeyService)

type leaderKeyService struct{}

// 检查绑定记录
func (s *leaderKeyService) Info(where interface{}) (res *leader_key.Entity, err error) {
	err = leader_key.M.Where(where).Limit(1).Scan(&res)
	return
}

//创建绑定信息
func (s *leaderKeyService) Create(leader_serial_number string, user_data []*user.UserData) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//删除旧的绑定记录
	if _, err = tx.Model(leader_key.Table).Where(g.Map{"status": 1}).Data(g.Map{"status": 0}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	//添加新记录
	if _, err = tx.Model(leader_key.Table).Data(g.Map{
		"serial_number": leader_serial_number,
		"status":        1,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	//删除旧的用户数据
	if _, err = tx.Model(user.Table).Delete(); err != nil {
		tx.Rollback()
		return err
	}

	//添加用户信息
	for _, v := range user_data {
		if _, err = tx.Model(user.Table).Data(g.Map{
			"id":        v.UserId,
			"username":  v.Username,
			"role_type": v.RoleType,
			"password":  v.Password,
			"status":    1,
		}).Insert(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
