package controllers

import (
	"github.com/guerillagrow/gobox/models"

	//"github.com/asdine/storm/q"

	//"context"
	"errors"
	"fmt"
	"strings"
	"time"
	//"encoding/json"
	//"errors"

	"regexp"

	arrow "github.com/bmuller/arrow/lib"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/guerillagrow/beego"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// !TODO: remove some of the aweful duplicate code. Make a custom base controller class an add a JSONOutput method on it

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

type FormUser struct {
	//State bool   `json:"status"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
}

func (f FormUser) Validate() error {

	err := validation.Errors{
		"name": func() error {
			err := validation.Validate(f.Name, validation.Required, validation.Length(4, 15))
			if err != nil {
				return err
			}
			return nil
		}(),
		"email": func() error {
			if f.Email == "root@localhost" {
				return nil
			}
			err := validation.Validate(f.Email, validation.Required, is.Email)
			if err != nil {
				return err
			}

			return nil
		}(),
		"password": func() error {
			err := validation.Validate(f.Password, validation.Length(4, 255))
			if err != nil {
				return err
			}
			return nil
		}(),
		"current_password": func() error {
			err := validation.Validate(f.CurrentPassword, validation.Required, validation.Length(4, 255))
			if err != nil {
				return err
			}
			ok := models.UserAuth(f.Email, f.CurrentPassword)
			if !ok {
				return errors.New("Invalid password!")
			}
			return nil
		}(),
	}.Filter()

	return err
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

	if statsSource {
		cachVal, _ := models.GoBox.SensorCache.GetFloat64(fmt.Sprintf("%s/temp", strings.ToLower(sensor)))
		sT := models.Temperature{
			ID:      0,
			Sensor:  sensor,
			Created: time.Now(),
			Value:   cachVal,
		}
		res = append(models.TemperatureSlice{sT}, res...)
	}

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

	if statsSource {
		cachVal, _ := models.GoBox.SensorCache.GetFloat64(fmt.Sprintf("%s/hum", strings.ToLower(sensor)))
		sT := models.Humidity{
			ID:      0,
			Sensor:  sensor,
			Created: time.Now(),
			Value:   cachVal,
		}
		res = append(models.HumiditySlice{sT}, res...)
	}

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

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}
