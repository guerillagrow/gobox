package routers

import (
	"github.com/guerillagrow/gobox/controllers"
	"github.com/guerillagrow/gobox/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/auth"
)

func CustomAuthFilter(ctx *context.Context) {
	a := &auth.BasicAuth{Secrets: models.UserAuth, Realm: "GoBox"}
	email := a.CheckAuth(ctx.Request)
	if email == "" {
		a.RequireAuth(ctx.ResponseWriter, ctx.Request)
	}

	_, iok := ctx.Input.Session("user/email").(string)
	if !iok {
		user, aerr := models.GetUserByEmail(email)
		if aerr != nil {
			a.RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
		ctx.Input.CruSession.Set("user/id", user.ID)
		ctx.Input.CruSession.Set("user/name", user.Name)
		ctx.Input.CruSession.Set("user/email", user.Email)
		ctx.Input.CruSession.Set("user/isadmin", user.IsAdmin)

		//ctx.Input.CruSession.SessionRelease()
	}

}

func init() {

	//authPlugin := auth.NewBasicAuthenticator(models.UserAuth, "GoBox Backoffice")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	beego.InsertFilter("*", beego.BeforeRouter, CustomAuthFilter)

	beego.Router("/", &controllers.DashboarController{})
	serviceNS := beego.NewNamespace("svc",
		beego.NSRouter("/sensors/temperature", &controllers.ServiceSensors{}, "get:GetTemp"),
		beego.NSRouter("/sensors/humidity", &controllers.ServiceSensors{}, "get:GetHumidity"),
		beego.NSRouter("/sys/time", &controllers.ServiceSys{}, "get:GetTime"),
		beego.NSRouter("/sys/pistats", &controllers.ServiceSys{}, "get:GetPiStats"),
		beego.NSRouter("/user", &controllers.ServiceUser{}),
		beego.NSRouter("/relay", &controllers.ServiceRelay{}))
	beego.AddNamespace(serviceNS)
	//beego.ErrorHandler("404", controllers.PageNotFound)
}
