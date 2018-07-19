package controllers

import (
	"github.com/guerillagrow/gobox/models"

	//"github.com/guerillagrow/beego"
	"github.com/guerillagrow/beego/context"
)

func init() {
	CSRF = &CSRFManager{}
}

func GetUserInfo(c *context.Context) models.User {
	u := models.User{}

	id, _ := c.Input.Session("user/id").(int64)
	name, _ := c.Input.Session("user/name").(string)
	email, _ := c.Input.Session("user/email").(string)
	isAdmin, _ := c.Input.Session("user/isadmin").(bool)

	u.ID = id
	u.Name = name
	u.Email = email
	u.IsAdmin = isAdmin

	return u

}
