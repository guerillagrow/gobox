package controllers

import (
	"fmt"
	//"log"

	"github.com/guerillagrow/gobox/models"

	//"github.com/asdine/storm/q"

	"context"
	"time"

	"encoding/json"
	//"errors"

	"regexp"

	"github.com/astaxie/beego"
	arrow "github.com/bmuller/arrow/lib"
	"github.com/go-ozzo/ozzo-validation"
	//"github.com/go-ozzo/ozzo-validation/is"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"gobot.io/x/gobot/drivers/gpio"
)

type JSONResp struct {
	Meta map[string]interface{} `json:"meta"`
	Data interface{}            `json:"data"`
}

type FormRelayL1 struct {
	//State bool   `json:"status"`
	Force string `json:"force"`
	TOn   string `json:"ton"`
	TOff  string `json:"toff"`
	Cond  string `json:"cond""`
}

func (f FormRelayL1) Validate() error {

	err := validation.Errors{
		"ton": func() error {
			err := validation.Validate(f.TOn, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$")))
			if err != nil {
				return err
			}
			_, err = arrow.CParse("%H:%M", f.TOn)
			return err
		}(),
		"toff": func() error {
			err := validation.Validate(f.TOff, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$")))
			if err != nil {
				return err
			}
			_, err = arrow.CParse("%H:%M", f.TOff)
			return err
		}(),
		"force": func() error {
			err := validation.Validate(f.Force, validation.In("-1", "0", "1"))
			if err != nil {
				return err
			}
			return nil
		}(),
		"cond": func() error {
			if f.Cond != "" {
				_, err := models.GoBox.EvalRelayExpression(f.Cond)
				return err
			}
			return nil
		}(),
	}.Filter()

	return err
}

type ServiceRelay struct {
	beego.Controller
}

func (c *ServiceRelay) Prepare() {
	c.Data["json"] = map[interface{}]interface{}{
		"data": nil,
	}

}

func (c *ServiceRelay) Post() {

	/*csrfToken := c.GetString("__csrf__")
	csrfErr := CSRF.ValidateToken("svc/relay_l1", csrfToken, c.Ctx)

	if csrfErr != nil {
		c.Abort("500")
	}*/

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
	reqs := FormRelayL1{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil {
		c.Abort("500")
	}

	verr := reqs.Validate()

	if verr != nil {

		tOn, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/on", relayName))
		tOff, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/off", relayName))

		res := JSONResp{
			Data: map[string]interface{}{
				"state": relayDevice.State(),
				"ton":   tOn,
				"toff":  tOff,
			},
			Meta: map[string]interface{}{
				"status": 500,
				"errors": verr,
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

	res := JSONResp{
		Data: map[string]interface{}{
			"status": relayDevice.State(),
			"ton":    tOn,
			"toff":   tOff,
			"cond":   cond,
			"force":  force,
		},
		Meta: map[string]interface{}{
			"status": 200,
		},
	}
	c.Data["json"] = res

	c.ServeJSON()

}

func (c *ServiceRelay) Get() {
	c.StartSession()

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

	tOn, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/on", relayName))
	tOff, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/off", relayName))
	force, _ := models.BoxConfig.GetInt64(fmt.Sprintf("devices/relay_%s/settings/force", relayName))
	cond, _ := models.BoxConfig.GetString(fmt.Sprintf("devices/relay_%s/settings/condition", relayName))

	res := JSONResp{
		Data: map[string]interface{}{
			"state": relayDevice.State(),
			"ton":   tOn,
			"toff":  tOff,
			"force": force,
			"cond":  cond,
		},
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken("svc/relay_l1", 1*time.Hour, c.Ctx),
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// !TODO: use SensorsParams in GetTemp() & GetHumidity()
type SensorsParams struct {
	Sensor     string `form:"sensor" json:"sensor"`
	Limit      int    `form:"limit" json:"limit"`
	Graph      int    `form:"g" json:"g"`
	Order      string `form:"order" json:"order"`
	Stats      bool   `form:"stats" json:"stats"`
	StatsTime  bool   `form:"ts" json:"ts"`
	TimeLength string `form:"tl" json:"tl"`
	DateFrom   string `form:"from" json:"from"`
	DateTo     string `form:"to" json:"to"`
}

func (c SensorsParams) Validate() error {
	// !TODO
	return nil
}

type ServiceSensors struct {
	beego.Controller
}

func (c *ServiceSensors) GetTemp() {
	c.StartSession()

	var res models.TemperatureSlice

	limit, _ := c.GetInt("limit", 1000)
	graph, _ := c.GetInt("g", 0)
	order := c.GetString("order", "desc")
	statsSource, _ := c.GetBool("stats", false)
	//statsTime, _ := c.GetInt64("ts", 0)
	tl := c.GetString("tl", "")
	date1 := c.GetString("to", arrow.Now().CFormat("%Y-%m-%d %H:%M:%S"))
	date2 := c.GetString("from", arrow.Yesterday().CFormat("%Y-%m-%d %H:%M:%S"))

	sensor := c.GetString("sensor", "T1")
	rbuckets := models.DB.From("sensd", "temperature").PrefixScan("T")
	buckets := []string{}
	for _, b := range rbuckets {
		bl := b.Bucket()
		buckets = append(buckets, bl[len(bl)-1])
	}

	if !contains(buckets, sensor) {
		c.Abort("500")
	}
	//node := models.DB.From("sensd", "temperature", sensor)

	statsTime := time.Duration(0)

	if tl == "day" {
		date1 = arrow.Now().CFormat("%Y-%m-%d %H:%M:%S")
		date2 = arrow.Yesterday().CFormat("%Y-%m-%d %H:%M:%S")
		statsTime = time.Duration(10) * time.Minute

	} else if tl == "hour" {
		date1 = arrow.Now().CFormat("%Y-%m-%d %H:%M:%S")
		date2 = arrow.Now().AddHours(-3).CFormat("%Y-%m-%d %H:%M:%S")
		statsTime = time.Duration(1) * time.Minute

	}

	dateTo, d1err := arrow.CParse("%Y-%m-%d %H:%M:%S", date1)
	dateFrom, d2err := arrow.CParse("%Y-%m-%d %H:%M:%S", date2)
	if d1err != nil || d2err != nil {
		c.Abort("500")
	}
	orderAsc := false
	if order == "asc" {
		orderAsc = true
	}

	res, _ = models.QueryTemperatureData(
		sensor,
		dateFrom.Time,
		dateTo.Time,
		limit,
		orderAsc,
		statsSource,
		statsTime,
	)

	if graph == 1 {
		var resGraph [][2]uint64
		for _, r := range res {
			resGraph = append(resGraph, [2]uint64{uint64(r.Created.Unix() * 1000), uint64(r.Value)})
		}
		c.Data["json"] = JSONResp{
			Meta: map[string]interface{}{
				"from":  date2,
				"to":    date1,
				"limit": limit,
			},
			Data: resGraph,
		}
	} else {
		c.Data["json"] = JSONResp{
			Meta: map[string]interface{}{
				"from":  date2,
				"to":    date1,
				"limit": limit,
			},
			Data: res,
		}
	}
	c.ServeJSON()
	return
}

func (c *ServiceSensors) GetHumidity() {
	c.StartSession()

	var res models.HumiditySlice

	limit, _ := c.GetInt("limit", 1000)
	graph, _ := c.GetInt("g", 0)
	order := c.GetString("order", "desc")
	statsSource, _ := c.GetBool("stats", false)
	tl := c.GetString("tl", "")
	date1 := c.GetString("to", arrow.Now().CFormat("%Y-%m-%d %H:%M:%S"))
	date2 := c.GetString("from", arrow.Yesterday().CFormat("%Y-%m-%d %H:%M:%S"))

	sensor := c.GetString("sensor", "T1")
	rbuckets := models.DB.From("sensd", "temperature").PrefixScan("T")
	buckets := []string{}
	for _, b := range rbuckets {
		bl := b.Bucket()
		buckets = append(buckets, bl[len(bl)-1])
	}

	if !contains(buckets, sensor) {
		c.Abort("500")
	}

	statsTime := time.Duration(0)

	if tl == "day" {
		date1 = arrow.Now().CFormat("%Y-%m-%d %H:%M:%S")
		date2 = arrow.Yesterday().CFormat("%Y-%m-%d %H:%M:%S")
		statsTime = time.Duration(10) * time.Minute

	} else if tl == "hour" {
		date1 = arrow.Now().CFormat("%Y-%m-%d %H:%M:%S")
		date2 = arrow.Now().AddHours(-3).CFormat("%Y-%m-%d %H:%M:%S")
		statsTime = time.Duration(1) * time.Minute

	}

	dateTo, d1err := arrow.CParse("%Y-%m-%d %H:%M:%S", date1)
	dateFrom, d2err := arrow.CParse("%Y-%m-%d %H:%M:%S", date2)
	if d1err != nil || d2err != nil {
		c.Abort("500")
	}
	orderAsc := false
	if order == "asc" {
		orderAsc = true
	}

	res, _ = models.QueryHumidityData(
		sensor,
		dateFrom.Time,
		dateTo.Time,
		limit,
		orderAsc,
		statsSource,
		statsTime,
	)

	if graph == 1 {
		var resGraph [][2]uint64
		for _, r := range res {
			resGraph = append(resGraph, [2]uint64{uint64(r.Created.Unix() * 1000), uint64(r.Value)})
		}
		c.Data["json"] = JSONResp{
			Meta: map[string]interface{}{
				"from":  date2,
				"to":    date1,
				"limit": limit,
			},
			Data: resGraph,
		}
	} else {
		c.Data["json"] = JSONResp{
			Meta: map[string]interface{}{
				"from":  date2,
				"to":    date1,
				"limit": limit,
			},
			Data: res,
		}
	}
	c.ServeJSON()
	return

}

type ServiceUser struct {
	beego.Controller
}

func (c *ServiceUser) Get() {
	c.Abort("500") // ! BLOCKED

}

func (c *ServiceUser) ChangeUser() {
	c.Abort("500") // ! BLOCKED

}

func (c *ServiceUser) CreateUser() {
	c.Abort("500") // ! BLOCKED
	// models.NewUser(name, email, password, isAdmin)

	// !TODO: add isAdmin check but allow users to change their own password, email and username

	cuser := GetUserInfo(c.Ctx)

	if cuser.Email == "" { // check again just for safty
		c.Abort("500")
	}

	reqs := models.User{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil {
		c.Abort("500")
	}

	verr := reqs.Validate()

	if verr != nil {

		res := JSONResp{
			Data: map[string]interface{}{},
			Meta: map[string]interface{}{
				"status": 500,
				"errors": verr,
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	tuser := models.User{}
	models.DB.One("Email", reqs.Email, &tuser)

	// Check permission

	if !cuser.IsAdmin && reqs.Email != cuser.Email {
		c.Abort("403")
	}

	serr := reqs.Save()

	if serr != nil {

		res := JSONResp{
			Data: map[string]interface{}{},
			Meta: map[string]interface{}{
				"status": 500,
				"errors": serr,
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res := JSONResp{
		Data: reqs,
		Meta: map[string]interface{}{
			"status": 200,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}

type ServiceSys struct {
	beego.Controller
}

func (c *ServiceSys) Get() {
	c.Abort("500") // ! BLOCKED
}

func (c *ServiceSys) GetPiStats() {
	dstats, derr := disk.Usage("/")
	mstats, merr := mem.VirtualMemory()
	if derr != nil || merr != nil {
		c.Abort("500")
	}
	res := JSONResp{
		Data: map[string]interface{}{
			"disk_total":        dstats.Total,
			"disk_used":         dstats.Used,
			"disk_used_percent": dstats.UsedPercent,
			"disk_free":         dstats.Free,
			"mem_total":         mstats.Total,
			"mem_used":          mstats.Used,
			"mem_used_percent":  mstats.UsedPercent,
			"mem_free":          mstats.Free,
		},
		Meta: map[string]interface{}{
			"status": 200,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func (c *ServiceSys) GetTime() {
	res := JSONResp{
		Data: map[string]interface{}{
			"t0": time.Now().Unix(),
			"t1": time.Now().UnixNano(),
			"t2": time.Now(),
			"t3": arrow.Now().CFormat("%Y-%m-%d %H:%M:%S"),
		},
		Meta: map[string]interface{}{
			"status": 200,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func (c *ServiceSys) Post() {
	c.Abort("500") // ! BLOCKED
}

type ServiceExport struct {
	beego.Controller
}

func (c *ServiceExport) Get() {
	c.Abort("500") // ! BLOCKED
}

func (c *ServiceExport) Post() {
	// !TODO
	c.Abort("500") // ! BLOCKED

	sensorType := c.GetString("type", "")
	sensors := c.GetStrings("sensors", []string{})
	fromRaw := c.GetString("from")
	toRaw := c.GetString("from")
	deleteExported, _ := c.GetBool("delete_exported", false)
	statsTimeRaw, _ := c.GetInt64("stats_time", 0)
	statsTime := time.Duration(statsTimeRaw)
	fromDate, fdErr := arrow.CParse("%Y-%m-%d %H:%M:%S", fromRaw)
	toDate, tdErr := arrow.CParse("%Y-%m-%d %H:%M:%S", toRaw)

	if (sensorType != "temperature" && sensorType != "humidity") || len(sensors) < 1 || fdErr != nil || tdErr != nil {
		c.Abort("500")
	}

	ctx, cancle := context.WithDeadline(
		context.Background(),
		time.Now().Add(1*time.Minute),
	)
	defer cancle()

	var resp string
	var rErr error

	if sensorType == "temperature" {
		resp, rErr = models.ExportTemperature(ctx, sensors, "./export", fromDate.Time, toDate.Time, deleteExported, statsTime)
	} else {
		resp, rErr = models.ExportHumidity(ctx, sensors, "./export", fromDate.Time, toDate.Time, deleteExported, statsTime)
	}

	if rErr != nil {
		// !TODO
		c.Abort("500")
	}

	res := JSONResp{
		Data: map[string]interface{}{
			"file": resp,
		},
		Meta: map[string]interface{}{
			"status": 200,
		},
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}
