package controllers

import (
	"fmt"
	//"log"

	"github.com/guerillagrow/gobox/models"

	//"github.com/asdine/storm/q"

	"encoding/json"
	//"errors"

	"time"

	"github.com/astaxie/beego"
	"gobot.io/x/gobot/drivers/gpio"
)

type ServiceRelay struct {
	beego.Controller
}

func (c *ServiceRelay) Prepare() {
	c.Data["json"] = map[interface{}]interface{}{
		"data": nil,
	}

}

func (c *ServiceRelay) Post() {

	var res JSONResp

	relayName := c.GetString("target")
	var relayDevice *gpio.GroveRelayDriver

	if relayName != "l1" && relayName != "l2" {
		c.Abort("500")
	}

	if relayName == "l1" {
		relayDevice = models.GoBox.RelayL1
	} else {
		relayDevice = models.GoBox.RelayL2
	}
	if relayDevice == nil {
		c.Abort("500")
	}

	csrfErr := CSRF.ValidateToken(fmt.Sprintf("svc"), c.GetString("__csrf__"), c.Ctx)

	if csrfErr != nil {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   500,
				"errors":   csrfErr,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	reqs := FormRelayL1{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil {
		c.Abort("500")
	}

	verr := reqs.Validate()

	if verr != nil {

		tOn, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/on", relayName))
		tOff, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/off", relayName))

		res = JSONResp{
			Data: map[string]interface{}{
				"state": relayDevice.State(),
				"ton":   tOn,
				"toff":  tOff,
			},
			Meta: map[string]interface{}{
				"status":   500,
				"errors":   verr,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()

		return
	}

	models.BoxConfig.Set(fmt.Sprintf("devices/relay_%s/settings/on", relayName), reqs.TOn)
	models.BoxConfig.Set(fmt.Sprintf("devices/relay_%s/settings/off", relayName), reqs.TOff)
	models.BoxConfig.Set(fmt.Sprintf("devices/relay_%s/settings/condition", relayName), reqs.Cond)
	models.BoxConfig.SetInt64(fmt.Sprintf("devices/relay_%s/settings/force", relayName), reqs.Force)
	models.BoxConfig.SaveFilePretty(models.BoxConfig.File)
	tOn, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/on", relayName))
	tOff, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/off", relayName))
	force, _ := models.BoxConfig.GetInt64(fmt.Sprintf("devices/relay_%s/settings/force", relayName))
	cond, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/condition", relayName))

	res = JSONResp{
		Data: map[string]interface{}{
			"status": relayDevice.State(),
			"ton":    tOn,
			"toff":   tOff,
			"cond":   cond,
			"force":  force,
		},
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
		},
	}
	c.Data["json"] = res

	c.ServeJSON()

}

func (c *ServiceRelay) Get() {
	c.StartSession()

	var res JSONResp

	relayName := c.GetString("target")
	var relayDevice *gpio.GroveRelayDriver

	if relayName == "l1" {
		relayDevice = models.GoBox.RelayL1
	} else if relayName == "l2" {
		relayDevice = models.GoBox.RelayL2
	} else {
		c.Abort("500")
	}

	tOn, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/on", relayName))
	tOff, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/off", relayName))
	force, _ := models.BoxConfig.GetInt64(fmt.Sprintf("devices/relay_%s/settings/force", relayName))
	cond, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/condition", relayName))

	res = JSONResp{
		Data: map[string]interface{}{
			"state": relayDevice.State(),
			"ton":   tOn,
			"toff":  tOff,
			"force": force,
			"cond":  cond,
		},
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}
