package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guerillagrow/beego"
	"github.com/guerillagrow/gobox/models"
)

type ServiceUser struct {
	beego.Controller
}

func (c *ServiceUser) Get() {

	user, _ := c.GetSession("user").(models.User)

	var res JSONResp
	if user.ID < 1 {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   500,
				"userID":   user.ID,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
	}
	user.PwHash = ""
	res = JSONResp{
		Data: user,
		Meta: map[string]interface{}{
			"status":   200,
			"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
		},
	}
	c.Data["json"] = res
	c.ServeJSON()

}

func (c *ServiceUser) Post() {

	var res JSONResp

	csrfErr := CSRF.ValidateToken(fmt.Sprintf("svc"), c.GetString("__csrf__"), c.Ctx)

	if csrfErr != nil {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":          500,
				"errors":          csrfErr,
				"__csrf_error___": true,
				"__csrf__":        CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	user, _ := c.GetSession("user").(models.User)

	if user.ID < 1 {
		res = JSONResp{
			Data: nil,
			Meta: map[string]interface{}{
				"status":   500,
				"userID":   user.ID,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	reqs := FormUser{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqs)

	if err != nil || user.ID < 1 {
		c.Abort("500")
	}

	err = reqs.Validate()
	fmt.Println("svc/user validate:", err)

	if err != nil {
		res = JSONResp{
			Data: reqs,
			Meta: map[string]interface{}{
				"status":   500,
				"errors":   err,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc"), 24*time.Hour, c.Ctx),
			},
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	user.Name = reqs.Name
	user.Email = reqs.Email

	if reqs.Password != "" {
		user.Password = reqs.Password
	}

	err = user.Save()
	reqs.Password = ""
	fmt.Println("svc/user saved user:", user.Email)

	if err != nil {
		res = JSONResp{
			Data: reqs,
			Meta: map[string]interface{}{
				"status":   500,
				"error":    err,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/user"), 24*time.Hour, c.Ctx),
			},
		}
	} else {

		res = JSONResp{
			Data: user,
			Meta: map[string]interface{}{
				"status":   200,
				"__csrf__": CSRF.SetToken(fmt.Sprintf("svc/user"), 24*time.Hour, c.Ctx),
			},
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
}
