package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	arrow "github.com/bmuller/arrow/lib"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/guerillagrow/beego"
	"github.com/guerillagrow/gobox/models"
	//"github.com/go-ozzo/ozzo-validation/is"
)

type ExportFile struct {
	File    string    `json:"file"`
	Created time.Time `json:"created"`
}

type ExportRequest struct {
	From           string        `json:"from"`
	To             string        `json:"to"`
	StatsTime      time.Duration `json:"stats_time"`
	Sensors        []string      `json:"sensors"`
	DeleteExported bool          `json:"delete_exported"`
	SensorType     string        `json:"sensor_type"`
}

func (f ExportRequest) Validate() error {

	err := validation.Errors{
		"from": func() error {
			err := validation.Validate(f.From, validation.Required)
			if err != nil {
				return err
			}
			_, err = arrow.CParse("%Y-%m-%d %H:%M:%S", f.From)
			return err
		}(),
		"to": func() error {
			err := validation.Validate(f.To, validation.Required)
			if err != nil {
				return err
			}
			_, err = arrow.CParse("%Y-%m-%d %H:%M:%S", f.To)
			return err
		}(),
		"sensors": func() error {
			if len(f.Sensors) < 1 {
				return errors.New("Missing sensors. Min one sensor required!")
			}
			return nil
		}(),
		"sensor_type": func() error {
			err := validation.Validate(f.To, validation.Required, validation.In("tempterature", "humidity"))
			if err != nil {
				return err
			}
			return nil
		}(),
	}.Filter()

	return err
}

type ServiceExport struct {
	beego.Controller
}

func (c *ServiceExport) Get() {
	var res JSONResp

	files, err := ioutil.ReadDir("./export")
	if err != nil {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   200,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/export"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	exportFiles := []ExportFile{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".csv") {
			exportFiles = append(exportFiles, ExportFile{
				File:    f.Name(),
				Created: f.ModTime(),
			})
		}
	}

	res = JSONResp{
		Data: exportFiles,
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/export"), 24*time.Hour, c.Ctx),
		},
	}
	c.Data["json"] = res
	c.ServeJSON()

}

func (c *ServiceExport) Post() {

	var res JSONResp

	// !TODO
	//w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT") // make response  as auto download in browser
	//w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	reqs := ExportRequest{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil {
		c.Abort("500")
	}

	err = reqs.Validate()
	if err != nil {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   500,
				"errors":   err,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/export"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	fromDate, _ := arrow.CParse("%Y-%m-%d %H:%M:%S", reqs.From)
	toDate, _ := arrow.CParse("%Y-%m-%d %H:%M:%S", reqs.From)

	ctx, cancle := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))

	defer cancle()
	var exportFile string

	if reqs.SensorType == "temperature" {
		exportFile, err = models.ExportTemperature(ctx, reqs.Sensors, "./export", fromDate.Time, toDate.Time, reqs.DeleteExported, time.Duration(reqs.StatsTime))
	} else { // humidity
		exportFile, err = models.ExportHumidity(ctx, reqs.Sensors, "./export", fromDate.Time, toDate.Time, reqs.DeleteExported, time.Duration(reqs.StatsTime))
	}

	if err != nil {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   500,
				"errors":   err,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/export"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res = JSONResp{
		Data: map[string]interface{}{
			"export_file": exportFile,
		},
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/export"), 24*time.Hour, c.Ctx),
		},
	}

	c.Data["json"] = res
	c.ServeJSON()
}
