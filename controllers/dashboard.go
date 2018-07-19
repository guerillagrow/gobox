package controllers

import (
	"github.com/guerillagrow/beego"
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
	c.Data["relay_l1"], _ = models.BoxConfig.GetBool("devices/relay_l1/status")
	c.Data["relay_l2"], _ = models.BoxConfig.GetBool("devices/relay_l2/status")

	c.Data["sensor_t1"], _ = models.BoxConfig.GetBool("devices/t1/status")
	c.Data["sensor_t2"], _ = models.BoxConfig.GetBool("devices/t2/status")

	c.Data["sensor_d1"], _ = models.BoxConfig.GetBool("devices/d1/status")
	c.Data["sensor_d2"], _ = models.BoxConfig.GetBool("devices/d2/status")

	c.Data["sensor_s1"], _ = models.BoxConfig.GetBool("devices/s1/status")
	c.Data["sensor_s2"], _ = models.BoxConfig.GetBool("devices/s2/status")

	c.Data["metric_source"], _ = models.BoxConfig.GetString("ui/metric_source")

	c.Data["user_isadmin"] = isAdmin
	c.TplName = "dashboard.tpl"
}
