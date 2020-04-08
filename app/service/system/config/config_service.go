package config

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtime"
	configModel "yj-app/app/model/system/config"
	userService "yj-app/app/service/system/user"
	"yj-app/app/utils/convert"
	"yj-app/app/utils/excel"
	"yj-app/app/utils/page"
)

//根据键获取值
func GetValueByKey(key string) string {
	resultStr := ""
	result := gcache.Get(key)
	if result == nil {
		configRecord, err := configModel.FindOne("config_key", key)
		if err != nil {
			return ""
		}

		resultStr = configRecord.ConfigValue
		gcache.Set(key, resultStr, 0)
	} else {
		resultStr = result.(string)
	}

	return resultStr
}

func GetOssUrl() string {
	v := GetValueByKey("sys.resource.url")
	if v == "null" {
		return ""
	}
	return v
}

//根据主键查询数据
func SelectRecordById(id int64) (*configModel.Entity, error) {
	return configModel.FindOne("config_id", id)
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	entity, _ := configModel.FindOne("config_id=?", id)
	if entity != nil {
		result, err := configModel.Delete("config_id=?", id)
		if err == nil {
			affected, _ := result.RowsAffected()
			if affected > 0 {
				gcache.Remove(entity.ConfigKey)
				return true
			}
		}
	}
	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	list, _ := configModel.FindAll("config_id in (?)", idarr)
	result, err := configModel.Delete("config_id in (?)", idarr)
	if err != nil {
		return 0
	}

	if len(list) > 0 {
		for _, item := range list {
			gcache.Remove(item.ConfigKey)
		}
	}

	nums, _ := result.RowsAffected()

	return nums
}

//添加数据
func AddSave(req *configModel.AddReq, session *ghttp.Session) (int64, error) {
	var entity configModel.Entity
	entity.ConfigName = req.ConfigName
	entity.ConfigKey = req.ConfigKey
	entity.ConfigType = req.ConfigType
	entity.ConfigValue = req.ConfigValue
	entity.Remark = req.Remark
	entity.CreateTime = gtime.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(session)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	result, err := entity.Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}

//修改数据
func EditSave(req *configModel.EditReq, session *ghttp.Session) (int64, error) {

	entity, err := configModel.FindOne("config_id=?", req.ConfigId)

	if err != nil {
		return 0, err
	}

	if entity == nil {
		return 0, gerror.New("数据不存在")
	}

	entity.ConfigName = req.ConfigName
	entity.ConfigKey = req.ConfigKey
	entity.ConfigValue = req.ConfigValue
	entity.Remark = req.Remark
	entity.ConfigType = req.ConfigType
	entity.UpdateTime = gtime.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(session)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	result, err := entity.Update()

	if err != nil {
		return 0, err
	}

	rs, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	gcache.Set(entity.ConfigKey, entity.ConfigValue, 0)

	return rs, nil
}

//根据条件分页查询角色数据
func SelectListAll(params *configModel.SelectPageReq) ([]configModel.Entity, error) {
	return configModel.SelectListAll(params)
}

//根据条件分页查询角色数据
func SelectListByPage(params *configModel.SelectPageReq) ([]configModel.Entity, *page.Paging, error) {
	return configModel.SelectListByPage(params)
}

// 导出excel
func Export(param *configModel.SelectPageReq) (string, error) {
	result, err := configModel.SelectListExport(param)
	if err != nil {
		return "", err
	}

	head := []string{"参数主键", "参数名称", "参数键名", "参数键值", "系统内置（Y是 N否）", "状态"}
	key := []string{"config_id", "config_name", "config_key", "config_value", "config_type"}
	url, err := excel.DownlaodExcel(head, key, result)

	if err != nil {
		return "", err
	}

	return url, nil
}

//检查角色名是否唯一
func CheckConfigKeyUniqueAll(configKey string) string {
	entity, err := configModel.CheckPostCodeUniqueAll(configKey)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.ConfigId > 0 {
		return "1"
	}
	return "0"
}

//检查岗位名称是否唯一
func CheckConfigKeyUnique(configKey string, configId int64) string {
	entity, err := configModel.CheckPostCodeUnique(configKey, configId)
	if err != nil {
		return "1"
	}
	if entity != nil && entity.ConfigId > 0 {
		return "1"
	}
	return "0"
}
