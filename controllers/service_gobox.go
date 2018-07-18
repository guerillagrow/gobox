package controllers

import (
	"github.com/guerillagrow/beego"
	"github.com/guerillagrow/gobox/lib/utils"
	gversion "github.com/mcuadros/go-version"
	//"github.com/guerillagrow/gobox/models"
)

type ServiceGoBox struct {
	beego.Controller
}

func (c *ServiceGoBox) Get() {

	currentVersion, errs := utils.CurrentGoBoxVersion()
	localVersion := utils.LocalGoBoxVersion()

	updateable := gversion.Compare(localVersion, currentVersion, "<")

	res := JSONResp{
		Meta: map[string]interface{}{
			"status": 200,
			"errors": errs,
		},
		Data: map[string]interface{}{
			"current_version":  currentVersion,
			"local_version":    localVersion,
			"update_available": updateable,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()

}
