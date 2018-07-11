package controllers

import (
	//"fmt"
	"github.com/guerillagrow/gobox/models"

	//"github.com/asdine/storm/q"

	"time"

	"encoding/json"
	//"errors"

	"regexp"

	"github.com/astaxie/beego"
	arrow "github.com/bmuller/arrow/lib"
	"github.com/go-ozzo/ozzo-validation"
	//"github.com/go-ozzo/ozzo-validation/is"
	//"github.com/go-ozzo/ozzo-validation/is"
)

type JSONResp struct {
	Meta map[string]interface{} `json:"meta"`
	Data interface{}            `json:"data"`
}

type FormRelayL1 struct {
	//State bool   `json:"status"`
	TOn  string `json:"ton"`
	TOff string `json:"toff"`
	Cond string `json:cond"`
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
	}.Filter()
	return err

	/*return validation.ValidateStruct(&f,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&f.State, validation.Required),
		// City cannot be empty, and the length must between 5 and 50
		//validation.Field(&f.TOn, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&f.TOn, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$"))),
		validation.Field(&f.TOff, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)*/
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

	reqs := FormRelayL1{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil {
		c.Abort("500")
	}

	verr := reqs.Validate()

	if verr != nil {

		tOn, _ := models.BoxConfig.GetString("devices/relay_l1/settings/on")
		tOff, _ := models.BoxConfig.GetString("devices/relay_l1/settings/off")

		res := JSONResp{
			Data: map[string]interface{}{
				"state": models.GoBox.LightState(),
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

	models.BoxConfig.Set("devices/relay_l1/settings/on", reqs.TOn)
	models.BoxConfig.Set("devices/relay_l1/settings/off", reqs.TOff)
	models.BoxConfig.SaveFilePretty(models.BoxConfig.File)
	tOn, _ := models.BoxConfig.GetString("devices/relay_l1/settings/on")
	tOff, _ := models.BoxConfig.GetString("devices/relay_l1/settings/off")

	res := JSONResp{
		Data: map[string]interface{}{
			"status": models.GoBox.LightState(),
			"ton":    tOn,
			"toff":   tOff,
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

	tOn, _ := models.BoxConfig.GetString("devices/relay_l1/settings/on")
	tOff, _ := models.BoxConfig.GetString("devices/relay_l1/settings/off")

	res := JSONResp{
		Data: map[string]interface{}{
			"state": models.GoBox.LightState(),
			"ton":   tOn,
			"toff":  tOff,
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
		var resGraph [][2]int
		for _, r := range res {
			resGraph = append(resGraph, [2]int{int(r.Created.Unix() * 1000), int(r.Value)})
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
		var resGraph [][2]int
		for _, r := range res {
			resGraph = append(resGraph, [2]int{int(r.Created.Unix() * 1000), int(r.Value)})
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

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}
