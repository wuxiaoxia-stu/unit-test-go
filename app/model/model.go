package model

import (
	"aiyun_local_srv/app/model/case_info"
	"aiyun_local_srv/app/model/case_label"
	"aiyun_local_srv/app/model/case_label_relate"
	"aiyun_local_srv/app/model/case_measured"
	"aiyun_local_srv/app/model/case_patient"
	"aiyun_local_srv/app/model/case_plane"
	"aiyun_local_srv/app/model/case_plane_structure"
	qc_score_log "aiyun_local_srv/app/model/case_score_log"
	"aiyun_local_srv/app/model/case_shot"
	"aiyun_local_srv/app/model/case_shot_plane"
	"aiyun_local_srv/app/model/case_shot_plane_structure"
	"aiyun_local_srv/app/model/check"
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_feature_atlas"
	"aiyun_local_srv/app/model/kl_syndrome"
	"aiyun_local_srv/app/model/kl_syndrome_feature"
	"aiyun_local_srv/app/model/leader_key"
	"aiyun_local_srv/app/model/licence"
	"aiyun_local_srv/app/model/meeting_device"
	"aiyun_local_srv/app/model/pre_pair"
	"aiyun_local_srv/app/model/qc_upload"
	"aiyun_local_srv/app/model/qc_user"
	"aiyun_local_srv/app/model/sys_admin"
	"aiyun_local_srv/app/model/sys_role"
	"aiyun_local_srv/app/model/user"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"os"
	"reflect"
	"strings"
)

type Table struct {
	TableName  string `orm:"table_name"`
	ColumnName string `orm:"column_name"`
	UdtName    string `orm:"udt_name"`
}

type Unique struct {
	Tablename string `orm:"tablename"`
	Indexname string `orm:"indexname"`
}

//分页查询
type PageReqParams struct {
	KeyWord   string   `p:"keyword"`
	RegionId  []string `p:"region_id"`
	Status    int      `p:"status" default:"-1"`
	Page      int      `p:"page" default:"1"`
	PageSize  int      `p:"page_size" default:"10"`
	Order     string   `p:"order" default:"id"`
	Sort      string   `p:"sort" default:"DESC"`
	StartTime string   `p:"start_time"`
	EndTime   string   `p:"end_time"`
}

//状态设置
type SetStatusParams struct {
	Id     int `json:"id" p:"id" v:"required#参数错误"`
	Status int `json:"status" p:"status" v:"required|in:0,1#参数错误|参数错误"`
}

//批量操作提交信息绑定
type Ids struct {
	Ids []int `json:"ids"`
}

func DbInit() {
	//ImportData()
	if err := dbAutoMigrate(
		&sys_admin.Entity{},
		&sys_role.Entity{},
		&licence.Entity{},
		&pre_pair.Entity{},
		&leader_key.Entity{},
		&user.Entity{},
		&case_info.Entity{},
		&check.Entity{},
		&case_label.Entity{},
		&case_label_relate.Entity{},
		&case_patient.Entity{},
		&case_measured.Entity{},
		&case_plane.Entity{},
		&case_plane_structure.Entity{},
		&case_shot.Entity{},
		&case_shot_plane.Entity{},
		&case_shot_plane_structure.Entity{},
		&kl_feature.Entity{},
		&kl_feature_atlas.Entity{},
		&kl_syndrome.Entity{},
		&kl_syndrome_feature.Entity{},
		&qc_upload.Entity{},
		&qc_user.Entity{},
		&qc_score_log.Entity{},
		&meeting_device.Entity{},
	); err != nil {
		panic("数据库初始化失败")
	}

	if err := initAdminData(); err != nil {
		panic("Admin表数据初始化失败")
	}
}

//导入数据
func ImportData() (err error) {
	file, err := os.Open("./data/public_v1.sql")
	if err != nil {
		g.Log().Error(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		g.Log().Error(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	sql := string(buffer)

	_, err = g.DB().Exec(sql)
	if err != nil {
		g.Log().Error("数据导入失败", err)
		return
	} else {
		g.Log().Info("数据导入成功")
	}

	return
}

//初始化数据
func initAdminData() (err error) {
	var res *sys_admin.Entity
	err = sys_admin.M.Where("username", "admin").Scan(&res)
	if err != nil {
		g.Log().Error(err)
		return
	}

	if res == nil {
		_, err = sys_admin.M.Data(g.Map{
			"role_id":  1,
			"username": "admin",
			"password": "315193380d2fedffa677b7fe236fadb5",
			"salt":     "sinb",
			"email":    "778774780@qq.com",
			"status":   1,
		}).Insert()
		if err != nil {
			g.Log().Error("初始化admin表数据失败", err)
			return
		} else {
			g.Log().Info("初始化admin表数据成功")
		}
	}
	return
}

//自动更新表结构，添加新字段
func dbAutoMigrate(models ...interface{}) (err error) {

	if len(models) <= 0 {
		return
	}

	for _, model := range models {
		t := reflect.TypeOf(model)

		if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
			err = gerror.New("模型参数应为结构体指针")
			g.Log().Error(err)
			return err
		}

		model_name := t.String()
		model_arr := strings.Split(strings.TrimLeft(model_name, "*"), ".")
		if len(model_arr) != 2 {
			err = gerror.New("模型参数错误")
			g.Log().Error(err)
			return err
		}

		//判断表是否存在
		//select * from pg_tables where tablename = 'lpm_sys_role'
		//select * from information_schema.TABLES where TABLE_NAME = 'lpm_sys_role';

		table_prefix := g.Cfg().GetString("database.prefix")
		full_table_name := table_prefix + model_arr[0]

		table := Table{}
		err = g.DB().GetScan(&table, "select * from information_schema.TABLES where TABLE_NAME = ?", full_table_name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
			} else {
				g.Log().Error(err)
				return
			}
		}

		if table.TableName == "" {
			//CREATE TABLE "public"."lpm_sys_role" ();
			_, err = g.DB().Exec("CREATE TABLE \"" + full_table_name + "\" ()")
			if err != nil {
				g.Log().Error(err)
				return
			} else {
				g.Log().Infof("%s表创建成功", full_table_name)
			}
		}

		v := reflect.ValueOf(model).Elem()

		//查询表字段
		//select * from information_schema.COLUMNS where table_name = 'lpm_sys_admin'
		columns := []Table{}
		err = g.DB().GetScan(&columns, "select * from information_schema.COLUMNS where table_name = ?", full_table_name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
			} else {
				g.Log().Error(err)
				return
			}
		}

		add_colums_sql := []string{}       //字段信息
		alter_colums_comment := []string{} //备注信息
		table_comment := ""
		unique := []string{}
		for i := 0; i < v.NumField(); i++ {
			tagOrmInfo := v.Type().Field(i).Tag.Get("orm")
			if tagOrmInfo == "" {
				break
			}

			tag_orms := strings.Split(tagOrmInfo, ",")
			if len(tag_orms) <= 0 {
				err = gerror.New("")
				g.Log().Error(err)
				return err
			}

			column := tag_orms[0]
			column_type := v.Field(i).Type().String()
			size := "4"
			default_val := ""
			is_primary := false
			not_null := ""
			primary_start := "1"
			for _, v2 := range tag_orms {
				if strings.Contains(v2, "primary") {
					is_primary = true
				}

				if strings.Contains(v2, "start") {
					primary_start = strings.TrimLeft(v2, "start:")
				}

				if strings.Contains(v2, "table_comment:") {
					table_comment = strings.TrimLeft(v2, "table_comment:")
				}

				if strings.Contains(v2, "unique") {
					unique = append(unique, column)
				}

				if strings.Contains(v2, "size:") {
					size = strings.TrimLeft(v2, "size:")
				}

				if strings.Contains(strings.ToUpper(v2), "NOT NULL") {
					not_null = "NOT NULL"
				}

				if strings.Contains(v2, "default:") {
					default_val = strings.TrimLeft(v2, "default:")
				}

				if strings.Contains(v2, "comment:") && !strings.Contains(v2, "table_comment:") {
					comment := strings.TrimLeft(v2, "comment:")
					if comment != "" {
						alter_colums_comment = append(alter_colums_comment, fmt.Sprintf("COMMENT ON COLUMN %s.%s IS %s", full_table_name, column, comment))
					}
				}
			}

			column_exist := false
			for _, v2 := range columns {
				if column == v2.ColumnName {
					column_exist = true
					break
				}
			}

			//如果字段你不存在，就创建字段， 预定义字段类型
			// int => int(size)
			// string =>  varchar(size)
			// gtime.Time =>  timestamptz(6)
			if !column_exist {
				t := "int"
				if column_type == "string" {
					t = "varchar"
				} else if strings.Contains(column_type, "gtime.Time") {
					t = "timestamptz"
					size = "6"
				} else if strings.Contains(column_type, "float") {
					t = "float"
				}

				if size == "text" {
					t = "text"
				}

				sql := "ADD COLUMN "
				if t == "int" {
					sql += column + " " + t + size
				} else if t == "text" {
					sql += column + " " + t
				} else if t == "float" {
					sql += column + " numeric(" + strings.Replace(size, ":", ",", 1) + ")"
				} else {
					sql += column + " " + t + "(" + size + ")"
				}

				if is_primary {
					sql += " NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 MINVALUE  1 START " + primary_start + " CACHE 1)"
				}

				if not_null != "" {
					sql += " NOT NULL"
				}

				if default_val != "" {
					sql += " DEFAULT " + default_val
				}

				add_colums_sql = append(add_colums_sql, sql)
			}

		}

		//ALTER TABLE lpm_sys_role
		//	ADD COLUMN id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY ( INCREMENT 1 MINVALUE  1 START 1 CACHE 1),
		//	ADD COLUMN name VARCHAR(50) NOT NULL DEFAULT '',
		//	ADD COLUMN create_at timestamptz(6)
		if len(add_colums_sql) > 0 {
			_, err = g.DB().Exec("ALTER TABLE " + full_table_name + " " + strings.Join(add_colums_sql, ","))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}

		//修改行注释
		//comment on column table_name.column_name is '名称';
		if len(alter_colums_comment) > 0 {
			_, err = g.DB().Exec(strings.Join(alter_colums_comment, ";"))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}
		//修改表注释
		//comment on table table_name is '表名称';
		if table_comment != "" {
			_, err = g.DB().Exec(fmt.Sprintf("COMMENT ON TABLE %s IS %s;", full_table_name, table_comment))
			if err != nil {
				g.Log().Error(err)
				return
			}
		}

		// 添加唯一索引
		//alter table lpm_auth_licence ADD CONSTRAINT unique_author_number UNIQUE(author_number)
		if len(unique) > 0 {
			//查询索引是否存在,不存在则创建
			//select * from pg_indexes where tablename = 'lpm_auth_licence';
			unique_data := []Unique{}
			err = g.DB().GetScan(&unique_data, "select * from pg_indexes where tablename = ?", full_table_name)
			if err != nil {
				g.Log().Error(err)
				return
			}
			for _, v := range unique {
				index_exist := false
				for _, v2 := range unique_data {
					if v2.Indexname == "unique_"+v {
						index_exist = true
					}
				}

				if !index_exist {
					_, err = g.DB().Exec(fmt.Sprintf("alter table %s ADD CONSTRAINT unique_%s UNIQUE(%s);", full_table_name, v, v))
					if err != nil {
						g.Log().Error(err)
						return
					}
				}

			}
		}
	}

	return
}
