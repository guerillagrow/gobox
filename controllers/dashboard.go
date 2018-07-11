package controllers

import (
	"github.com/astaxie/beego"
	"github.com/guerillagrow/gobox/models"
)

type DashboarController struct {
	beego.Controller
}

func (c *DashboarController) Get() {
	isAdmin, ok := c.GetSession("user/isadmin").(bool)
	if !ok {
		c.Abort("500")
	}
	c.Data["sensor_t1"], _ = models.BoxConfig.GetBool("devices/t1/status")
	c.Data["sensor_t2"], _ = models.BoxConfig.GetBool("devices/t2/status")
	c.Data["sensor_d1"], _ = models.BoxConfig.GetBool("devices/d1/status")
	c.Data["sensor_d2"], _ = models.BoxConfig.GetBool("devices/d2/status")
	c.Data["user_isadmin"] = isAdmin
	c.TplName = "dashboard.tpl"
}
