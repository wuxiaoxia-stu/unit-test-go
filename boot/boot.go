package boot

import (
	"aiyun_local_srv/app/model"
	"aiyun_local_srv/library/utils"
	_ "aiyun_local_srv/packed"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	//清理日志
	clearLog()

	//创建DB
	createDb()

	//初始化数据库
	model.DbInit()

	// 系统定时器
	Timer()
}

//清理一个月之前的日志
func clearLog() {
	g.Log().Info("开始清理日志！")

	log_path := []string{}
	log_path = append(log_path, g.Cfg().GetString("database.logger.Path"))
	log_path = append(log_path, g.Cfg().GetString("server.LogPath"))
	log_path = append(log_path, g.Cfg().GetString("logger.Path"))

	// 日志存活时间
	log_expire := g.Cfg().GetInt("logger.Expire", 30)

	for _, v := range log_path {
		path_exist, _ := utils.PathExists(v)
		if strings.Contains(v, "log/") && path_exist {
			filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
				if info.ModTime().Before(time.Now().AddDate(0, 0, -1*log_expire)) {
					if err := os.Remove(path); err != nil {
						g.Log().Errorf("日志文件【%s】删除失败", path)
					} else {
						g.Log().Infof("日志文件【%s】已删除", path)
					}
				}
				return nil
			})
		}
	}
	g.Log().Info("日志清理完成！")
}

// 自动创建DB
func createDb() {
	host := g.Cfg().GetString("database.host")
	post := g.Cfg().GetString("database.port")
	user := g.Cfg().GetString("database.user")
	pass := g.Cfg().GetString("database.pass")
	default_user := g.Cfg().GetString("database.default_user")
	default_pass := g.Cfg().GetString("database.default_pass")
	db_name := g.Cfg().GetString("database.name")
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", host, post, default_user, default_pass))
	if err != nil {
		g.Log().Error(err)
	}
	err = db.Ping()
	if err != nil {
		g.Log().Error(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", db_name))
	if err != nil {
		g.Log().Info(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", user, pass))
	if err != nil {
		g.Log().Info(err)
	}

	_, err = db.Exec(fmt.Sprintf("GRANT ALL ON DATABASE %s to %s", db_name, user))
	if err != nil {
		g.Log().Error(err)
	}

	db.Close()
}

func Timer() {
	// 心跳检查
	interval := time.Second
	gtimer.AddSingleton(interval, func() {
		glog.Println("doing")
		time.Sleep(5 * time.Second)
	})

}
