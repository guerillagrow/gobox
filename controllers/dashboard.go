package controllers

import (
	"github.com/astaxie/beego"
)

type DashboarController struct {
	beego.Controller
}

func (c *DashboarController) Get() {
	isAdmin, ok := c.GetSession("user/isadmin").(bool)
	if !ok {
		c.Abort("500")
	}
	c.Data["user_isadmin"] = isAdmin
	c.TplName = "dashboard.tpl"
}
