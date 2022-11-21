package service

import (
	"aiyun_local_srv/app/model/user"
	"github.com/gogf/gf/frame/g"
)

var UserService = new(userService)

type userService struct{}

// 检查绑定记录
func (s *userService) Info(where interface{}) (res *user.Entity, err error) {
	err = user.M.Where(where).Limit(1).Scan(&res)
	return
}

func (s *userService) List(where interface{}) (res []*user.Entity, err error) {
	err = user.M.Where(where).Scan(&res)
	return
}

func (s *userService) Delete(id int) (row int64, err error) {
	ret, err := user.M.Where("id", id).Delete()
	if err != nil {
		return 0, err
	}

	return ret.RowsAffected()
}

//添加用户信息
func (s *userService) Save(user_list []*user.UserData) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if _, err := tx.Model(user.Table).Delete(g.Map{"status": 1}); err != nil {
		tx.Rollback()
		return err
	}

	//添加用户信息
	list := g.Array{}
	deleta_ids := []int{}
	for _, v := range user_list {
		d := g.Map{
			"id":        v.UserId,
			"username":  v.Username,
			"role_type": v.RoleType,
			"password":  v.Password,
			"status":    1,
		}

		if v.DeleteAt > 0 {
			//d["delete_at"] = gtime.Datetime()
			deleta_ids = append(deleta_ids, v.UserId)
		}
		list = append(list, d)
	}

	if _, err = tx.Model(user.Table).Insert(list); err != nil {
		tx.Rollback()
		return err
	}

	if len(deleta_ids) > 0 {
		if _, err = tx.Model(user.Table).Delete("id", deleta_ids); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
