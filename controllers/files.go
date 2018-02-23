//域名管理
package controllers

import (
	"Webgroup/models"
	"fmt"
	"tsEngine/tsDb"
	"tsEngine/tsFile"
	"tsEngine/tsString"
	"tsEngine/tsTime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type FilesController struct {
	BaseController
}

//类似构造函数
func (this *FilesController) Prepare() {
	//权限判断
	this.CheckPermission()
}

//软文列表
func (this *FilesController) List() {

	Keyword := this.GetString("Keyword")
	Page, _ := this.GetInt64("Page")
	PageSize, _ := this.GetInt64("PageSize")

	o := models.Files{}

	items, pagination, err := o.List(Page, PageSize, Keyword)

	if err != nil {
		beego.Error(err)
		this.Code = 0
		this.Msg = "数据库异常错误，请联系管理员"
		this.TraceJson()
	}

	this.Code = models.SuccessProto
	this.Result = map[string]interface{}{"Items": items, "Pagination": pagination}
	this.TraceJson()

}

func (this *FilesController) Add() {

	o := models.Files{}
	o.Title = this.GetString("Title")
	o.Content = this.GetString("Content")
	o.Sort, _ = this.GetInt64("Sort")
	o.Time = tsTime.CurrSe()

	//数据验证
	valid := validation.Validation{}

	valid.Required(o.Title, "Title").Message("标题不能为空")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			this.Code = 0
			this.Msg = err.Message
			this.TraceJson()
		}

	}

	db := tsDb.NewDbBase()
	_, err := db.DbInsert(&o)
	if err != nil {
		beego.Error(err)
		this.Code = 0
		this.Msg = "数据库操作异常，请联系管理员"
		this.TraceJson()
	}

	html := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
%s
</body>
</html>`
	html = fmt.Sprintf(html, o.Content)
	tsFile.WriteFile("./static/files/"+"1.html", html)
	this.Code = 1
	this.Result = o
	this.TraceJson()
}

func (this *FilesController) Edit() {

	//初始化
	db := tsDb.NewDbBase()
	o := models.Files{}

	//获取get数据
	o.Id = tsString.ToInt64(this.Ctx.Input.Param("0"))
	if o.Id > 0 {
		err := db.DbGet(&o)
		if err != nil {
			beego.Error(err)
			this.Code = 0
			this.Msg = "没有该记录"
			this.TraceJson()
		}

		this.Code = 1
		this.Result = o
		this.TraceJson()
	}

	//获取post数据
	o.Id, _ = this.GetInt64("Id")
	o.Title = this.GetString("Title")
	o.Sort, _ = this.GetInt64("Sort")
	o.Content = this.GetString("Content")

	//****************************************************
	//数据验证
	valid := validation.Validation{}

	valid.Required(o.Title, "Title").Message("标题不能为空")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			this.Code = 0
			this.Msg = err.Message
			this.TraceJson()
		}

	}

	//****************************************************
	err := db.DbUpdate(&o, "Title", "Content", "Sort")

	if err != nil {
		beego.Error(err)
		this.Code = 0
		this.Msg = "参数错误或没有任何修改"
		this.TraceJson()
	}

	this.Code = 1
	this.Result = o
	this.TraceJson()

}

func (this *FilesController) Del() {

	o := models.Files{}
	o.Id, _ = this.GetInt64("Id")
	db := tsDb.NewDbBase()
	err := db.DbRead(&o)
	if err != nil {
		this.Code = 0
		this.Msg = "参数错误"
		this.TraceJson()
	}

	err = db.DbDel(&o)
	if err != nil {
		beego.Error(err)
		this.Code = 0
		this.Msg = "数据库异常错误，请联系管理员"
		this.TraceJson()
	}
	this.Code = 1
	this.TraceJson()

}
