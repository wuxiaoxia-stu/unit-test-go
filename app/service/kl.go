package service

import (
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_feature_atlas"
	"aiyun_local_srv/app/model/kl_syndrome"
	"aiyun_local_srv/app/model/kl_syndrome_feature"
	"aiyun_local_srv/library/utils"
	"aiyun_local_srv/library/utils/encrypt"
	"encoding/json"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"os"
	"path/filepath"
	"strings"
)

var KlService = new(klService)

type klService struct{}

//获取列表数据
func (s *klService) Import(zip_file string) (err error) {
	//解压文件输出路径
	dst := "public/attachment/"
	//删除路径下所有文件
	if err := os.RemoveAll(dst + "feature_ledge"); err != nil {
		return err
	}

	//解压缩文件到对应目录
	if err := utils.UnzipFile(zip_file, dst); err != nil {
		return err
	}

	//读取所有解压文件并解密文件
	file_list_all, err := utils.GetAllFile("public/attachment/feature_ledge")
	if err != nil {
		return
	}

	//去除后缀为data的文件
	file_list := []string{}
	for _, v := range file_list_all {
		if strings.Contains(v, ".data") {
			file_list = append(file_list, v)
		}
	}

	//对目录下.data所有文件解密
	for _, v := range file_list {
		dir, fname := filepath.Split(v)
		fname = strings.TrimRight(fname, ".data")
		des := filepath.Join(dir, fname)
		//文件加密
		err = encrypt.DecryptFile(v, des, "123456")
		if err != nil {
			return
		}
	}

	//导入数据库数据
	err = s.ImportData()

	return
}

func (s *klService) ImportData() error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//导入数据库数据
	if err := s.ImportFeatureData(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.ImportFeatureAtlasData(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.ImportSyndromeData(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.ImportSyndromeFeatureData(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 导入特征数据到数据库
func (s *klService) ImportFeatureData(tx *gdb.TX) error {
	b, err := os.ReadFile("public/attachment/feature_ledge/data/kl_feature.json")
	if err != nil {
		return err
	}

	var kl_feature_list []*kl_feature.Entity
	if err := json.Unmarshal(b, &kl_feature_list); err != nil {
		return err
	}

	if _, err := tx.Model(kl_feature.Table).Delete(g.Map{"id>": 0}); err != nil {
		return err
	}

	_, err = tx.Model(kl_feature.Table).Insert(kl_feature_list)
	return err
}

// 导入特征图例数据到数据库
func (s *klService) ImportFeatureAtlasData(tx *gdb.TX) error {
	b, err := os.ReadFile("public/attachment/feature_ledge/data/kl_feature_atlas.json")
	if err != nil {
		return err
	}

	var kl_feature_atlas_list []*kl_feature_atlas.Entity
	if err := json.Unmarshal(b, &kl_feature_atlas_list); err != nil {
		return err
	}

	if _, err := tx.Model(kl_feature_atlas.Table).Delete(g.Map{"id>": 0}); err != nil {
		return err
	}

	_, err = tx.Model(kl_feature_atlas.Table).Insert(kl_feature_atlas_list)
	return err
}

// 导入综合征数据到数据库
func (s *klService) ImportSyndromeData(tx *gdb.TX) error {
	b, err := os.ReadFile("public/attachment/feature_ledge/data/kl_syndrome.json")
	if err != nil {
		return err
	}

	var kl_syndrome_list []*kl_syndrome.Entity
	if err := json.Unmarshal(b, &kl_syndrome_list); err != nil {
		return err
	}

	if _, err := tx.Model(kl_syndrome.Table).Delete(g.Map{"id>": 0}); err != nil {
		return err
	}

	_, err = tx.Model(kl_syndrome.Table).Insert(kl_syndrome_list)
	return err
}

// 导入综合征特征关系数据到数据库
func (s *klService) ImportSyndromeFeatureData(tx *gdb.TX) error {
	b, err := os.ReadFile("public/attachment/feature_ledge/data/kl_syndrome_feature.json")
	if err != nil {
		return err
	}

	var kl_syndrome_feature_list []*kl_syndrome_feature.Entity
	if err := json.Unmarshal(b, &kl_syndrome_feature_list); err != nil {
		return err
	}

	if _, err := tx.Model(kl_syndrome_feature.Table).Delete(g.Map{"id>": 0}); err != nil {
		return err
	}

	_, err = tx.Model(kl_syndrome_feature.Table).Insert(kl_syndrome_feature_list)
	return err
}
