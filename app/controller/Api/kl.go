package Api

import (
	"aiyun_local_srv/app/model/kl_feature"
	"aiyun_local_srv/app/model/kl_syndrome"
	"aiyun_local_srv/app/model/kl_syndrome_feature"
	"aiyun_local_srv/app/service"
	"aiyun_local_srv/library/response"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"strconv"
)

var Kl = klApi{}

type klApi struct{}

// 特征列表
func (*klApi) FeatureList(r *ghttp.Request) {
	list, err := service.KlFeatureService.Tree(g.Map{"status": 1, "invisible": 0})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

// 综合征列表
func (*klApi) SyndromeList(r *ghttp.Request) {
	list, err := service.KlSyndromeService.Tree(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

// 树结构第三层数据uuid
// 根据特征搜索综合征
func (*klApi) SyndromeFeature(r *ghttp.Request) {
	var req *kl_syndrome.FeatureSyndromeReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	var id_arr []int
	for _, v := range req.FeatureUuidArray {
		int, err := strconv.Atoi(v)
		if err != nil {
			response.ErrorSys(r, err)
		}
		id_arr = append(id_arr, int)
	}

	if id_arr == nil {
		response.Success(r)
	}

	var feature_list []*kl_feature.Entity
	if err := kl_feature.M.Fields("id").WhereIn("pid", id_arr).Order("pid,id").Scan(&feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	feature_ids := []int{}
	for _, v := range feature_list {
		feature_ids = append(feature_ids, v.Id)
	}

	//SELECT syndrome_id,count(*) FROM "lpm_kl_syndrome_feature" WHERE feature_id IN (267,646)  GROUP BY "syndrome_id" ORDER BY "count" DESC
	var syndrome_count_list []*kl_syndrome_feature.SyndromeCount
	if err := kl_syndrome_feature.M.
		Fields("syndrome_id,count(*) as count").
		WhereIn("feature_id", feature_ids).
		Group("syndrome_id").Order("count DESC").
		Scan(&syndrome_count_list); err != nil {
		response.ErrorDb(r, err)
	}

	syndrome_ids := []int{}
	for _, v := range syndrome_count_list {
		syndrome_ids = append(syndrome_ids, v.SyndromeId)
	}

	var syndrome_feature_list []*kl_syndrome_feature.Entity
	if err := kl_syndrome_feature.M_alias.
		Fields("ksf.*,kf.name as feature_name,kf.name_en as feature_name_en").
		LeftJoin("kl_feature kf", "kf.id = ksf.feature_id").
		WhereIn("ksf.syndrome_id", syndrome_ids).
		Scan(&syndrome_feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	var syndrome_list []*kl_syndrome.Entity
	if err := kl_syndrome.M.
		Fields("id,name").
		WhereIn("id", syndrome_ids).
		Scan(&syndrome_list); err != nil {
		response.ErrorDb(r, err)
	}

	var list []*kl_syndrome_feature.FeatureSyndromeReq
	for _, v := range syndrome_ids {
		for _, v2 := range syndrome_list {
			if v == v2.Id {
				count := 0
				for _, v3 := range syndrome_count_list {
					if v2.Id == v3.SyndromeId {
						count = v3.Count
					}
				}
				list = append(list, &kl_syndrome_feature.FeatureSyndromeReq{
					SyndromeUuid:    strconv.Itoa(v2.Id),
					SyndromeSerial:  strconv.Itoa(v2.Id),
					SyndromeName:    v2.Name,
					SyndromeNameEn:  v2.NameEn,
					HitGroupCount:   count,
					SyndromeQuality: 1,
				})

			}
		}
	}

	type_options := map[string]string{
		"1": "SY01",
		"2": "SY02",
	}

	for _, v := range list {
		for _, v2 := range syndrome_feature_list {
			if strconv.Itoa(v2.SyndromeId) == v.SyndromeUuid {
				v.Features = append(v.Features, &kl_syndrome_feature.Feature{
					PartUuid:           strconv.Itoa(v2.FeatureRootId),
					PartSerial:         strconv.Itoa(v2.FeatureRootId),
					PartName:           "",
					PartNameEn:         "",
					FeatureUuid:        strconv.Itoa(v2.FeatureId),
					FeatureSerial:      strconv.Itoa(v2.FeatureId),
					FeatureName:        v2.FeatureName,
					FeatureNameEn:      v2.FeatureNameEn,
					SyndromeFeatureOpt: type_options[v2.Type],
				})
				if v2.Type == "1" {
					v.Type1Count++
				}
			}
		}
	}

	// 特征病变总数排序
	//sort.Slice(list, func(i, j int) bool {
	//	return list[i].Type1Count > list[j].Type1Count
	//})

	response.Success(r, list)
}

//特征详情
func (*klApi) FeatureDetail(r *ghttp.Request) {
	feature_uuid := r.GetQueryInt("feature_uuid")

	info, err := service.KlFeatureService.Detail(feature_uuid)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Success(r)
	}

	atlas := g.Array{}
	for _, v := range info.Atlas {
		atlas = append(atlas, g.Map{
			"feature_uuid":       strconv.Itoa(info.Id),
			"legend_file_path":   v.Url,
			"legend_origin_name": v.Name,
			"legend_type":        fmt.Sprintf("FL0%d", v.Type),
			"legend_type_name":   "典型超声图",
			"legend_url":         g.Cfg().GetString("server.Domain") + v.Url,
		})
	}

	//path,err := service.KlFeatureService.FullPath(feature_uuid, "")
	//id_arr := strings.Split(path, ",")
	//part_uuid := ""
	//if len(id_arr) > 0 {
	//	part_uuid = id_arr[len(id_arr) - 1]
	//}

	response.Success(r, g.Map{
		"feature_serial":   strconv.Itoa(info.Id),
		"feature_uuid":     strconv.Itoa(info.Id),
		"feature_name":     info.Name,
		"feature_name_en":  info.NameEn,
		"feature_other":    info.Other,
		"feature_other_en": info.OtherEn,
		//"part_uuid":           part_uuid,
		"feature_consult":     info.Consult,
		"feature_consult_en":  info.ConsultEn,
		"feature_define":      info.Define,
		"feature_define_en":   info.DefineEn,
		"feature_diagnose":    info.Diagnose,
		"feature_diagnose_en": info.DiagnoseEn,
		"feature_legends":     atlas,
	})
}

//综合征详情
func (*klApi) SyndromeDetail(r *ghttp.Request) {
	syndrome_uuid := r.GetQueryInt("syndrome_uuid")

	info, err := service.KlSyndromeService.Detail(syndrome_uuid)
	if err != nil {
		response.ErrorDb(r, err)
	}

	part_list, err := service.KlFeatureService.PartList()
	if err != nil {
		response.ErrorDb(r, err)
	}

	type_options := map[string]string{
		"1": "SY01",
		"2": "SY02",
	}

	morphology_list := g.Array{}
	for _, v := range part_list {
		features := g.Array{}
		others := g.Array{}
		exist := false
		for _, v2 := range info.Features {
			if v2.RootId == v.Id {
				exist = true
				if v2.Type == "1" {
					features = append(features, g.Map{
						"feature_name":         v2.Name,
						"feature_name_en":      v2.NameEn,
						"feature_serial":       strconv.Itoa(v2.Id),
						"feature_uuid":         strconv.Itoa(v2.Id),
						"part_serial":          strconv.Itoa(v.Id),
						"part_uuid":            strconv.Itoa(v.Id),
						"syndrome_feature_opt": type_options[v2.Type],
					})
				} else {
					others = append(others, g.Map{
						"feature_name":         v2.Name,
						"feature_name_en":      v2.NameEn,
						"feature_serial":       strconv.Itoa(v2.Id),
						"feature_uuid":         strconv.Itoa(v2.Id),
						"part_serial":          strconv.Itoa(v.Id),
						"part_uuid":            strconv.Itoa(v.Id),
						"syndrome_feature_opt": type_options[v2.Type],
					})
				}
			}
		}

		if exist {
			morphology_list = append(morphology_list, g.Map{
				"part_serial":  strconv.Itoa(v.RootId),
				"part_uuid":    strconv.Itoa(v.RootId),
				"part_name":    v.Name,
				"part_name_en": v.NameEn,
				"created_time": 0,
				"features":     features,
				"others":       others,
			})
		}
	}

	response.Success(r, g.Map{
		"syndrome_serial":      strconv.Itoa(info.Id),
		"syndrome_uuid":        strconv.Itoa(info.Id),
		"syndrome_name":        info.Name,
		"syndrome_name_en":     info.NameEn,
		"gene_location":        info.GeneLocation,
		"gene_location_en":     info.GeneLocationEn,
		"genetics_desc":        info.GeneticsDesc,
		"genetics_desc_en":     info.GeneticsDescEn,
		"syndrome_consult":     info.Consult,
		"syndrome_consult_en":  info.ConsultEn,
		"syndrome_diagnose":    info.Diagnose,
		"syndrome_diagnose_en": info.DiagnoseEn,
		"morphology_list":      morphology_list,
	})
}

//综合征对比
func (*klApi) SyndromeMorphologies(r *ghttp.Request) {
	var req *kl_syndrome.SyndromeUuidArrayReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	part_list, err := service.KlFeatureService.PartList()
	if err != nil {
		response.ErrorDb(r, err)
	}
	type_options := map[string]string{
		"1": "SY01",
		"2": "SY02",
	}

	list := g.Array{}
	for _, v := range req.SyndromeUuidArray {
		int, err := strconv.Atoi(v)
		if err != nil {
			response.ErrorSys(r, err)
		}
		info, err := service.KlSyndromeService.Detail(int)
		if err != nil {
			response.ErrorDb(r, err)
		}

		morphology_list := g.Array{}
		if info == nil {
			break
		}

		for _, v := range part_list {
			features := g.Array{}
			others := g.Array{}
			exist := false

			for _, v2 := range info.Features {
				if v2.RootId == v.Id {
					exist = true
					if v2.Type == "1" {
						features = append(features, g.Map{
							"feature_name":         v2.Name,
							"feature_name_en":      v2.NameEn,
							"feature_serial":       strconv.Itoa(v2.Id),
							"feature_uuid":         strconv.Itoa(v2.Id),
							"part_serial":          strconv.Itoa(v.Id),
							"part_uuid":            strconv.Itoa(v.Id),
							"syndrome_feature_opt": type_options[v2.Type],
						})
					} else {
						others = append(others, g.Map{
							"feature_name":         v2.Name,
							"feature_name_en":      v2.NameEn,
							"feature_serial":       strconv.Itoa(v2.Id),
							"feature_uuid":         strconv.Itoa(v2.Id),
							"part_serial":          strconv.Itoa(v.Id),
							"part_uuid":            strconv.Itoa(v.Id),
							"syndrome_feature_opt": type_options[v2.Type],
						})
					}
				}
			}

			if exist {
				morphology_list = append(morphology_list, g.Map{
					"part_serial":  strconv.Itoa(v.RootId),
					"part_uuid":    strconv.Itoa(v.RootId),
					"part_name":    v.Name,
					"part_name_en": v.NameEn,
					"created_time": 0,
					"features":     features,
					"others":       others,
				})
			}
		}

		list = append(list, g.Map{
			"syndrome_name":    info.Name,
			"syndrome_name_en": info.NameEn,
			"syndrome_serial":  strconv.Itoa(info.Id),
			"syndrome_uuid":    strconv.Itoa(info.Id),
			"morphology_list":  morphology_list,
		})
	}

	response.Success(r, list)
}
