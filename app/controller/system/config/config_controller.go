package config

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"yj-app/app/model"
	configModel "yj-app/app/model/system/config"
	configService "yj-app/app/service/system/config"
	"yj-app/app/utils/response"
)

//列表页
func List(r *ghttp.Request) {
	response.BuildTpl(r, "system/config/list.html").WriteTpl()
}

//列表分页数据
func ListAjax(r *ghttp.Request) {
	var req *configModel.SelectPageReq
	//获取参数
	if err := r.Parse(&req); err != nil {
		response.ErrorResp(r).SetMsg(err.Error()).Log("参数管理", req).WriteJsonExit()
	}
	rows := make([]configModel.Entity, 0)
	result, page, err := configService.SelectListByPage(req)

	if err == nil && len(result) > 0 {
		rows = result
	}
	response.BuildTable(r, page.Total, rows).WriteJsonExit()
}

//新增页面
func Add(r *ghttp.Request) {
	response.BuildTpl(r, "system/config/add.html").WriteTpl()
}

//新增页面保存
func AddSave(r *ghttp.Request) {
	var req *configModel.AddReq
	//获取参数
	if err := r.Parse(&req); err != nil {
		response.ErrorResp(r).SetBtype(model.Buniss_Add).SetMsg(err.Error()).Log("参数管理", req).WriteJsonExit()
	}

	if configService.CheckConfigKeyUniqueAll(req.ConfigKey) == "1" {
		response.ErrorResp(r).SetBtype(model.Buniss_Add).SetMsg("参数键名已存在").Log("参数管理", req).WriteJsonExit()
	}

	rid, err := configService.AddSave(req, r.Session)

	if err != nil || rid <= 0 {
		response.ErrorResp(r).SetBtype(model.Buniss_Add).Log("参数管理", req).WriteJsonExit()
	}
	response.SucessResp(r).SetData(rid).Log("参数管理", req).WriteJsonExit()
}

//修改页面
func Edit(r *ghttp.Request) {
	id := r.GetQueryInt64("id")

	if id <= 0 {
		response.ErrorTpl(r).WriteTpl(g.Map{
			"desc": "参数错误",
		})
		return
	}

	entity, err := configService.SelectRecordById(id)

	if err != nil || entity == nil {
		response.ErrorTpl(r).WriteTpl(g.Map{
			"desc": "数据不存在",
		})
		return
	}

	response.BuildTpl(r, "system/config/edit.html").WriteTpl(g.Map{
		"config": entity,
	})
}

//修改页面保存
func EditSave(r *ghttp.Request) {
	var req *configModel.EditReq
	//获取参数
	if err := r.Parse(&req); err != nil {
		response.ErrorResp(r).SetBtype(model.Buniss_Edit).SetMsg(err.Error()).Log("参数管理", req).WriteJsonExit()
	}

	if configService.CheckConfigKeyUnique(req.ConfigKey, req.ConfigId) == "1" {
		response.ErrorResp(r).SetBtype(model.Buniss_Edit).SetMsg("参数键名已存在").Log("参数管理", req).WriteJsonExit()
	}

	rs, err := configService.EditSave(req, r.Session)

	if err != nil || rs <= 0 {
		response.ErrorResp(r).SetBtype(model.Buniss_Edit).Log("参数管理", req).WriteJsonExit()
	}
	response.SucessResp(r).SetBtype(model.Buniss_Edit).Log("参数管理", req).WriteJsonExit()
}

//删除数据
func Remove(r *ghttp.Request) {
	var req *model.RemoveReq
	//获取参数
	if err := r.Parse(&req); err != nil {
		response.ErrorResp(r).SetBtype(model.Buniss_Del).SetMsg(err.Error()).Log("参数管理", req).WriteJsonExit()
	}

	rs := configService.DeleteRecordByIds(req.Ids)

	if rs > 0 {
		response.SucessResp(r).SetBtype(model.Buniss_Del).Log("参数管理", req).WriteJsonExit()
	} else {
		response.ErrorResp(r).SetBtype(model.Buniss_Del).Log("参数管理", req).WriteJsonExit()
	}
}

//导出
func Export(r *ghttp.Request) {
	var req *configModel.SelectPageReq
	//获取参数
	if err := r.Parse(&req); err != nil {
		response.ErrorResp(r).Log("参数管理", req).WriteJsonExit()
	}
	url, err := configService.Export(req)

	if err != nil {
		response.ErrorResp(r).SetBtype(model.Buniss_Other).Log("参数管理", req).WriteJsonExit()
	}

	response.SucessResp(r).SetBtype(model.Buniss_Other).SetMsg(url).WriteJsonExit()
}

//检查参数键名是否已经存在不包括本参数
func CheckConfigKeyUnique(r *ghttp.Request) {
	var req *configModel.CheckConfigKeyReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteflnExit("1")
	}

	result := configService.CheckConfigKeyUnique(req.ConfigKey, req.ConfigId)

	r.Response.WritefExit(result)
}

//检查参数键名是否已经存在
func CheckConfigKeyUniqueAll(r *ghttp.Request) {
	var req *configModel.CheckPostCodeALLReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteflnExit("1")
	}

	result := configService.CheckConfigKeyUniqueAll(req.ConfigKey)

	r.Response.WritefExit(result)
}
